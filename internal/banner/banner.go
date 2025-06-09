package banner

import (
	"fmt"
	"time"

	shared "github.com/sparrow-hkr/oto/internal/Shared"
)

var (
	Red    = shared.ColorRed
	Green  = shared.ColorGreen
	Yellow = shared.ColorYellow
	Cyan   = shared.ColorCyan
	Reset  = shared.ColorReset
)

func PrintBanner() {
	banner := fmt.Sprintf(`
	     ot%so%s
	  6%sF%s 74 6F
	111 116 11%s1%s   
011%so%s1111 0%st%s110100 011%so%s1111	%sv1.0.0%s
`, Yellow, Reset, Green, Reset, Red, Reset, Red, Reset, Yellow, Reset, Cyan, Reset, Cyan, Reset)
	fmt.Println(banner)
}
func PrintProcessMessage(urls []string, resultTypes []string, outputFile string, concurrency int, timeout time.Duration, verbose bool, debug bool) {
	fmt.Printf("[%s*%s] Processing URLs: %d\n", Green, Reset, len(urls))
	fmt.Printf("[%s*%s] Result Types: %s\n", Green, Reset, resultTypes)
	fmt.Printf("[%s*%s] Output File: %s\n", Green, Reset, outputFile)
	fmt.Printf("[%s*%s] Concurrency: %d\n", Green, Reset, concurrency)
	fmt.Printf("[%s*%s] Timeout: %d ms\n", Green, Reset, timeout)
	if verbose {
		fmt.Printf("[%s*%s] Verbose mode %senabled%s\n", Green, Reset, Green, Reset)
	} else {
		fmt.Printf("[%s*%s] Verbose mode %sdisabled%s\n", Red, Reset, Red, Reset)
	}
	if debug {
		fmt.Printf("[%s*%s] Debug mode %senabled%s\n", Green, Reset, Green, Reset)
	} else {
		fmt.Printf("[%s*%s] Debug mode %sdisabled%s\n", Red, Reset, Red, Reset)
	}
	fmt.Printf("%sStarting processing...%s\n", Green, Reset)
	time.Sleep(1 * time.Second)
	if len(urls) == 0 {
		fmt.Printf("%sNo URLs provided. Exiting.%s\n", Red, Reset)
		return
	}
}
