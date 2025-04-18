package jira

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock function for JIRAGetRequest
func mockJIRAGetRequest(url string, response interface{}) error {
	if url == "error" {
		return errors.New("failed to fetch issues")
	}

	// Create time.Time values for due dates using relative dates
	now := time.Now()
	dueDate1 := now.AddDate(0, 0, 7)  // 7 days in the future
	dueDate2 := now.AddDate(0, 0, 14) // 14 days in the future

	mockResponse := JiraResponse{
		Issues: []Issue{
			{
				Key: "ISSUE-1",
				Fields: Fields{
					Summary:    "First issue",
					DueDate:    &dueDate1,
					StoryPoint: 3.0,
					Priority: Priority{
						Name: "High",
					},
					Status: Status{
						Name: "In Progress",
					},
				},
			},
			{
				Key: "ISSUE-2",
				Fields: Fields{
					Summary:    "Second issue",
					DueDate:    &dueDate2,
					StoryPoint: 5.0,
					Priority: Priority{
						Name: "Medium",
					},
					Status: Status{
						Name: "To Do",
					},
				},
			},
		},
	}

	// Type assertion
	if r, ok := response.(*JiraResponse); ok {
		*r = mockResponse
		return nil
	}

	return errors.New("invalid response type")
}

func TestFetchIssuesFromJQL_Success(t *testing.T) {
	jql := "project=TEST"

	issues, err := FetchIssuesFromJQL(jql, mockJIRAGetRequest)
	assert.NoError(t, err)
	assert.Len(t, issues, 2)
	assert.Equal(t, "ISSUE-1", issues[0].Key)
	assert.Equal(t, "First issue", issues[0].Fields.Summary)

	// Test date handling
	assert.NotNil(t, issues[0].Fields.DueDate)
	now := time.Now()
	assert.True(t, issues[0].Fields.DueDate.After(now))

	// Test helper methods
	assert.False(t, issues[0].Fields.IsOverdue())
	assert.Greater(t, issues[0].Fields.DaysUntilDue(), 0)
}

func TestFetchIssuesFromJQL_APIError(t *testing.T) {
	jql := "error"

	// Mock function that returns an error
	mockErrorFunc := func(url string, response interface{}) error {
		return errors.New("API failure")
	}

	issues, err := FetchIssuesFromJQL(jql, mockErrorFunc)
	assert.Error(t, err)
	assert.Nil(t, issues)
}

func TestFetchIssuesFromJQL_WithOverdueIssue(t *testing.T) {
	// Create a mock function that returns an overdue issue
	mockOverdueFunc := func(url string, response interface{}) error {
		// Create a date in the past
		pastDate := time.Now().AddDate(0, 0, -1) // Yesterday

		mockResponse := JiraResponse{
			Issues: []Issue{
				{
					Key: "OVERDUE-1",
					Fields: Fields{
						Summary: "Overdue issue",
						DueDate: &pastDate,
					},
				},
			},
		}

		if r, ok := response.(*JiraResponse); ok {
			*r = mockResponse
			return nil
		}

		return errors.New("invalid response type")
	}

	jql := "project=TEST"
	issues, err := FetchIssuesFromJQL(jql, mockOverdueFunc)

	assert.NoError(t, err)
	assert.Len(t, issues, 1)
	assert.True(t, issues[0].Fields.IsOverdue())
	assert.Less(t, issues[0].Fields.DaysUntilDue(), 0)
}
