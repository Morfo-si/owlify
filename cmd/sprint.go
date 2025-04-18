package cmd

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
	"github.com/morfo-si/owlify/pkg/reports"

	"github.com/spf13/cobra"
)

var (
	boardId int

	sprintCmd = &cobra.Command{
		Use:   "sprint",
		Short: "Fetch JIRA issues from sprints",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project == "" {
				return fmt.Errorf("project is required")
			}

			issues, err := jira.FetchCurrentSprintIssues(project, component, sprint, jira.JIRAGetRequest)
			if err != nil {
				return fmt.Errorf("error fetching JIRA issues: %v", err)
			}
			if len(issues) > 0 {
				reports.GenerateReport(issues, reports.OutputFormat(output))
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
				return fmt.Errorf("boardId is required")
			}

			sprints, err := jira.FetchOpenSprints(boardId, jira.JIRAGetRequest)
			if err != nil {
				return fmt.Errorf("error fetching sprints: %v", err)
			}

			reports.GenerateReport(sprints, reports.OutputFormat(output))
			return nil
		},
	}
)

func init() {
	// Add list command to sprint command
	sprintCmd.AddCommand(sprintListCmd)

	// Add required flags
	sprintListCmd.Flags().IntVarP(&boardId, "boardId", "b", 0, "JIRA board ID (required)")

	// Mark flags as required
	sprintListCmd.MarkFlagRequired("boardId")
}
