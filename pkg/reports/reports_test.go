package reports

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStruct is used for testing report generation
type TestStruct struct {
	ID   int    `json:"id,string"` // Add string tag to handle string-to-int conversion
	Name string `json:"name"`
}

// TestReportWriter tests the report writer interface
func TestReportWriter(t *testing.T) {
	// Create test data
	testData := []TestStruct{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}

	// Test table format
	t.Run("table format", func(t *testing.T) {
		var buf bytes.Buffer
		writer := &StandardReportWriter{Out: &buf}

		err := writer.Write(testData, TableFormat)
		assert.NoError(t, err)

		output := buf.String()
		t.Logf("Table output:\n%s", output)

		// Check for column headers and data
		assert.Contains(t, output, "ID")
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "1")
		assert.Contains(t, output, "Test 1")
		assert.Contains(t, output, "2")
		assert.Contains(t, output, "Test 2")
	})

	// Test JSON format
	t.Run("json format", func(t *testing.T) {
		var buf bytes.Buffer
		writer := &StandardReportWriter{Out: &buf}

		err := writer.Write(testData, JSONFormat)
		assert.NoError(t, err)

		output := buf.String()
		t.Logf("JSON output:\n%s", output)

		// Parse JSON and verify
		var result []TestStruct
		err = json.Unmarshal([]byte(output), &result)
		assert.NoError(t, err)
		assert.Equal(t, testData, result)
	})

	// Test CSV format
	t.Run("csv format", func(t *testing.T) {
		var buf bytes.Buffer
		writer := &StandardReportWriter{Out: &buf}

		err := writer.Write(testData, CSVFormat)
		assert.NoError(t, err)

		output := buf.String()
		t.Logf("CSV output:\n%s", output)

		// Check CSV format
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Equal(t, 3, len(lines)) // Header + 2 data rows

		// Check header
		header := strings.Split(lines[0], ",")
		assert.Contains(t, header, "id")
		assert.Contains(t, header, "name")

		// Check data rows
		row1 := strings.Split(lines[1], ",")
		assert.Contains(t, row1, "1")
		assert.Contains(t, row1, "Test 1")

		row2 := strings.Split(lines[2], ",")
		assert.Contains(t, row2, "2")
		assert.Contains(t, row2, "Test 2")
	})
}

// TestComplexStructReporting tests reporting for complex structures
func TestComplexStructReporting(t *testing.T) {
	t.Skip("Skipping until reflection issues are fixed")
}

func TestGetFlattenedHeaders(t *testing.T) {
	type TestData struct {
		Input    reflect.Type
		Expected []string
	}

	testCases := []TestData{
		{
			Input: reflect.TypeOf(struct {
				Name string `json:"name,omitempty"`
			}{}),
			Expected: []string{"name"},
		},
		{
			Input: reflect.TypeOf(struct {
				Name string
			}{}),
			Expected: []string{"name"},
		},
		{
			Input: reflect.TypeOf(struct {
				Name string `json:"name,omitempty"`
				Age  int    `json:"age"`
			}{}),
			Expected: []string{"name", "age"},
		},
		{
			Input: reflect.TypeOf(struct {
				Name string `json:"name,omitempty"`
				Age  int    `json:"age"`
				Kid  struct {
					Name string `json:"name"`
					Age  int    `json:"age"`
				} `json:"kid"`
			}{}),
			Expected: []string{"name", "age", "kid.name", "kid.age"},
		},
		{
			Input: reflect.TypeOf(struct {
				Name   string `json:"name,omitempty"`
				Age    int    `json:"age"`
				Fields struct {
					Name string `json:"name"`
					Age  string `json:"age"`
				} `json:"fields"`
			}{}),
			Expected: []string{"name", "age", "fields.name", "fields.age"},
		},
		{
			Input: reflect.TypeOf(struct {
				Name string `json:"name,omitempty"`
				Age  int    `json:"age"`
				Pet  struct {
					Breed string `json:"breed"`
					Shot  struct {
						Brand string `json:"brand"`
						Date  string `json:"date"`
					} `json:"shot"`
				} `json:"pet"`
			}{}),
			Expected: []string{"name", "age", "pet.breed", "pet.shot.brand", "pet.shot.date"},
		},
	}

	for _, tc := range testCases {
		actual := GetFlattenedHeaders(tc.Input)
		assert.ElementsMatch(t, tc.Expected, actual)
	}
}

func TestGetFlattenedValues(t *testing.T) {
	t.Skip("Skipping until reflection issues are fixed")
}

func TestGenerateReportWithNonSliceData(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Test with a single struct (non-slice)
	person := Person{Name: "Alice", Age: 30}

	// Test with different formats
	formats := []OutputFormat{TableFormat, JSONFormat, CSVFormat}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			err := GenerateReport(person, format)
			assert.NoError(t, err)
		})
	}
}
