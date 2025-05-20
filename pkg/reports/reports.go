package reports

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

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
func GenerateReport(data any, format OutputFormat) error {
	val := reflect.ValueOf(data)
	
	// Check if data is a slice, if not, wrap it in a slice
	if val.Kind() != reflect.Slice {
		// Create a slice of the same type as data
		sliceType := reflect.SliceOf(val.Type())
		slice := reflect.MakeSlice(sliceType, 1, 1)
		slice.Index(0).Set(val)
		val = slice
	}

	switch format {
	case TableFormat:
		table := tablewriter.NewWriter(os.Stdout)

		// Get flattened headers from struct fields
		if val.Len() > 0 {
			firstItem := val.Index(0)
			headers := getFlattenedHeaders(firstItem.Type(), "")
			table.Header(headers)
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

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
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
		if strings.EqualFold(fieldName, "fields") {
			nestedHeaders := getFlattenedHeaders(field.Type, "")
			headers = append(headers, nestedHeaders...)
			continue
		}

		// Always capitalize the field name
		fieldName = capitalizeFirst(fieldName)

		// Special handling for Epic field - create separate headers for its fields
		if fieldName == "Epic" && field.Type.Kind() == reflect.Ptr {
			// Get the Epic struct fields as separate headers
			structType := field.Type.Elem() // Get the type the pointer points to
			for j := 0; j < structType.NumField(); j++ {
				epicField := structType.Field(j)
				epicFieldName := epicField.Name
				jsonTag := epicField.Tag.Get("json")
				
				if jsonTag != "" {
					parts := strings.Split(jsonTag, ",")
					if parts[0] != "" {
						epicFieldName = parts[0]
					}
				}
				
				epicFieldName = capitalizeFirst(epicFieldName)
				if prefix != "" {
					headers = append(headers, prefix+"Epic."+epicFieldName)
				} else {
					headers = append(headers, "Epic."+epicFieldName)
				}
			}
			continue
		// Special handling for Feature field - create separate headers for its fields
		} else if fieldName == "Feature" && field.Type.Kind() == reflect.Ptr {
			// Get the Feature struct fields as separate headers
			structType := field.Type.Elem() // Get the type the pointer points to
			for j := 0; j < structType.NumField(); j++ {
				featureField := structType.Field(j)
				featureFieldName := featureField.Name
				jsonTag := featureField.Tag.Get("json")
				
				if jsonTag != "" {
					parts := strings.Split(jsonTag, ",")
					if parts[0] != "" {
						featureFieldName = parts[0]
					}
				}
				
				featureFieldName = capitalizeFirst(featureFieldName)
				if prefix != "" {
					headers = append(headers, prefix+"Feature."+featureFieldName)
				} else {
					headers = append(headers, "Feature."+featureFieldName)
				}
			}
			continue
		} else if field.Type.Kind() == reflect.Struct {
			// Recursively get headers for nested struct
			var nestedPrefix string
			if prefix != "" {
				nestedPrefix = prefix + fieldName + "."
			} else {
				nestedPrefix = fieldName + "."
			}
			nestedHeaders := getFlattenedHeaders(field.Type, nestedPrefix)
			headers = append(headers, nestedHeaders...)
		} else {
			// Add the field name with prefix
			if prefix != "" {
				headers = append(headers, prefix+fieldName)
			} else {
				headers = append(headers, fieldName)
			}
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

		// Handle Epic struct field by adding each of its fields separately
		if fieldName == "epic" && field.Kind() == reflect.Ptr && !field.IsNil() {
			epicStruct := field.Elem()
			
			// Instead of combining fields, get values for each field individually
			for j := 0; j < epicStruct.NumField(); j++ {
				// Get the value as string for each field
				fieldValue := fmt.Sprintf("%v", epicStruct.Field(j).Interface())
				values = append(values, fieldValue)
			}
			continue
		}
		
		// Handle Feature struct field by adding each of its fields separately
		if fieldName == "feature" && field.Kind() == reflect.Ptr && !field.IsNil() {
			featureStruct := field.Elem()
			
			// Instead of combining fields, get values for each field individually
			for j := 0; j < featureStruct.NumField(); j++ {
				// Get the value as string for each field
				fieldValue := fmt.Sprintf("%v", featureStruct.Field(j).Interface())
				values = append(values, fieldValue)
			}
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
			
			// Special formatting for time.Time values
			switch v := val.(type) {
			case time.Time:
				// Format time as YYYY-MM-DD
				values = append(values, v.Format("2006-01-02"))
			default:
				values = append(values, fmt.Sprintf("%v", val))
			}
		}
	}
	return values
}
