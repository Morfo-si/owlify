package main

import (
	"fmt"
	"os"

	"github.com/morfo-si/owlify/cmd"
)

// Version information - will be set during build
var (
	version = "0.0.5"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Pass version info to command package
	cmd.SetVersionInfo(version, commit, date)
	
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
