package jira

import (
	"fmt"
	"net/url"
)

func FetchCurrentSprintIssues(project, component string, sprintNumber int) ([]Issue, error) {
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

	url := fmt.Sprintf("%s/rest/api/2/search?jql=%s", jiraBaseURL, url.QueryEscape(jql))

	var jiraResponse JiraResponse
	if err := makeGetRequest(url, &jiraResponse); err != nil {
		return nil, err
	}

	return jiraResponse.Issues, nil
}

func FetchOpenSprints(boardId int) ([]Sprint, error) {
	// JQL to find boards for the project and component
	// For each board, fetch active sprints
	var allSprints SprintResponse

	sprintURL := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?state=active", jiraBaseURL, boardId)

	if err := makeGetRequest(sprintURL, &allSprints); err != nil {
		return nil, err
	}

	return allSprints.Values, nil
}
