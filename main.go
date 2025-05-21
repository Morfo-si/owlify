package main

import (
	"fmt"
	"os"

	"github.com/morfo-si/owlify/cmd"
	"github.com/morfo-si/owlify/pkg/config"
)

// Version information - will be set during build
var (
	version = "0.0.5"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Initialize configuration
	if err := config.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}
	
	// Pass version info to command package
	cmd.SetVersionInfo(version, commit, date)
	
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
