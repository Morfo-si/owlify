// Assisted by watsonx Code Assistant

package jira

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock function for JiraRequestFunc for FetchOpenSprints call
func mockOpenSprintsRequest(url string, response interface{}) error {
	if url == "error" {
		return errors.New("failed to fetch sprints")
	}

	mockResponse := SprintResponse{
		Values: []Sprint{
			{ID: 1, Name: "Sprint 1", State: "active"},
			{ID: 2, Name: "Sprint 2", State: "active"},
		},
	}

	// Type assertion to correctly populate response
	if r, ok := response.(*SprintResponse); ok {
		*r = mockResponse
		return nil
	}

	return errors.New("invalid response type")
}

func mockOpenSprintIssuesRequest(url string, response interface{}) error {
	if url == "error" {
		return errors.New("failed to fetch sprints")
	}

	mockResponse := JiraResponse{
		[]Issue{
			{
				Key: "issue1",
				Fields: Fields{
					Summary: "Issue 1",
					Assignee: Assignee{
						Name: "John Doe",
					},
					StoryPoint: 3.0,
					DueDate:    "2025-03-15",
					Priority: Priority{
						Name: "High",
					},
					Status: Status{
						Name: "In Progress",
					},
				},
			},
			{
				Key: "issue2",
				Fields: Fields{
					Summary: "Issue 2",
					Assignee: Assignee{
						Name: "Jane Smith",
					},
					StoryPoint: 5.0,
					DueDate:    "2025-03-20",
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

	// Type assertion to correctly populate response
	if r, ok := response.(*JiraResponse); ok {
		*r = mockResponse
		return nil
	}

	return errors.New("invalid response type")
}

func TestFetchOpenSprints_Success(t *testing.T) {
	boardID := 123

	sprints, err := FetchOpenSprints(boardID, mockOpenSprintsRequest)
	assert.NoError(t, err)
	assert.Len(t, sprints, 2)
	assert.Equal(t, "Sprint 1", sprints[0].Name)
	assert.Equal(t, "Sprint 2", sprints[1].Name)
}

func TestFetchOpenSprints_APIError(t *testing.T) {
	boardID := 123

	// Mock function to simulate an API error
	mockErrorFunc := func(url string, response interface{}) error {
		return errors.New("API failure")
	}

	sprints, err := FetchOpenSprints(boardID, mockErrorFunc)
	assert.Error(t, err)
	assert.Nil(t, sprints)
}

func TestFetchOpenSprints_EmptyResponse(t *testing.T) {
	boardID := 123

	// Mock function returning an empty sprint list
	mockEmptyFunc := func(url string, response interface{}) error {
		if r, ok := response.(*SprintResponse); ok {
			*r = SprintResponse{Values: []Sprint{}}
			return nil
		}
		return errors.New("invalid response type")
	}

	sprints, err := FetchOpenSprints(boardID, mockEmptyFunc)
	assert.NoError(t, err)
	assert.Empty(t, sprints)
}

func TestFetchCurrentSprintIssues_Success(t *testing.T) {
	project := "project1"
	component := "component1"
	sprintNumber := 1

	issues, err := FetchCurrentSprintIssues(project, component, sprintNumber, mockOpenSprintIssuesRequest)
	assert.NoError(t, err)
	assert.Len(t, issues, 2)
	assert.Equal(t, "issue1", issues[0].Key)
	assert.Equal(t, "issue2", issues[1].Key)
}

func TestSprintUnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": 123,
		"self": "https://jira.example.com/rest/agile/1.0/sprint/123",
		"state": "active",
		"name": "Sprint 1",
		"startDate": "2024-03-01T00:00:00.000-0700",
		"endDate": "2024-03-15T00:00:00.000-0700",
		"activatedDate": "2024-03-01T09:00:00.000-0700",
		"originBoardId": 456,
		"synced": true,
		"autoStartStop": false
	}`

	var sprint Sprint
	if err := json.Unmarshal([]byte(jsonData), &sprint); err != nil {
		t.Fatalf("Failed to unmarshal Sprint: %v", err)
	}

	// Verify fields
	if sprint.ID != 123 {
		t.Errorf("Expected ID 123, got %d", sprint.ID)
	}

	if sprint.State != SprintStateActive {
		t.Errorf("Expected state %s, got %s", SprintStateActive, sprint.State)
	}

	// Verify dates
	expectedStart, _ := time.Parse("2006-01-02T15:04:05.000-0700", "2024-03-01T00:00:00.000-0700")
	if !sprint.StartDate.Equal(expectedStart) {
		t.Errorf("Expected start date %v, got %v", expectedStart, sprint.StartDate)
	}

	// Test helper methods
	if !sprint.IsActive() {
		t.Error("Expected sprint to be active")
	}

	duration := sprint.Duration()
	expectedDuration := 14 * 24 * time.Hour // 14 days
	if duration != expectedDuration {
		t.Errorf("Expected duration %v, got %v", expectedDuration, duration)
	}
}

func TestSprintWithMissingDates(t *testing.T) {
	jsonData := `{
		"id": 123,
		"self": "https://jira.example.com/rest/agile/1.0/sprint/123",
		"state": "future",
		"name": "Future Sprint",
		"originBoardId": 456,
		"synced": true,
		"autoStartStop": false
	}`

	var sprint Sprint
	if err := json.Unmarshal([]byte(jsonData), &sprint); err != nil {
		t.Fatalf("Failed to unmarshal Sprint: %v", err)
	}

	// Verify dates are nil
	if sprint.StartDate != nil {
		t.Error("Expected StartDate to be nil")
	}
	if sprint.EndDate != nil {
		t.Error("Expected EndDate to be nil")
	}
	if sprint.ActivatedDate != nil {
		t.Error("Expected ActivatedDate to be nil")
	}

	// Duration should be 0 for sprints without dates
	if sprint.Duration() != 0 {
		t.Errorf("Expected duration 0, got %v", sprint.Duration())
	}
}
