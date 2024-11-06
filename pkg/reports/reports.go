package reports

import (
	"encoding/json"
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
)

// OutputFormat represents the type of report output
type OutputFormat string

const (
	TextFormat  OutputFormat = "text"
	JSONFormat  OutputFormat = "json"
)

// GenerateJSONReport generates a JSON report
func GenerateJSONReport(issues []jira.JiraIssue) {
	jsonData, err := json.MarshalIndent(issues, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

// GenerateReport generates a text report
func GenerateReport(issues []jira.JiraIssue) {
	fmt.Println("\nJIRA Issues Report")
	fmt.Println("==========================")
	for _, issue := range issues {
		fmt.Printf("Issue: %s\n", issue.Key)
		fmt.Printf("Summary: %s\n", issue.Fields.Summary)
		fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
		fmt.Printf("Due Date: %s\n", issue.Fields.DueDate)
		fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
		fmt.Println("--------------------------")
	}
}

// GenerateOutput creates a report in the specified format
func GenerateOutput(issues []jira.JiraIssue, format OutputFormat) {
	switch format {
	case JSONFormat:
		GenerateJSONReport(issues)
	default:
		GenerateReport(issues)
	}
}
