package cmd

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
	"github.com/morfo-si/owlify/pkg/reports"

	"github.com/spf13/cobra"
)

var (
	project string
	
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
			boards, err := jira.FetchBoards(project, jira.JIRAGetRequest)
			if err != nil {
				return fmt.Errorf("error fetching boards: %v", err)
			}
			
			// Handle empty boards result
			if len(boards) == 0 {
				if output == "json" {
					fmt.Println("[]")
					return nil
				}
				fmt.Printf("No boards found for project %s\n", project)
				return nil
			}
			
			if err := reports.GenerateReport(boards, reports.OutputFormat(output)); err != nil {
				return fmt.Errorf("error generating report: %v", err)
			}
			return nil
		},
	}
)

func init() {
	boardListCmd.PersistentFlags().StringVarP(&project, "project", "p", "", "JIRA project key (required)")
	boardCmd.AddCommand(boardListCmd)
}
