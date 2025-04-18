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

			if err := reports.GenerateReport(issue, reports.OutputFormat(output)); err != nil {
				fmt.Printf("Error generating report: %v\n", err)
				return
			}
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
	issueCmd.Flags().StringVarP(&issueKey, "key", "k", "", "JIRA issue key (required)")
	if err := issueCmd.MarkFlagRequired("key"); err != nil {
		fmt.Printf("Error marking key flag as required: %v\n", err)
		return
	}

	updateCmd.Flags().StringVarP(&issueKey, "key", "k", "", "JIRA issue key (required)")
	updateCmd.Flags().StringVarP(&newStatus, "status", "s", "", "New status (required)")
	if err := updateCmd.MarkFlagRequired("key"); err != nil {
		fmt.Printf("Error marking key flag as required: %v\n", err)
		return
	}
	if err := updateCmd.MarkFlagRequired("status"); err != nil {
		fmt.Printf("Error marking status flag as required: %v\n", err)
		return
	}

	issueCmd.AddCommand(updateCmd)
}
