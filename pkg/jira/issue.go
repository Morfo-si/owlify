package jira

import "fmt"

func GetIssue(issueKey string) ([]Issue, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", jiraBaseURL, issueKey)

	var issueData Issue
	if err := makeGetRequest(url, &issueData); err != nil {
		return nil, fmt.Errorf("error fetching issue %s: %v", issueKey, err)
	}

	return []Issue{issueData}, nil
}
