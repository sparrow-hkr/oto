/*
Copyright © 2025 Chandra HERE chandra@gmeil.com
*/
package main

import (
	"github.com/sparrow-hkr/oto/cmd"
	"github.com/sparrow-hkr/oto/internal/banner"
)

func main() {
	banner.PrintBanner()
	cmd.Execute()
}
