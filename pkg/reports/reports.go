package reports

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// OutputFormat represents the type of report output
type OutputFormat string

const (
	TableFormat OutputFormat = "table"
	JSONFormat  OutputFormat = "json"
	CSVFormat   OutputFormat = "csv"
)

// GenerateReport generates a report in the specified format
func GenerateReport(data interface{}, format OutputFormat) error {
	// Handle nil data
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)

	// Handle non-slice data by wrapping it in a slice
	if v.Kind() != reflect.Slice {
		// Create a slice of the same type as data
		sliceType := reflect.SliceOf(v.Type())
		slice := reflect.MakeSlice(sliceType, 1, 1)
		slice.Index(0).Set(v)
		v = slice
	} else if v.Len() == 0 {
		// For JSON format, print empty array instead of null
		if format == JSONFormat {
			fmt.Println("[]")
			return nil
		}
		// For other formats, print a message
		fmt.Println("No data available.")
		return nil
	}

	// Continue with normal report generation
	switch format {
	case TableFormat:
		return generateTableReport(v)
	case JSONFormat:
		return generateJSONReport(data)
	case CSVFormat:
		return generateCSVReport(v)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// For JSON output, ensure numeric fields are properly marshaled as numbers
func WriteJSONReport(data interface{}, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
