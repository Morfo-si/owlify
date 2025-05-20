package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	output    string
	showVersion bool

	// Version information
	versionInfo struct {
		Version string
		Commit  string
		Date    string
	}

	rootCmd = &cobra.Command{
		Use:   "owlify",
		Short: "A CLI tool to fetch JIRA issues",
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("Owlify version %s\n", versionInfo.Version)
				fmt.Printf("Commit: %s\n", versionInfo.Commit)
				fmt.Printf("Built: %s\n", versionInfo.Date)
				return
			}
			_ = cmd.Help()
		},
	}
)

// SetVersionInfo sets the version information
func SetVersionInfo(version, commit, date string) {
	versionInfo.Version = version
	versionInfo.Commit = commit
	versionInfo.Date = date
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format: table or json or csv")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version information")

	// Add commands to root command
	rootCmd.AddCommand(sprintCmd)
	rootCmd.AddCommand(boardCmd)
	rootCmd.AddCommand(jqlCmd)
	rootCmd.AddCommand(issueCmd)

	if err := viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output")); err != nil {
		fmt.Printf("Error binding output flag: %v\n", err)
	}
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
