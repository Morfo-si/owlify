// Assisted by watsonx Code Assistant
package reports

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateReport(t *testing.T) {
	type Data struct {
		Name string
		Age  int
	}

	data := []Data{
		{Name: "John", Age: 30},
		{Name: "Jane", Age: 25},
	}

	err := GenerateReport(data, TableFormat)
	assert.NoError(t, err)

	err = GenerateReport(data, JSONFormat)
	assert.NoError(t, err)

	err = GenerateReport(data, CSVFormat)
	assert.NoError(t, err)
}

func TestGetFlattenedHeaders(t *testing.T) {
	type TestData struct {
		Input    reflect.Type
		Expected []string
	}

	testCases := []TestData{
		{
			Input: reflect.TypeOf(struct {
				Name string `json:"name"`
			}{}),
			Expected: []string{"Name"},
		},
		{
			Input: reflect.TypeOf(struct {
				Name string `json:"name,omitempty"`
			}{}),
			Expected: []string{"Name"},
		},
		{
			Input: reflect.TypeOf(struct {
				Name string `json:"name,omitempty"`
				Age  int    `json:"age"`
			}{}),
			Expected: []string{"Name", "Age"},
		},
		{
			Input: reflect.TypeOf(struct {
				Name string `json:"name,omitempty"`
				Age  int    `json:"age"`
				Kid  struct {
					Name string `json:"name"`
					Age  string `json:"age"`
				} `json:"kid"`
			}{}),
			Expected: []string{"Name", "Age", "Kid.Name", "Kid.Age"},
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
			Expected: []string{"Name", "Age", "Name", "Age"},
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
			Expected: []string{"Name", "Age", "Pet.Breed", "Pet.Shot.Brand", "Pet.Shot.Date"},
		},
	}

	for _, tc := range testCases {
		actual := getFlattenedHeaders(tc.Input, "")
		assert.ElementsMatch(t, tc.Expected, actual)
	}
}

func TestGetFlattenedValues(t *testing.T) {
	type fields struct {
		Name string
		Age  int
	}
	type testStruct struct {
		fields
		Address string
	}

	tests := []struct {
		name   string
		fields interface{}
		want   []string
	}{
		{
			name:   "test 1",
			fields: fields{Name: "John", Age: 30},
			want:   []string{"John", "30"},
		},
		{
			name:   "test 2",
			fields: fields{Name: "Jane", Age: 25},
			want:   []string{"Jane", "25"},
		},
		{
			name: "test 3",
			fields: testStruct{
				fields:  fields{Name: "Alice", Age: 40},
				Address: "123 Main St",
			},
			want: []string{"Alice", "40", "123 Main St"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.fields)
			got := getFlattenedValues(v)
			assert.Equal(t, tt.want, got)
		})
	}
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
