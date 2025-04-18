// Assisted by watsonx Code Assistant

package jira

import (
	"errors"
	"testing"

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
