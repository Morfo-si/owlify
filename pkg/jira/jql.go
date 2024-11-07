package jira

import (
	"fmt"
	"net/url"
)

func FetchIssuesFromJQL(jql string) ([]Issue, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=%s", jiraBaseURL, url.QueryEscape(jql))

	var jiraResponse JiraResponse
	if err := makeGetRequest(url, &jiraResponse); err != nil {
		return nil, err
	}

	return jiraResponse.Issues, nil
}
