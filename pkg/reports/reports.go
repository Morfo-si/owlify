package reports

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

// OutputFormat represents the type of report output
type OutputFormat string

const (
	TableFormat OutputFormat = "table"
	JSONFormat  OutputFormat = "json"
	CSVFormat   OutputFormat = "csv"
)

// GenerateReport generates a report in the specified format for any slice of structs
func GenerateReport(data interface{}, format OutputFormat) error {
	// Get the value and verify it's a slice
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	switch format {
	case TableFormat:
		table := tablewriter.NewWriter(os.Stdout)

		// Get headers from struct fields
		if val.Len() > 0 {
			firstItem := val.Index(0)
			t := firstItem.Type()
			headers := make([]string, 0)

			for i := 0; i < t.NumField(); i++ {
				headers = append(headers, t.Field(i).Name)
			}
			table.SetHeader(headers)
		}

		// Add rows
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			row := make([]string, 0)

			for j := 0; j < item.NumField(); j++ {
				field := item.Field(j)
				row = append(row, fmt.Sprintf("%v", field.Interface()))
			}
			table.Append(row)
		}

		table.Render()

	case JSONFormat:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(data); err != nil {
			return fmt.Errorf("error encoding JSON: %v", err)
		}

	case CSVFormat:
		writer := csv.NewWriter(os.Stdout)

		// Write headers
		if val.Len() > 0 {
			firstItem := val.Index(0)
			t := firstItem.Type()
			headers := make([]string, 0)

			for i := 0; i < t.NumField(); i++ {
				headers = append(headers, t.Field(i).Name)
			}
			if err := writer.Write(headers); err != nil {
				return fmt.Errorf("error writing CSV headers: %v", err)
			}
		}

		// Write rows
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			row := make([]string, 0)

			for j := 0; j < item.NumField(); j++ {
				field := item.Field(j)
				row = append(row, fmt.Sprintf("%v", field.Interface()))
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("error writing CSV row: %v", err)
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			return fmt.Errorf("error flushing CSV writer: %v", err)
		}

	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	return nil
}
