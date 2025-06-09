package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	shared "github.com/sparrow-hkr/oto/internal/Shared"
	"github.com/sparrow-hkr/oto/internal/processUrls"
	"github.com/spf13/cobra"
)

var (
	domain      string
	domainFile  string
	outputFile  string
	resultTypes = []string{}
	concurrency int
	timeout     time.Duration
	verbose     bool
	debug       bool
)

var endpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Extract endpoints from HTML/JS source.",
	Long:  `Extract endpoints from HTML/JS source for a single or multiple domains.`,
	Run: func(cmd *cobra.Command, args []string) {
		if domain == "" && domainFile == "" {
			fmt.Println("Please provide either --domain (-d) or --list (-l) flag.")
			os.Exit(1)
		}
		if domain != "" {
			processSingleDomain(domain)
		} else if domainFile != "" {
			processDomainList(domainFile)
		}
	},
}

func processDomainList(domainsFile string) {
	fmt.Printf("[%s*%s] Processing domain list from file: %s\n", shared.ColorGreen, shared.ColorReset, domainsFile)
	file, err := os.Open(domainsFile)
	if err != nil {
		fmt.Printf("Error opening domain file: %v\n", err)
		return
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Normalize: add https:// if missing
		if !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
			line = "https://" + line
		}
		urls = append(urls, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading domain file: %v\n", err)
		return
	}

	processUrls.ProcessURLs(urls, resultTypes, outputFile, concurrency, timeout, verbose, debug)
}
func processSingleDomain(domain string) {
	var urls []string
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		urls = append(urls, "https://"+domain)
	} else {
		urls = append(urls, domain)
	}
	processUrls.ProcessURLs(urls, resultTypes, outputFile, concurrency, timeout, verbose, debug)
}

func init() {
	endpointCmd.Flags().StringVarP(&domain, "domain", "d", "", "Target domain (required if --list not set)")
	endpointCmd.Flags().StringVarP(&domainFile, "list", "l", "", "File containing domains list (required if --domain not set)")
	endpointCmd.Flags().StringSliceVarP(&resultTypes, "result-types", "t", []string{"endpoint", "path", "info", "critical", "sensitive"}, "Types of results to extract [Posibale Value: endpoint, path, info, critical, sensitive]")
	endpointCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 5, "Number of concurrent threads (default: 10)")
	endpointCmd.Flags().DurationVarP(&timeout, "timeout", "T", 5*time.Second, "Timeout for HTTP requests (default: 5s)")
	endpointCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	endpointCmd.Flags().BoolVarP(&debug, "debug", "D", false, "Enable debug output")
	endpointCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file to save results (default: stdout)")
	rootCmd.AddCommand(endpointCmd)
}
