package reports

import (
    "fmt"
)

// ReportError represents an error that occurred during report generation
type ReportError struct {
    Format OutputFormat
    Op     string
    Err    error
}

// Error implements the error interface
func (e *ReportError) Error() string {
    return fmt.Sprintf("report error (%s format, %s operation): %v", e.Format, e.Op, e.Err)
}

// Unwrap returns the underlying error
func (e *ReportError) Unwrap() error {
    return e.Err
}

// newReportError creates a new ReportError
func newReportError(format OutputFormat, op string, err error) *ReportError {
    return &ReportError{
        Format: format,
        Op:     op,
        Err:    err,
    }
}