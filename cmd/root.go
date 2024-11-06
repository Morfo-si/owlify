package cmd

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
	"github.com/morfo-si/owlify/pkg/reports"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	component string
	project   string
	sprint    int
	rootCmd   = &cobra.Command{
		Use:   "owlify",
		Short: "A CLI tool to fetch JIRA issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project == "" {
				return fmt.Errorf("project is required")
			}

			issues, err := jira.FetchCurrentSprintIssues(project, component, sprint)
			if err != nil {
				return fmt.Errorf("error fetching JIRA issues: %v", err)
			}
			if len(issues) > 0 {
				reports.GenerateReport(issues)
			} else {
				fmt.Println("No issues found for the specified criteria.")
			}
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&component, "component", "c", "", "JIRA component (optional)")
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", "", "JIRA project key (required)")
	rootCmd.PersistentFlags().IntVarP(&sprint, "sprint", "s", 0, "Sprint number (optional)")

	viper.BindPFlag("component", rootCmd.PersistentFlags().Lookup("component"))
	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	viper.BindPFlag("sprint", rootCmd.PersistentFlags().Lookup("sprint"))
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
