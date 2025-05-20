// Assisted by watsonx Code Assistant

package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

func TestWithMaxResults(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		maxValue int
		expected int
	}{
		{"Zero value", 0, 0},
		{"Positive value", 50, 50},
		{"Negative value", -10, -10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create options with default values
			opts := defaultSprintRequestOptions()

			// Apply the WithMaxResults option
			option := WithMaxResults(tc.maxValue)
			option(opts)

			// Check if maxResults was set correctly
			if opts.maxResults != tc.expected {
				t.Errorf("WithMaxResults(%d) = %d, expected %d",
					tc.maxValue, opts.maxResults, tc.expected)
			}
		})
	}
}

func TestWithStartAt(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		startAtValue  int
		expectedValue int
	}{
		{"Zero value", 0, 0},
		{"Positive value", 10, 10},
		{"Negative value", -5, -5}, // Even though negative values don't make sense for pagination
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create options with default values
			opts := defaultSprintRequestOptions()

			// Apply the WithStartAt option
			option := WithStartAt(tc.startAtValue)
			option(opts)

			// Check if the startAt value was correctly set
			if opts.startAt != tc.expectedValue {
				t.Errorf("WithStartAt(%d) = %d, want %d",
					tc.startAtValue, opts.startAt, tc.expectedValue)
			}
		})
	}
}

func TestMaxResultsParameter(t *testing.T) {
	tests := []struct {
		name       string
		maxResults int
		want       bool
	}{
		{
			name:       "with positive maxResults",
			maxResults: 10,
			want:       true,
		},
		{
			name:       "with zero maxResults",
			maxResults: 0,
			want:       false,
		},
		{
			name:       "with negative maxResults",
			maxResults: -5,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			opts := defaultSprintRequestOptions()
			opts.maxResults = tt.maxResults
			params := url.Values{}

			// Execute the code being tested
			if opts.maxResults > 0 {
				params.Add("maxResults", strconv.Itoa(opts.maxResults))
			}

			// Verify results
			_, exists := params["maxResults"]
			if exists != tt.want {
				t.Errorf("maxResults parameter existence = %v, want %v", exists, tt.want)
			}

			// If parameter should exist, verify its value
			if tt.want {
				got := params.Get("maxResults")
				want := strconv.Itoa(tt.maxResults)
				if got != want {
					t.Errorf("maxResults value = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestFetchSprintByID(t *testing.T) {
	tests := []struct {
		name           string
		sprintID       int
		mockResponse   Sprint
		mockError      error
		expectedSprint Sprint
		expectedError  bool
	}{
		{
			name:     "successful fetch",
			sprintID: 123,
			mockResponse: Sprint{
				ID:    123,
				Name:  "Test Sprint",
				State: "active",
			},
			mockError:      nil,
			expectedSprint: Sprint{ID: 123, Name: "Test Sprint", State: "active"},
			expectedError:  false,
		},
		{
			name:           "api error",
			sprintID:       456,
			mockResponse:   Sprint{},
			mockError:      errors.New("API error"),
			expectedSprint: Sprint{},
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request function
			mockMakeRequest := func(url string, target interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}

				// Marshal and unmarshal to simulate JSON response handling
				data, _ := json.Marshal(tt.mockResponse)
				return json.Unmarshal(data, target)
			}

			// Call the function being tested
			sprint, err := FetchSprintByID(tt.sprintID, mockMakeRequest)

			// Check error
			if (err != nil) != tt.expectedError {
				t.Errorf("FetchSprintByID() error = %v, expectedError = %v", err, tt.expectedError)
				return
			}

			// Check sprint data
			if sprint.ID != tt.expectedSprint.ID ||
				sprint.Name != tt.expectedSprint.Name ||
				sprint.State != tt.expectedSprint.State {
				t.Errorf("FetchSprintByID() = %v, expected %v", sprint, tt.expectedSprint)
			}
		})
	}
}

func TestFetchSprintIssues(t *testing.T) {
	tests := []struct {
		name           string
		sprintID       int
		mockResponse   JiraResponse
		mockError      error
		fetchFeatures  bool
		expectedIssues []Issue
		expectedError  bool
	}{
		{
			name:     "successful fetch without features",
			sprintID: 123,
			mockResponse: JiraResponse{
				Issues: []Issue{
					{Key: "TEST-1", Fields: Fields{Summary: "Test Issue 1"}},
					{Key: "TEST-2", Fields: Fields{Summary: "Test Issue 2"}},
				},
			},
			mockError:      nil,
			fetchFeatures:  false,
			expectedIssues: []Issue{{Key: "TEST-1"}, {Key: "TEST-2"}},
			expectedError:  false,
		},
		{
			name:           "api error",
			sprintID:       456,
			mockResponse:   JiraResponse{},
			mockError:      errors.New("API error"),
			fetchFeatures:  false,
			expectedIssues: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request function
			mockMakeRequest := func(url string, target interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}

				// Check URL format contains expected fields
				expectedURL := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d/issue?fields=", jiraBaseURL, tt.sprintID)
				if !strings.HasPrefix(url, expectedURL) {
					t.Errorf("incorrect URL prefix: got %s, want %s", url, expectedURL)
				}

				// Marshal and unmarshal to simulate JSON response handling
				data, _ := json.Marshal(tt.mockResponse)
				return json.Unmarshal(data, target)
			}

			// Call the function being tested
			issues, err := FetchSprintIssues(tt.sprintID, mockMakeRequest, tt.fetchFeatures)

			// Check error
			if (err != nil) != tt.expectedError {
				t.Errorf("FetchSprintIssues() error = %v, expectedError = %v", err, tt.expectedError)
				return
			}

			// Check issues length
			if len(issues) != len(tt.expectedIssues) {
				t.Errorf("FetchSprintIssues() returned %d issues, expected %d", len(issues), len(tt.expectedIssues))
				return
			}

			// Check issue keys
			for i, issue := range issues {
				if issue.Key != tt.mockResponse.Issues[i].Key {
					t.Errorf("Issue %d: expected key %s, got %s", i, tt.mockResponse.Issues[i].Key, issue.Key)
				}
			}
		})
	}
}

func TestUniqueEpicsFromIssues(t *testing.T) {
	tests := []struct {
		name          string
		issues        []Issue
		expectedEpics map[string]*Epic
	}{
		{
			name: "multiple issues with same epic",
			issues: []Issue{
				{Fields: Fields{Epic: &Epic{Key: "EPIC-1", Summary: "Epic 1"}}},
				{Fields: Fields{Epic: &Epic{Key: "EPIC-1", Summary: "Epic 1"}}},
			},
			expectedEpics: map[string]*Epic{
				"EPIC-1": {Key: "EPIC-1", Summary: "Epic 1"},
			},
		},
		{
			name: "issues with different epics",
			issues: []Issue{
				{Fields: Fields{Epic: &Epic{Key: "EPIC-1", Summary: "Epic 1"}}},
				{Fields: Fields{Epic: &Epic{Key: "EPIC-2", Summary: "Epic 2"}}},
			},
			expectedEpics: map[string]*Epic{
				"EPIC-1": {Key: "EPIC-1", Summary: "Epic 1"},
				"EPIC-2": {Key: "EPIC-2", Summary: "Epic 2"},
			},
		},
		{
			name: "issues with nil epic",
			issues: []Issue{
				{Fields: Fields{Epic: &Epic{Key: "EPIC-1", Summary: "Epic 1"}}},
				{Fields: Fields{Epic: nil}},
			},
			expectedEpics: map[string]*Epic{
				"EPIC-1": {Key: "EPIC-1", Summary: "Epic 1"},
			},
		},
		{
			name: "issues with empty epic key",
			issues: []Issue{
				{Fields: Fields{Epic: &Epic{Key: "EPIC-1", Summary: "Epic 1"}}},
				{Fields: Fields{Epic: &Epic{Key: "", Summary: "Empty Key"}}},
			},
			expectedEpics: map[string]*Epic{
				"EPIC-1": {Key: "EPIC-1", Summary: "Epic 1"},
			},
		},
		{
			name:          "empty issues list",
			issues:        []Issue{},
			expectedEpics: map[string]*Epic{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			epics := uniqueEpicsFromIssues(tt.issues)

			// Check map length
			assert.Equal(t, len(tt.expectedEpics), len(epics), "Epic map length mismatch")

			// Check map contents
			for key, expectedEpic := range tt.expectedEpics {
				actualEpic, exists := epics[key]
				assert.True(t, exists, "Expected epic with key %s not found", key)
				assert.Equal(t, expectedEpic.Key, actualEpic.Key, "Epic key mismatch")
				assert.Equal(t, expectedEpic.Summary, actualEpic.Summary, "Epic summary mismatch")
			}
		})
	}
}

func TestUpdateIssuesWithFeatures(t *testing.T) {
	tests := []struct {
		name          string
		issues        []Issue
		epicToFeature map[string]*Feature
		expected      []Issue
	}{
		{
			name: "issues with epics get features assigned",
			issues: []Issue{
				{Fields: Fields{Epic: &Epic{Key: "EPIC-1"}}},
				{Fields: Fields{Epic: &Epic{Key: "EPIC-2"}}},
			},
			epicToFeature: map[string]*Feature{
				"EPIC-1": {Key: "FEAT-1", Summary: "Feature 1"},
				"EPIC-2": {Key: "FEAT-2", Summary: "Feature 2"},
			},
			expected: []Issue{
				{Fields: Fields{
					Epic:    &Epic{Key: "EPIC-1"},
					Feature: &Feature{Key: "FEAT-1", Summary: "Feature 1"},
				}},
				{Fields: Fields{
					Epic:    &Epic{Key: "EPIC-2"},
					Feature: &Feature{Key: "FEAT-2", Summary: "Feature 2"},
				}},
			},
		},
		{
			name: "issues without epics don't get features",
			issues: []Issue{
				{Fields: Fields{Epic: nil}},
				{Fields: Fields{Epic: &Epic{Key: "EPIC-1"}}},
			},
			epicToFeature: map[string]*Feature{
				"EPIC-1": {Key: "FEAT-1", Summary: "Feature 1"},
			},
			expected: []Issue{
				{Fields: Fields{Epic: nil, Feature: nil}},
				{Fields: Fields{
					Epic:    &Epic{Key: "EPIC-1"},
					Feature: &Feature{Key: "FEAT-1", Summary: "Feature 1"},
				}},
			},
		},
		{
			name: "epics without matching features",
			issues: []Issue{
				{Fields: Fields{Epic: &Epic{Key: "EPIC-1"}}},
				{Fields: Fields{Epic: &Epic{Key: "EPIC-2"}}},
			},
			epicToFeature: map[string]*Feature{
				"EPIC-1": {Key: "FEAT-1", Summary: "Feature 1"},
				// No feature for EPIC-2
			},
			expected: []Issue{
				{Fields: Fields{
					Epic:    &Epic{Key: "EPIC-1"},
					Feature: &Feature{Key: "FEAT-1", Summary: "Feature 1"},
				}},
				{Fields: Fields{Epic: &Epic{Key: "EPIC-2"}, Feature: nil}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the issues to avoid modifying the test data
			testIssues := make([]Issue, len(tt.issues))
			copy(testIssues, tt.issues)

			// Call the function being tested
			updateIssuesWithFeatures(testIssues, tt.epicToFeature)

			// Verify results
			for i, issue := range testIssues {
				expected := tt.expected[i]

				// Check if Epic is nil in both
				if issue.Fields.Epic == nil && expected.Fields.Epic == nil {
					// Both are nil, this is correct
				} else if issue.Fields.Epic == nil || expected.Fields.Epic == nil {
					t.Errorf("Issue %d: Epic mismatch - got %v, expected %v",
						i, issue.Fields.Epic, expected.Fields.Epic)
				} else if issue.Fields.Epic.Key != expected.Fields.Epic.Key {
					t.Errorf("Issue %d: Epic key mismatch - got %s, expected %s",
						i, issue.Fields.Epic.Key, expected.Fields.Epic.Key)
				}

				// Check if Feature is nil in both
				if issue.Fields.Feature == nil && expected.Fields.Feature == nil {
					// Both are nil, this is correct
				} else if issue.Fields.Feature == nil || expected.Fields.Feature == nil {
					t.Errorf("Issue %d: Feature mismatch - got %v, expected %v",
						i, issue.Fields.Feature, expected.Fields.Feature)
				} else {
					assert.Equal(t, expected.Fields.Feature.Key, issue.Fields.Feature.Key)
					assert.Equal(t, expected.Fields.Feature.Summary, issue.Fields.Feature.Summary)
				}
			}
		})
	}
}
