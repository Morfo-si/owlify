package cmd

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
	"github.com/morfo-si/owlify/pkg/reports"
	"github.com/spf13/cobra"
)

var (
	jql    string
	jqlCmd = &cobra.Command{
		Use:   "search",
		Short: "Fetch JIRA issues using JQL",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jql == "" {
				return fmt.Errorf("jql is required")
			}
			issues, err := jira.FetchIssuesFromJQL(jql, jira.JIRAGetRequest)
			if err != nil {
				return fmt.Errorf("error fetching JIRA issues: %v", err)
			}
			reports.GenerateReport(issues, reports.OutputFormat(output))
			return nil
		},
	}
)

func init() {
	jqlCmd.Flags().StringVarP(&jql, "jql", "j", "", "JQL query (required)")
	jqlCmd.MarkFlagRequired("jql")
}
