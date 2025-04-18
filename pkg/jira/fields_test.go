package jira

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFieldsUnmarshalJSON(t *testing.T) {
	// Test with different date formats
	testCases := []struct {
		name     string
		jsonData string
		expected time.Time
	}{
		{
			name:     "Simple date format",
			jsonData: `{"summary": "Test Issue", "duedate": "2024-03-15"}`,
			expected: time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "ISO format with timezone",
			jsonData: `{"summary": "Test Issue", "duedate": "2024-03-15T14:30:00.000-0700"}`,
			expected: time.Date(2024, 3, 15, 14, 30, 0, 0, time.FixedZone("PDT", -7*60*60)),
		},
		{
			name:     "ISO format with UTC",
			jsonData: `{"summary": "Test Issue", "duedate": "2024-03-15T14:30:00.000Z"}`,
			expected: time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var fields Fields
			err := json.Unmarshal([]byte(tc.jsonData), &fields)
			assert.NoError(t, err)
			assert.NotNil(t, fields.DueDate)
			assert.True(t, fields.DueDate.Equal(tc.expected))
		})
	}
}

func TestFieldsWithMissingDueDate(t *testing.T) {
	jsonData := `{"summary": "Test Issue"}`

	var fields Fields
	err := json.Unmarshal([]byte(jsonData), &fields)
	assert.NoError(t, err)
	assert.Nil(t, fields.DueDate)
}

func TestFieldsIsOverdue(t *testing.T) {
	// Create a date in the past
	pastDate := time.Now().AddDate(0, 0, -1)  // Yesterday
	futureDate := time.Now().AddDate(0, 0, 1) // Tomorrow

	testCases := []struct {
		name     string
		dueDate  *time.Time
		expected bool
	}{
		{
			name:     "Overdue issue",
			dueDate:  &pastDate,
			expected: true,
		},
		{
			name:     "Not overdue issue",
			dueDate:  &futureDate,
			expected: false,
		},
		{
			name:     "No due date",
			dueDate:  nil,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := Fields{DueDate: tc.dueDate}
			assert.Equal(t, tc.expected, fields.IsOverdue())
		})
	}
}

func TestFieldsDaysUntilDue(t *testing.T) {
	// Create dates for testing
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	nextWeek := now.AddDate(0, 0, 7)

	testCases := []struct {
		name     string
		dueDate  *time.Time
		expected int
	}{
		{
			name:     "Overdue issue",
			dueDate:  &yesterday,
			expected: -1,
		},
		{
			name:     "Due tomorrow",
			dueDate:  &tomorrow,
			expected: 1,
		},
		{
			name:     "Due next week",
			dueDate:  &nextWeek,
			expected: 7,
		},
		{
			name:     "No due date",
			dueDate:  nil,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := Fields{DueDate: tc.dueDate}
			assert.Equal(t, tc.expected, fields.DaysUntilDue())
		})
	}
}
