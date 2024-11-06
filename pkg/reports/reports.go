package reports

import (
	"fmt"

	"github.com/morfo-si/owlify/pkg/jira"
)

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
