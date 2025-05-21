/*
Package reports provides functionality for generating reports in various formats.

The package supports three output formats:
  - Table: Formatted ASCII tables for terminal output
  - JSON: Structured JSON for machine consumption
  - CSV: Comma-separated values for spreadsheet applications

The main entry point is the GenerateReport function, which takes any Go data
structure and converts it to the specified output format. The package uses
reflection to handle arbitrary struct types, including nested structures.

Special handling is provided for common types like time.Time, and for specific
struct fields like Epic and Feature which are expanded into separate columns.

Example usage:

    type Person struct {
        Name string
        Age  int
    }

    data := []Person{
        {Name: "John", Age: 30},
        {Name: "Jane", Age: 25},
    }

    // Generate a table report
    err := reports.GenerateReport(data, reports.TableFormat)
    if err != nil {
        log.Fatalf("Error generating report: %v", err)
    }

For testing or redirecting output, use the ReportWriter interface:

    var buf bytes.Buffer
    writer := &reports.StandardReportWriter{Out: &buf}
    err := writer.Write(data, reports.JSONFormat)
*/
package reports