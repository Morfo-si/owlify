package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIssue(t *testing.T) {
	tests := []struct {
		name           string
		issueKey       string
		mockResponse   Issue
		mockError      error
		expectedIssue  Issue
		expectedErrMsg string
	}{
		{
			name:     "successful fetch",
			issueKey: "TEST-123",
			mockResponse: Issue{
				Key: "TEST-123",
				Fields: Fields{
					Summary: "Test issue",
				},
			},
			mockError:     nil,
			expectedIssue: Issue{Key: "TEST-123", Fields: Fields{Summary: "Test issue"}},
		},
		{
			name:           "api error",
			issueKey:       "INVALID-456",
			mockResponse:   Issue{},
			mockError:      errors.New("API connection error"),
			expectedErrMsg: "error fetching issue INVALID-456: API connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request function
			mockRequest := func(url string, target interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}

				// Check URL format
				expectedURL := fmt.Sprintf("%s/rest/api/2/issue/%s", jiraBaseURL, tt.issueKey)
				assert.Equal(t, expectedURL, url)

				// Marshal and unmarshal to simulate JSON response
				data, _ := json.Marshal(tt.mockResponse)
				return json.Unmarshal(data, target)
			}

			// Call the function
			issue, err := GetIssue(tt.issueKey, mockRequest)

			// Check error
			if tt.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrMsg, err.Error())
				assert.Empty(t, issue)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIssue.Key, issue.Key)
				assert.Equal(t, tt.expectedIssue.Fields.Summary, issue.Fields.Summary)
			}
		})
	}
}

func TestGetValidTransitionID(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		transitions []Transition
		expected    string
	}{
		{
			name:   "exact match",
			status: "In Progress",
			transitions: []Transition{
				{ID: "11", Name: "To Do"},
				{ID: "21", Name: "In Progress"},
				{ID: "31", Name: "Done"},
			},
			expected: "21",
		},
		{
			name:   "case insensitive match",
			status: "in progress",
			transitions: []Transition{
				{ID: "11", Name: "To Do"},
				{ID: "21", Name: "In Progress"},
				{ID: "31", Name: "Done"},
			},
			expected: "21",
		},
		{
			name:   "no match",
			status: "Blocked",
			transitions: []Transition{
				{ID: "11", Name: "To Do"},
				{ID: "21", Name: "In Progress"},
				{ID: "31", Name: "Done"},
			},
			expected: "",
		},
		{
			name:        "empty transitions",
			status:      "In Progress",
			transitions: []Transition{},
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetValidTransitionID(tt.status, tt.transitions)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUpdateIssueStatus(t *testing.T) {
	tests := []struct {
		name           string
		issueKey       string
		newStatus      string
		mockPostError  error
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "successful status update",
			issueKey:      "TEST-123",
			newStatus:     "10001", // Transition ID
			mockPostError: nil,
			expectedError: false,
		},
		{
			name:           "API error",
			issueKey:       "TEST-456",
			newStatus:      "10002",
			mockPostError:  fmt.Errorf("API connection error"),
			expectedError:  true,
			expectedErrMsg: "error transitioning issue: API connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock post request function
			mockPostRequest := func(url string, payload interface{}, response interface{}) error {
				// Verify URL format
				expectedURL := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", jiraBaseURL, tt.issueKey)
				if url != expectedURL {
					t.Errorf("incorrect URL: got %s, want %s", url, expectedURL)
				}

				// Verify payload
				transition, ok := payload.(UpdateTransition)
				if !ok {
					t.Errorf("incorrect payload type: got %T, want UpdateTransition", payload)
				}
				if transition.Transition.ID != tt.newStatus {
					t.Errorf("incorrect transition ID: got %s, want %s", transition.Transition.ID, tt.newStatus)
				}

				return tt.mockPostError
			}

			// Call the function being tested
			err := UpdateIssueStatus(tt.issueKey, tt.newStatus, mockPostRequest)

			// Check error
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err != nil)
			}

			if err != nil && tt.expectedErrMsg != "" && err.Error() != tt.expectedErrMsg {
				t.Errorf("expected error message: %s, got: %s", tt.expectedErrMsg, err.Error())
			}
		})
	}
}

func TestGetAvailableTransitions(t *testing.T) {
	tests := []struct {
		name           string
		issue          Issue
		mockResponse   TransitionResponse
		mockError      error
		expectedResult []Transition
		expectedError  bool
	}{
		{
			name:  "successful fetch",
			issue: Issue{Key: "TEST-123"},
			mockResponse: TransitionResponse{
				Transitions: []Transition{
					{ID: "1", Name: "To Do"},
					{ID: "2", Name: "In Progress"},
				},
			},
			mockError:      nil,
			expectedResult: []Transition{{ID: "1", Name: "To Do"}, {ID: "2", Name: "In Progress"}},
			expectedError:  false,
		},
		{
			name:           "API error",
			issue:          Issue{Key: "TEST-456"},
			mockResponse:   TransitionResponse{},
			mockError:      fmt.Errorf("API error"),
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request function
			mockRequest := func(url string, target interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				
				// Check URL format
				expectedURL := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", jiraBaseURL, tt.issue.Key)
				if url != expectedURL {
					t.Errorf("incorrect URL: got %s, want %s", url, expectedURL)
				}
				
				// Marshal and unmarshal to simulate JSON response
				data, _ := json.Marshal(tt.mockResponse)
				return json.Unmarshal(data, target)
			}

			// Call the function
			transitions, err := GetAvailableTransitions(tt.issue, mockRequest)

			// Check error
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err != nil)
			}

			// Check transitions
			if len(transitions) != len(tt.expectedResult) {
				t.Errorf("expected %d transitions, got %d", len(tt.expectedResult), len(transitions))
				return
			}

			for i, transition := range transitions {
				if transition.ID != tt.expectedResult[i].ID || transition.Name != tt.expectedResult[i].Name {
					t.Errorf("transition %d mismatch: expected %+v, got %+v", i, tt.expectedResult[i], transition)
				}
			}
		})
	}
}
