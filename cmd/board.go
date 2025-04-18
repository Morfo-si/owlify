package cmd

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
	"github.com/morfo-si/owlify/pkg/reports"

	"github.com/spf13/cobra"
)

var (
	boardCmd = &cobra.Command{
		Use:   "board",
		Short: "Fetch JIRA boards from project",
	}

	boardListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all boards from project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project == "" {
				return fmt.Errorf("project is required")
			}
			boards, err := jira.FetchBoards(project)
			if err != nil {
				return fmt.Errorf("error fetching boards: %v", err)
			}
			if err := reports.GenerateReport(boards, reports.OutputFormat(output)); err != nil {
				return fmt.Errorf("error generating report: %v", err)
			}
			return nil
		},
	}
)

func init() {
	boardCmd.AddCommand(boardListCmd)
}
