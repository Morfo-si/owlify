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
	state    string
	features bool

	sprintCmd = &cobra.Command{
		Use:   "sprint",
		Short: "Fetch JIRA issues from sprints",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sprintId == 0 {
				return fmt.Errorf("sprint id is required")
			}

			issues, err := jira.FetchSprintIssues(sprintId, jira.JIRAGetRequest, features)
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

			// Convert state to string
			if state == "a" {
				state = jira.SprintStateActive.String()
			} else if state == "c" {
				state = jira.SprintStateClosed.String()
			} else if state == "f" {
				state = jira.SprintStateFuture.String()
			} else {
				return fmt.Errorf("invalid sprint state: %s", state)
			}

			sprints, err := jira.FetchSprints(
				boardId,
				jira.JIRAGetRequest,
				jira.WithSprintState(jira.SprintState(state)))
			if err != nil {
				return fmt.Errorf("error fetching sprints: %v", err)
			}

			if err := reports.GenerateReport(sprints, reports.OutputFormat(output)); err != nil {
				return fmt.Errorf("error generating report: %v", err)
			}
			return nil
		},
	}

	sprintGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information about a specific sprint by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sprintId == 0 {
				return fmt.Errorf("sprint id is required")
			}

			sprint, err := jira.FetchSprintByID(sprintId, jira.JIRAGetRequest)
			if err != nil {
				return fmt.Errorf("error fetching sprint: %v", err)
			}

			// Wrap the single Sprint in a slice for the report generator
			sprintSlice := []jira.Sprint{sprint}
			if err := reports.GenerateReport(sprintSlice, reports.OutputFormat(output)); err != nil {
				return fmt.Errorf("error generating report: %v", err)
			}
			return nil
		},
	}
)

func init() {

	// Add required flags
	sprintCmd.Flags().IntVarP(&sprintId, "id", "i", 0, "JIRA sprint ID (required)")
	sprintListCmd.Flags().IntVarP(&boardId, "board", "b", 0, "JIRA board ID (required)")
	sprintListCmd.Flags().StringVarP(&state, "state", "s", "active", "Sprint state (a/active, c/closed, f/future)")
	sprintGetCmd.Flags().IntVarP(&sprintId, "id", "i", 0, "JIRA sprint ID (required)")

	// Add the fetch-features flag
	sprintCmd.Flags().BoolVar(&features, "features", false, "Also fetch Feature data for epics (default: false)")

	// Add subcommands to sprint command
	sprintCmd.AddCommand(sprintListCmd, sprintGetCmd)
}
