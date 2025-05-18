package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	output    string

	rootCmd = &cobra.Command{
		Use:   "owlify",
		Short: "A CLI tool to fetch JIRA issues",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format: table or json or csv")

	// Add commands to root command
	rootCmd.AddCommand(sprintCmd)
	rootCmd.AddCommand(boardCmd)
	rootCmd.AddCommand(jqlCmd)
	rootCmd.AddCommand(issueCmd)

	if err := viper.BindPFlag("component", rootCmd.PersistentFlags().Lookup("component")); err != nil {
		fmt.Printf("Error binding component flag: %v\n", err)
	}

	if err := viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output")); err != nil {
		fmt.Printf("Error binding output flag: %v\n", err)
	}
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
