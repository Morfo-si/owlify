
package jira

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	JIRA_URL_BOARD = "rest/agile/1.0/board"
	JIRA_URL_JQL   = "rest/api/2/search?jql"
)

// FetchCurrentSprintIssues retrieves issues from the current sprint for a given project and component.
// It uses the provided makeGetRequest function to make the API call to JIRA.
// Parameters:
//   - project: The project key or ID to filter issues by
//   - component: The component name to filter issues by
//   - sprintNumber: The sprint number to filter issues by (0 for all open sprints)
//   - makeGetRequest: Function to make the HTTP GET request to JIRA
//
// Returns:
//   - []Issue: A slice of Issue objects representing the issues in the current sprint
//   - error: An error if the request fails or the response cannot be parsed
func FetchCurrentSprintIssues(project, component string, sprintNumber int, makeGetRequest JiraRequestFunc) ([]Issue, error) {
	var jql string

	// If no sprint number is provided, fetch issues from all open sprints
	if sprintNumber == 0 {
		jql = "sprint in openSprints()"
	} else {
		jql = fmt.Sprintf("sprint = %d", sprintNumber)
	}

	if component != "" {
		jql += fmt.Sprintf(" AND component = '%s'", component)
	}

	if project != "" {
		jql += fmt.Sprintf(" AND project = '%s'", project)
	}

	url := fmt.Sprintf("%s/%s=%s", jiraBaseURL, JIRA_URL_JQL, url.QueryEscape(jql))

	var jiraResponse JiraResponse
	if err := makeGetRequest(url, &jiraResponse); err != nil {
		return nil, err
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
