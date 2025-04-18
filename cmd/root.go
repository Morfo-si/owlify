package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	component string
	project   string
	sprint    int
	output    string

	rootCmd = &cobra.Command{
		Use:   "owlify",
		Short: "A CLI tool to fetch JIRA issues",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&component, "component", "c", "", "JIRA component (optional)")
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", "", "JIRA project key (required)")
	rootCmd.PersistentFlags().IntVarP(&sprint, "sprint", "s", 0, "Sprint number (optional)")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format: table or json or csv")

	// Add sprint command to root
	rootCmd.AddCommand(sprintCmd)
	rootCmd.AddCommand(boardCmd)
	rootCmd.AddCommand(jqlCmd)
	rootCmd.AddCommand(issueCmd)

	if err := viper.BindPFlag("component", rootCmd.PersistentFlags().Lookup("component")); err != nil {
		fmt.Printf("Error binding component flag: %v\n", err)
	}
	if err := viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project")); err != nil {
		fmt.Printf("Error binding project flag: %v\n", err)
	}
	if err := viper.BindPFlag("sprint", rootCmd.PersistentFlags().Lookup("sprint")); err != nil {
		fmt.Printf("Error binding sprint flag: %v\n", err)
	}
	if err := viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output")); err != nil {
		fmt.Printf("Error binding output flag: %v\n", err)
	}
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
