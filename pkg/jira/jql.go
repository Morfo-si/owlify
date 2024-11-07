package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func FetchIssuesFromJQL(jql string) ([]Issue, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=%s", jiraBaseURL, url.QueryEscape(jql))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jiraToken))
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
