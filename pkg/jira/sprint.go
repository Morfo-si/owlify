package jira

import (
	"fmt"
	"net/url"
)

const (
	JIRA_URL_BOARD = "rest/agile/1.0/board"
	JIRA_URL_JQL   = "rest/api/2/search?jql"
)

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

func FetchOpenSprints(boardId int, makeGetRequest JiraRequestFunc) ([]Sprint, error) {
	// JQL to find boards for the project and component
	// For each board, fetch active sprints
	var allSprints SprintResponse

	sprintURL := fmt.Sprintf("%s/%s/%d/sprint?state=active", jiraBaseURL, JIRA_URL_BOARD, boardId)

	if err := makeGetRequest(sprintURL, &allSprints); err != nil {
		return nil, err
	}

	return allSprints.Values, nil
}
