package reports

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

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
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	switch format {
	case TableFormat:
		table := tablewriter.NewWriter(os.Stdout)

		// Get flattened headers from struct fields
		if val.Len() > 0 {
			firstItem := val.Index(0)
			headers := getFlattenedHeaders(firstItem.Type(), "")
			table.SetHeader(headers)
		}

		// Add rows with flattened values
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			row := getFlattenedValues(item)
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

		// Write flattened headers
		if val.Len() > 0 {
			firstItem := val.Index(0)
			headers := getFlattenedHeaders(firstItem.Type(), "")
			if err := writer.Write(headers); err != nil {
				return fmt.Errorf("error writing CSV headers: %v", err)
			}
		}

		// Write rows with flattened values
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			row := getFlattenedValues(item)
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

// getFlattenedHeaders returns a slice of headers for all fields including nested structs
func getFlattenedHeaders(t reflect.Type, prefix string) []string {
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == "-" {
			continue
		}

		// Get the field name from JSON tag or use struct field name
		fieldName := field.Name
		if jsonTag != "" {
			// Split the json tag to get the name part (before any comma)
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		// Skip "fields" prefix for better readability
		if fieldName == "fields" {
			nestedHeaders := getFlattenedHeaders(field.Type, "")
			headers = append(headers, nestedHeaders...)
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			// Recursively get headers for nested struct
			nestedPrefix := prefix
			if prefix != "" {
				nestedPrefix = prefix + fieldName + "."
			} else {
				nestedPrefix = fieldName + "."
			}
			nestedHeaders := getFlattenedHeaders(field.Type, nestedPrefix)
			headers = append(headers, nestedHeaders...)
		} else {
			header := fieldName
			if prefix != "" {
				header = prefix + header
			}
			// Capitalize first letter for better readability
			if len(header) > 0 {
				header = strings.ToUpper(header[:1]) + header[1:]
			}
			headers = append(headers, header)
		}
	}
	return headers
}

// getFlattenedValues returns a slice of values for all fields including nested structs
func getFlattenedValues(v reflect.Value) []string {
	var values []string
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		jsonTag := t.Field(i).Tag.Get("json")

		if jsonTag == "-" {
			continue
		}

		fieldName := t.Field(i).Name
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		// Skip "fields" struct level
		if fieldName == "fields" {
			nestedValues := getFlattenedValues(field)
			values = append(values, nestedValues...)
			continue
		}

		if field.Kind() == reflect.Struct {
			// Recursively get values for nested struct
			nestedValues := getFlattenedValues(field)
			values = append(values, nestedValues...)
		} else {
			val := field.Interface()
			if field.Kind() == reflect.Ptr && !field.IsNil() {
				val = field.Elem().Interface()
			}
			values = append(values, fmt.Sprintf("%v", val))
		}
	}
	return values
}
