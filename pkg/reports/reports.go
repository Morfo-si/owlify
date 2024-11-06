package reports

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/morfo-si/owlify/pkg/jira"
)

// OutputFormat represents the type of report output
type OutputFormat string

const (
	TextFormat OutputFormat = "text"
	JSONFormat OutputFormat = "json"
)

func GenerateTableReport(issues []jira.JiraIssue) {
	// Create headers slice starting with Key
	headers := []string{"Key", "Summary", "Status"}

	// Get all field names from the first issue (if available)
	if len(issues) > 0 {

		// Add remaining field names dynamically
		// Use reflection to get struct fields
		v := reflect.ValueOf(issues[0].Fields)
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			fieldName := t.Field(i).Name
			if fieldName != "Key" && fieldName != "Summary" && fieldName != "Status" {
				headers = append(headers, fieldName)
			}
		}
	}

	// Print table header
	fmt.Printf("%-12s", headers[0]) // Key
	fmt.Printf("%-28s", headers[1]) // Summary (starts at 13)
	fmt.Printf("%-15s", headers[2]) // Status (starts at 38)
	for i := 3; i < len(headers); i++ {
		fmt.Printf("\t%-15s", headers[i])
	}
	fmt.Println()

	// Print rows
	for _, issue := range issues {
		fmt.Printf("%-12s", issue.Key)
		
		summary := issue.Fields.Summary
		if len(summary) > 25 {
			summary = summary[:22] + "..."
		}
		fmt.Printf("%-28s", summary)
		
		fmt.Printf("%-15s", issue.Fields.Status.Name)
		fmt.Printf("%-15s", issue.Fields.Priority.Name)
		if issue.Fields.DueDate != "" {
			fmt.Printf("\t%-15s", issue.Fields.DueDate)
		}
		fmt.Println()
	}

}

// GenerateJSONReport generates a JSON report
func GenerateJSONReport(issues []jira.JiraIssue) {
	jsonData, err := json.MarshalIndent(issues, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

// GenerateOutput creates a report in the specified format
func GenerateOutput(issues []jira.JiraIssue, format OutputFormat) {
	switch format {
	case JSONFormat:
		GenerateJSONReport(issues)
	default:
		GenerateTableReport(issues)
	}
}
