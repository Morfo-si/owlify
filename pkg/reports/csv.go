package reports

import (
	"encoding/csv"
	"os"
	"reflect"
)

// generateCSVReport creates a CSV-formatted report
func generateCSVReport(val reflect.Value) error {
	writer := csv.NewWriter(os.Stdout)

	// Write flattened headers
	if val.Len() > 0 {
		firstItem := val.Index(0)
		headers := GetFlattenedHeaders(firstItem.Interface())
		if err := writer.Write(headers); err != nil {
			return newReportError(CSVFormat, "writing headers", err)
		}
	}

	// Write rows with flattened values
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		row := getFlattenedValues(item)
		if err := writer.Write(row); err != nil {
			return newReportError(CSVFormat, "writing row", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return newReportError(CSVFormat, "flushing writer", err)
	}

	return nil
}
