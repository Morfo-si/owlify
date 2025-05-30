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

			issue, err := jira.GetIssue(issueKey, jira.JIRAGetRequest)
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

	issueUpdateStatusCmd = &cobra.Command{
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

			// Fetch the issue to get the current status
			issue, err := jira.GetIssue(issueKey, jira.JIRAGetRequest)
			if err != nil {
				fmt.Printf("Error fetching issue %s: %v\n", issueKey, err)
				return
			}
			// Fetch the available status transitions
			transitions, err := jira.GetAvailableTransitions(issue, jira.JIRAGetRequest)
			if err != nil {
				fmt.Printf("Error fetching transitions for issue %s: %v\n", issueKey, err)
				return
			}

			// Validate transition name
			status := jira.GetValidTransitionID(newStatus, transitions)
			if status == "" {
				fmt.Printf("Invalid transition name: %s\n", newStatus)
				return
			}

			// Update the issue status
			err = jira.UpdateIssueStatus(issue.Key, status, jira.JIRAPostRequest)
			if err != nil {
				fmt.Printf("Error updating issue %s status to %s: %v\n", issueKey, newStatus, err)
				return
			}
			fmt.Printf("Successfully updated issue %s status to %s\n", issueKey, newStatus)
		},
	}
)

func init() {
	issueCmd.PersistentFlags().StringVarP(&issueKey, "key", "k", "", "JIRA issue key (required)")
	issueUpdateStatusCmd.PersistentFlags().StringVarP(&newStatus, "status", "s", "", "New status (required)")

	issueCmd.AddCommand(issueUpdateStatusCmd)
}
