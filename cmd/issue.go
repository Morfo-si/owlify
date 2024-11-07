package cmd

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
	"github.com/morfo-si/owlify/pkg/reports"
	"github.com/spf13/cobra"
)

var (
	issueKey string

	issueCmd = &cobra.Command{
		Use:   "issue",
		Short: "Fetch a JIRA issue",
		Run: func(cmd *cobra.Command, args []string) {
			if issueKey == "" {
				fmt.Println("Issue key is required")
				return
			}

			issue, err := jira.GetIssue(issueKey)
			if err != nil {
				fmt.Printf("Error fetching issue %s: %v\n", issueKey, err)
				return
			}

			fmt.Printf("Issue: %+v\n", issue)
			reports.GenerateReport(issue, reports.OutputFormat(output))
		},
	}
)

func init() {
	issueCmd.Flags().StringVarP(&issueKey, "issue", "i", "", "JIRA issue key")
}
