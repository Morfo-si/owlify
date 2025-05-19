package jira

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	JIRA_URL_BOARD = "rest/agile/1.0/board"
	JIRA_URL_JQL   = "rest/api/2/search?jql"
)

// FetchSprintByID retrieves a Jira sprint with the specified ID.
// It makes a GET request to the Jira API and returns the sprint.
//
// Parameters:
//   - id: The ID of the sprint to fetch
//   - makeGetRequest: Function to make the Jira API request
//
// Returns:
//   - Sprint: Sprint object if successful
//   - error: Error if the request fails
func FetchSprintByID(id int, makeGetRequest JiraRequestFunc) (Sprint, error) {
	var sprintResp Sprint
	// JQL to find boards for the project and component
	sprintSearchURL := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d", jiraBaseURL, id)

	if err := makeGetRequest(sprintSearchURL, &sprintResp); err != nil {
		return Sprint{}, err
	}
	return sprintResp, nil
}

// FetchSprintIssues retrieves issues from a given sprint.
// It uses the provided makeGetRequest function to make the API call to JIRA.
//
// Parameters:
//   - sprintID: The ID of the sprint to fetch issues from
//   - makeGetRequest: Function to make the HTTP GET request to JIRA
//
// Returns:
//   - []Issue: A slice of Issue objects representing the issues in the sprint
//   - error: An error if the request fails or the response cannot be parsed
func FetchSprintIssues(sprintID int, makeGetRequest JiraRequestFunc) ([]Issue, error) {
	var jiraResponse JiraResponse

	sprintSearchURL := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d/issue", jiraBaseURL, sprintID)

	if err := makeGetRequest(sprintSearchURL, &jiraResponse); err != nil {
		return nil, err
	}

	return jiraResponse.Issues, nil
}

// FetchSprintIssuesWithEpic retrieves issues from a given sprint with epic information included.
// It uses the provided makeGetRequest function to make the API call to JIRA.
//
// Parameters:
//   - sprintID: The ID of the sprint to fetch issues from
//   - makeGetRequest: Function to make the HTTP GET request to JIRA
//
// Returns:
//   - []Issue: A slice of Issue objects representing the issues in the sprint with epic information
//   - error: An error if the request fails or the response cannot be parsed
func FetchSprintIssuesWithEpic(sprintID int, makeGetRequest JiraRequestFunc) ([]Issue, error) {
	var jiraResponse JiraResponse
	fields := []string{
		"summary",
		"assignee",
		"status",
		"priority",
		"customfield_12310243", // Story Points
		"duedate",
		"epic",
		"issuetype",
	}
	fieldsStr := strings.Join(fields, ",")

	// Use the fields parameter to expand epic information
	// Using rest/agile/1.0 API which includes epic data in the response
	sprintSearchURL := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d/issue?fields=%s", jiraBaseURL, sprintID, fieldsStr)

	if err := makeGetRequest(sprintSearchURL, &jiraResponse); err != nil {
		return nil, err
	}

	for i, issue := range jiraResponse.Issues {
		if issue.Fields.Epic != nil && issue.Fields.Epic.Key != "" {
			// Fetch epic details
			// Add 1 second delay to avoid rate limiting
			time.Sleep(1 * time.Second)
			epicIssue, err := GetEpic(issue.Fields.Epic.Key)
			if err != nil {
				log.Printf("Failed to fetch epic details for issue %s: %v", issue.Key, err)
				return nil, err
			}
			// Check if Feature key exists (can't use nil check on struct)
			if epicIssue.Feature.Key != "" {
				// Safe to access Feature structure (which is an embedded struct, not a pointer)
				featureKey := epicIssue.Feature.Key
				if featureKey != "" {
					// We need to ensure Epic.Feature is properly initialized in the destination
					// Initialize a Feature struct explicitly if needed
					// Fields.Epic is guaranteed to be non-nil due to the outer check

					// Update Epic with the feature
					// We need to assign to Fields.Epic, which might require recreating the Epic struct
					// Get the current Epic

					newFeature := &Feature{
						Summary: epicIssue.Feature.Summary,
						Key:     featureKey,
					}

					// Assign the updated Epic back
					jiraResponse.Issues[i].Fields.Feature = newFeature
				}
			} else {
				log.Printf("No Feature found for issue %s", issue.Key)
			}
		}
	}
	return jiraResponse.Issues, nil
}

// FetchOpenSprints retrieves all active sprints for a given board ID.
// It uses the provided makeGetRequest function to make the API call to JIRA.
//
// Parameters:
//   - boardId: The ID of the JIRA board to fetch sprints from
//   - makeGetRequest: Function to make the HTTP GET request to JIRA
//   - options: Optional parameters for the request (maxResults, startAt)
//
// Returns:
//   - []Sprint: A slice of Sprint objects representing the active sprints
//   - error: An error if the request fails or the response cannot be parsed
func FetchOpenSprints(boardId int, makeGetRequest JiraRequestFunc, options ...SprintRequestOption) ([]Sprint, error) {
	// Default options
	opts := defaultSprintRequestOptions()
	for _, option := range options {
		option(opts)
	}

	// Build URL with query parameters
	baseURL := fmt.Sprintf("%s/%s/%d/sprint", jiraBaseURL, JIRA_URL_BOARD, boardId)

	// Build query parameters
	params := url.Values{}
	params.Add("state", SprintStateActive)
	if opts.maxResults > 0 {
		params.Add("maxResults", strconv.Itoa(opts.maxResults))
	}
	if opts.startAt > 0 {
		params.Add("startAt", strconv.Itoa(opts.startAt))
	}

	sprintURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	var allSprints SprintResponse
	if err := makeGetRequest(sprintURL, &allSprints); err != nil {
		return nil, fmt.Errorf("failed to fetch active sprints for board %d: %w", boardId, err)
	}

	return allSprints.Values, nil
}

// SprintRequestOptions holds optional parameters for sprint requests
type sprintRequestOptions struct {
	maxResults int
	startAt    int
}

// SprintRequestOption is a function that modifies sprintRequestOptions
type SprintRequestOption func(*sprintRequestOptions)

// defaultSprintRequestOptions returns the default options
func defaultSprintRequestOptions() *sprintRequestOptions {
	return &sprintRequestOptions{
		maxResults: 0, // Use API default
		startAt:    0,
	}
}

// WithMaxResults sets the maximum number of results to return
func WithMaxResults(max int) SprintRequestOption {
	return func(o *sprintRequestOptions) {
		o.maxResults = max // Use API default if max is negative or zero
	}
}

// WithStartAt sets the index to start at for pagination
func WithStartAt(start int) SprintRequestOption {
	return func(o *sprintRequestOptions) {
		o.startAt = start
	}
}
