package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

var jiraBaseURL = fmt.Sprintf("%s/rest/api/2/search", os.Getenv("JIRA_BASE_URL"))

func FetchCurrentSprintIssues(project, component string, sprintNumber int) ([]JiraIssue, error) {
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

	url := fmt.Sprintf("%s?jql=%s", jiraBaseURL, url.QueryEscape(jql))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("JIRA_API_TOKEN")))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jiraResponse JiraResponse
	if err := json.NewDecoder(resp.Body).Decode(&jiraResponse); err != nil {
		return nil, err
	}

	return jiraResponse.Issues, nil
}
