package cmd

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
	"github.com/morfo-si/owlify/pkg/reports"

	"github.com/spf13/cobra"
)

var (
	boardId  int
	sprintId int

	sprintCmd = &cobra.Command{
		Use:   "sprint",
		Short: "Fetch JIRA issues from sprints",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sprintId == 0 {
				return fmt.Errorf("sprint id is required")
			}

			issues, err := jira.FetchSprintIssues(sprintId, jira.JIRAGetRequest)
			if err != nil {
				return fmt.Errorf("error fetching JIRA issues: %v", err)
			}
			if len(issues) > 0 {
				if err := reports.GenerateReport(issues, reports.OutputFormat(output)); err != nil {
					return fmt.Errorf("error generating report: %v", err)
				}
			} else {
				fmt.Println("No issues found for the specified criteria.")
			}
			return nil
		},
	}

	sprintListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all open sprints",
		RunE: func(cmd *cobra.Command, args []string) error {
			if boardId == 0 {
				return fmt.Errorf("board id is required")
			}

			sprints, err := jira.FetchOpenSprints(boardId, jira.JIRAGetRequest)
			if err != nil {
				return fmt.Errorf("error fetching sprints: %v", err)
			}

			if err := reports.GenerateReport(sprints, reports.OutputFormat(output)); err != nil {
				return fmt.Errorf("error generating report: %v", err)
			}
			return nil
		},
	}

	sprintIssuesCmd = &cobra.Command{
		Use:   "issues",
		Short: "List issues from a sprint with epic information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sprintId == 0 {
				return fmt.Errorf("sprint id is required")
			}

			issues, err := jira.FetchSprintIssuesWithEpic(sprintId, jira.JIRAGetRequest)
			if err != nil {
				return fmt.Errorf("error fetching JIRA issues with epic information: %v", err)
			}

			if len(issues) > 0 {
				if err := reports.GenerateReport(issues, reports.OutputFormat(output)); err != nil {
					return fmt.Errorf("error generating report: %v", err)
				}
			} else {
				fmt.Println("No issues found for the specified sprint.")
			}
			return nil
		},
	}
)

func init() {

	// Add required flags
	sprintCmd.Flags().IntVarP(&sprintId, "sprint", "s", 0, "JIRA sprint ID (required)")
	sprintListCmd.Flags().IntVarP(&boardId, "board", "b", 0, "JIRA board ID (required)")
	sprintIssuesCmd.Flags().IntVarP(&sprintId, "sprint", "s", 0, "JIRA sprint ID (required)")

	// Add subcommands to sprint command
	sprintCmd.AddCommand(sprintListCmd)
	sprintCmd.AddCommand(sprintIssuesCmd)
}
