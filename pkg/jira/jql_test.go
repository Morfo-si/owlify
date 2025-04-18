package jira

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock function for JIRAGetRequest
func mockJIRAGetRequest(url string, response interface{}) error {
	if url == "error" {
		return errors.New("failed to fetch issues")
	}

	mockResponse := JiraResponse{
		Issues: []Issue{
			{
				Key: "ISSUE-1",
				Fields: Fields{
					Summary: "First issue",
				},
			},
			{
				Key: "ISSUE-2",
				Fields: Fields{
					Summary: "Second issue",
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
