package processUrls

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	shared "github.com/sparrow-hkr/oto/internal/Shared"
	"github.com/sparrow-hkr/oto/internal/banner"
)

var (
	Red    = shared.ColorRed
	Green  = shared.ColorGreen
	Yellow = shared.ColorYellow
	Cyan   = shared.ColorCyan
	Reset  = shared.ColorReset
)

var (
	htmlTagFilter = shared.Patterns.HTMLTagFilter
	reEndpoint    = shared.Patterns.Endpoint
	rePath        = shared.Patterns.Path
	reInfo        = shared.Patterns.Info
	reCritical    = shared.Patterns.CriticalPath
	reSensitive   = shared.Patterns.SensitiveKeywords
	reScript      = shared.Patterns.Script
)

func filterHtmlTags(items []string) []string {
	var filtered []string
	for _, item := range items {
		if !htmlTagFilter.MatchString(strings.ToLower(item)) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// ExtractJSUrls extracts all JS URLs from HTML body and normalizes them.
func ExtractJSUrls(htmlBody []byte, baseURL string, verbose bool) []string {
	Script := reScript
	matches := Script.FindAllStringSubmatch(string(htmlBody), -1)
	var jsUrls []string
	for _, match := range matches {
		jsUrl := match[1]
		// Normalize relative URLs
		if strings.HasPrefix(jsUrl, "//") {
			parsed, _ := url.Parse(baseURL)
			jsUrl = parsed.Scheme + ":" + jsUrl
		} else if strings.HasPrefix(jsUrl, "/") {
			parsed, _ := url.Parse(baseURL)
			jsUrl = parsed.Scheme + "://" + parsed.Host + jsUrl
		} else if !strings.HasPrefix(jsUrl, "http://") && !strings.HasPrefix(jsUrl, "https://") {
			parsed, _ := url.Parse(baseURL)
			jsUrl = parsed.Scheme + "://" + parsed.Host + "/" + strings.TrimLeft(jsUrl, "/")
		}
		jsUrls = append(jsUrls, jsUrl)
		if verbose {
			fmt.Printf("%s-%s Extracted JS URL: %s\n", Yellow, Reset, jsUrl)
		}
	}
	return jsUrls
}

// ProcessURLs fetches HTML source from given URLs, extracts endpoints/paths, and also processes JS files.
func ProcessURLs(urls []string, resultTypes []string, outputFile string, concurrency int, timeout time.Duration, verbose bool, debug bool) {
	var results []shared.Result
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	banner.PrintProcessMessage(urls, resultTypes, outputFile, concurrency, timeout, verbose, debug)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}

	// Deduplication map
	seen := make(map[string]struct{})
	var allUrls []string

	// First pass: fetch HTML, extract JS URLs, and collect all URLs
	for _, urlStr := range urls {
		urlStr = strings.TrimSpace(urlStr)
		time.Sleep(500 * time.Millisecond)
		parsed, err := url.Parse(urlStr)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			if debug {
				fmt.Printf("- Invalid URL: %s, skipping...\n", urlStr)
			}
			continue
		}
		if _, exists := seen[urlStr]; exists {
			continue
		}
		seen[urlStr] = struct{}{}
		allUrls = append(allUrls, urlStr)

		resp, err := client.Get(urlStr)
		if err != nil {
			if debug {
				fmt.Printf("- Error fetching URL: %s, %v\n", urlStr, err)
			}
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			if debug {
				fmt.Printf("- Error reading response body for URL: %s, %v\n", urlStr, err)
			}
			continue
		}
		// Extract JS URLs and add to allUrls if not seen
		jsUrls := ExtractJSUrls(body, urlStr, verbose)
		for _, jsUrl := range jsUrls {
			if _, exists := seen[jsUrl]; !exists {
				seen[jsUrl] = struct{}{}
				allUrls = append(allUrls, jsUrl)
			}
		}
	}

	// Second pass: process all deduplicated URLs (HTML and JS)
	for _, urlStr := range allUrls {
		wg.Add(1)
		sem <- struct{}{}
		go func(urlStr string) {
			defer wg.Done()
			defer func() { <-sem }()
			resp, err := client.Get(urlStr)
			if err != nil {
				if debug {
					fmt.Printf("- Error fetching URL: %s, %v\n", urlStr, err)
				}
				return
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				if debug {
					fmt.Printf("Error reading response body for URL: %s, %v\n", urlStr, err)
				}
				return
			}
			if verbose {
				fmt.Printf("+ Fetched source for URL: %s [%v]\n", urlStr, resp.StatusCode)
			}
			result := shared.Result{URL: urlStr, Endpoints: []string{}, Paths: []string{}}
			for _, t := range resultTypes {
				switch strings.ToLower(t) {
				case "endpoint":
					result.Endpoints = reEndpoint.FindAllString(string(body), -1)
				case "path":
					result.Paths = rePath.FindAllString(string(body), -1)
				case "info":
					result.Info = reInfo.FindAllString(string(body), -1)
				case "critical":
					result.CriticalPaths = reCritical.FindAllString(string(body), -1)
				case "sensitive":
					result.SensitiveKeywords = reSensitive.FindAllString(string(body), -1)
				}
			}
			result.Endpoints = deduplicate(filterHtmlTags(result.Endpoints))
			result.Paths = deduplicate(filterHtmlTags(result.Paths))
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(urlStr)
	}
	wg.Wait()

	// Print results as JSON
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		if debug {
			fmt.Printf("Error marshalling results to JSON: %v\n", err)
		}
		return
	}
	if outputFile != "" {
		err := os.WriteFile(outputFile, data, 0644)
		if err != nil {
			if debug {
				fmt.Printf("Error writing results to file %s: %v\n", outputFile, err)
			}
			return
		}
		if verbose {
			fmt.Println(string(data))
			fmt.Printf("[%s+%s] Results written to: %s\n", Green, Reset, outputFile)
			fmt.Printf("[%s+%s] Results crecter byte: %s%v%s\n", Green, Reset, Green, len(data), Reset)

		} else {
			fmt.Printf("[%s+%s] Results written to: %s\n", Green, Reset, outputFile)
			fmt.Printf("[%s+%s] Results crecter byte: %s%v%s\n", Green, Reset, Green, len(data), Reset)
		}
	} else {
		fmt.Println(string(data))
		fmt.Printf("[%s+%s] Results crecter byte: %s%v%s\n", Green, Reset, Green, len(data), Reset)

	}
}

func deduplicate(items []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
