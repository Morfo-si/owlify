package reports

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/morfo-si/owlify/pkg/jira"
)

// OutputFormat represents the type of report output
type OutputFormat string

const (
	TableFormat OutputFormat = "table"
	JSONFormat  OutputFormat = "json"
	CSVFormat   OutputFormat = "csv"
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

func GenerateCSVReport(issues []jira.JiraIssue) {
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("owlify-%s.csv", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"Key", "Summary", "Status", "Priority", "Due Date"}
	if err := writer.Write(headers); err != nil {
		fmt.Printf("Error writing CSV headers: %v\n", err)
		return
	}

	// Write data rows
	for _, issue := range issues {
		row := []string{
			issue.Key,
			issue.Fields.Summary,
			issue.Fields.Status.Name,
			issue.Fields.Priority.Name,
			issue.Fields.DueDate,
		}
		if err := writer.Write(row); err != nil {
			fmt.Printf("Error writing CSV row: %v\n", err)
			return
		}
	}

	fmt.Printf("CSV report generated: %s\n", filename)

}

// GenerateOutput creates a report in the specified format
func GenerateOutput(issues []jira.JiraIssue, format OutputFormat) {
	switch format {
	case JSONFormat:
		GenerateJSONReport(issues)
	case CSVFormat:
		GenerateCSVReport(issues)
	default:
		GenerateTableReport(issues)
	}
}
