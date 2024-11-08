package cmd

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
	"github.com/morfo-si/owlify/pkg/reports"
	"github.com/spf13/cobra"
)

var (
	issueKey  string
	newStatus string

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

			reports.GenerateReport(issue, reports.OutputFormat(output))
		},
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a JIRA issue fields",
		Run: func(cmd *cobra.Command, args []string) {
			if issueKey == "" {
				fmt.Println("Issue key is required")
				return
			}
			if newStatus == "" {
				fmt.Println("Status is required")
				return
			}

			err := jira.UpdateIssueStatus(issueKey, newStatus)
			if err != nil {
				fmt.Printf("Error updating issue %s status to %s: %v\n", issueKey, newStatus, err)
				return
			}

			fmt.Printf("Successfully updated issue %s status to %s\n", issueKey, newStatus)
		},
	}
)

func init() {
	issueCmd.PersistentFlags().StringVarP(&issueKey, "issue", "i", "", "JIRA issue key")
	updateCmd.Flags().StringVar(&newStatus, "status", "", "New status for the issue")
	issueCmd.AddCommand(updateCmd)
}
