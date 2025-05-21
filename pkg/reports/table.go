package reports

import (
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

// generateTableReport creates a table-formatted report
func generateTableReport(val reflect.Value) error {
	table := tablewriter.NewWriter(os.Stdout)

	// Get flattened headers from struct fields
	if val.Len() > 0 {
		firstItem := val.Index(0)
		headers := GetFlattenedHeaders(firstItem.Interface())
		table.Header(headers)
	}

	// Add rows with flattened values
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		row := getFlattenedValues(item)
		table.Append(row)
	}

	table.Render()
	return nil
}
