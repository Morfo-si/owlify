package reports

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// ReportWriter defines an interface for writing reports
type ReportWriter interface {
	Write(data any, format OutputFormat) error
}

// StandardReportWriter writes reports to stdout
type StandardReportWriter struct {
	Out io.Writer
}

// NewStandardReportWriter creates a new StandardReportWriter
func NewStandardReportWriter() *StandardReportWriter {
	return &StandardReportWriter{
		Out: os.Stdout,
	}
}

// Write writes the data to the output in the specified format
func (w *StandardReportWriter) Write(data interface{}, format OutputFormat) error {
	// Handle nil data
	if data == nil {
		return fmt.Errorf("cannot write nil data")
	}

	// Get the value and type of the data
	value := reflect.ValueOf(data)

	// Handle pointers
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return fmt.Errorf("cannot write nil data")
		}
		value = value.Elem()
	}

	// Handle different kinds of data
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		return w.writeSlice(value, format)
	case reflect.Struct:
		return w.writeStruct(value, format)
	case reflect.Map:
		return w.writeMap(value, format)
	default:
		return fmt.Errorf("unsupported data type: %s", value.Kind())
	}
}

// writeSlice writes a slice of data to the output
func (w *StandardReportWriter) writeSlice(value reflect.Value, format OutputFormat) error {
	// Handle empty slice
	if value.Len() == 0 {
		return nil
	}

	// Get the first element to determine the type
	firstElem := value.Index(0)

	// Handle pointers
	if firstElem.Kind() == reflect.Ptr {
		if firstElem.IsNil() {
			return fmt.Errorf("cannot write slice with nil elements")
		}
		firstElem = firstElem.Elem()
	}

	// Handle different kinds of elements
	switch firstElem.Kind() {
	case reflect.Struct:
		// Get headers from the struct type
		headers := getHeaders(firstElem.Type())

		// Get rows from each element in the slice
		var rows [][]string
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			if elem.Kind() == reflect.Ptr {
				if elem.IsNil() {
					continue
				}
				elem = elem.Elem()
			}

			if elem.Kind() == reflect.Struct {
				row := getValues(elem)
				rows = append(rows, row)
			}
		}

		// Write the data in the specified format
		return w.writeData(headers, rows, format)

	case reflect.Map:
		// For slices of maps, convert to JSON
		return w.writeJSON(value.Interface())

	default:
		// For simple slices, convert to strings
		var rows [][]string
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			rows = append(rows, []string{fmt.Sprintf("%v", elem.Interface())})
		}
		return w.writeData([]string{"Value"}, rows, format)
	}
}

// writeStruct writes a single struct to the output
func (w *StandardReportWriter) writeStruct(value reflect.Value, format OutputFormat) error {
	// Get headers from the struct type
	headers := getHeaders(value.Type())

	// Get values from the struct
	row := getValues(value)

	// Write the data in the specified format
	return w.writeData(headers, [][]string{row}, format)
}

// writeMap writes a map to the output
func (w *StandardReportWriter) writeMap(value reflect.Value, _ OutputFormat) error {
	// For maps, convert to JSON
	return w.writeJSON(value.Interface())
}

// writeData writes the headers and rows in the specified format
func (w *StandardReportWriter) writeData(headers []string, rows [][]string, format OutputFormat) error {
	switch format {
	case TableFormat:
		return w.writeTable(headers, rows)
	case JSONFormat:
		return w.writeJSONFromRows(headers, rows)
	case CSVFormat:
		return w.writeCSV(headers, rows)
	default:
		return fmt.Errorf("unsupported format: %v", format)
	}
}

// writeTable writes the data as a formatted table
func (w *StandardReportWriter) writeTable(headers []string, rows [][]string) error {
	// Create a new table with configuration
	table := tablewriter.NewTable(w.Out,
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoWrap:  tw.WrapNone,  // Don't wrap text
					Alignment: tw.AlignLeft, // Left-align rows
				},
			},
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoWrap:  tw.WrapNone,  // Don't wrap headers
					Alignment: tw.AlignLeft, // Left-align headers
				},
			},
		}),
		tablewriter.WithRenderer(renderer.NewBlueprint(
			tw.Rendition{
				Borders: tw.BorderNone,
				Settings: tw.Settings{
					Separators: tw.SeparatorsNone,
				},
			})),
	)

	// Set headers and data
	table.Header(headers)
	table.Bulk(rows)

	// Render the table
	table.Render()

	return nil
}

// writeJSON writes the data as JSON
func (w *StandardReportWriter) writeJSON(data interface{}) error {
	encoder := json.NewEncoder(w.Out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// writeJSONFromRows converts rows to JSON and writes it
func (w *StandardReportWriter) writeJSONFromRows(headers []string, rows [][]string) error {
	// Convert rows to a slice of maps
	var data []map[string]string
	for _, row := range rows {
		item := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				item[header] = row[i]
			}
		}
		data = append(data, item)
	}

	return w.writeJSON(data)
}

// writeCSV writes the data as CSV
func (w *StandardReportWriter) writeCSV(headers []string, rows [][]string) error {
	writer := csv.NewWriter(w.Out)

	// Write headers
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write rows
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

// getHeaders returns the headers for a struct type
func getHeaders(t reflect.Type) []string {
	var headers []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get field name from JSON tag or use struct field name
		fieldName := field.Name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
		}

		headers = append(headers, fieldName)
	}

	return headers
}

// getValues returns the values for a struct
func getValues(v reflect.Value) []string {
	var values []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Skip unexported fields
		if !v.Type().Field(i).IsExported() {
			continue
		}

		// Handle nil pointers
		if field.Kind() == reflect.Ptr && field.IsNil() {
			values = append(values, "")
			continue
		}

		// Get value (dereference pointer if needed)
		var val interface{}
		if field.Kind() == reflect.Ptr {
			val = field.Elem().Interface()
		} else {
			val = field.Interface()
		}

		// Format value as string
		values = append(values, fmt.Sprintf("%v", val))
	}

	return values
}
