package reports

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

// generateJSONReport creates a JSON-formatted report
func generateJSONReport(data any) error {
	// Handle nil data
	if data == nil {
		fmt.Println("[]")
		return nil
	}

	// Check if data is an empty slice
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.Len() == 0 {
		fmt.Println("[]")
		return nil
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return newReportError(JSONFormat, "encoding", err)
	}
	return nil
}
