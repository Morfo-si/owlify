package jira

import (
	"fmt"
	"strings"
)

func GetIssue(issueKey string, makeGetRequest JiraRequestFunc) (Issue, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", jiraBaseURL, issueKey)

	var issueData Issue
	if err := makeGetRequest(url, &issueData); err != nil {
		return Issue{}, fmt.Errorf("error fetching issue %s: %v", issueKey, err)
	}

	return issueData, nil
}

func GetEpic(issueKey string, makeGetRequest JiraRequestFunc) (EpicResponse, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", jiraBaseURL, issueKey)

	var issueData EpicResponse
	if err := makeGetRequest(url, &issueData); err != nil {
		return EpicResponse{}, fmt.Errorf("error fetching issue %s: %v", issueKey, err)
	}

	return issueData, nil
}

// GetValidTransitionValid checks if the transition is valid for the issue
func GetValidTransitionID(status string, transitions []Transition) string {
	var transitionName string
	for _, t := range transitions {
		if strings.EqualFold(t.Name, status) {
			transitionName = t.ID
			break
		}
	}
	return transitionName
}

func UpdateIssueStatus(issueKey string, newStatus string, makePostRequest JiraPostRequestFunc) error {
	// Create the transition payload
	payload := UpdateTransition{
		Transition: struct {
			ID string `json:"id"`
		}{ID: newStatus},
	}

	url := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", jiraBaseURL, issueKey)
	if err := makePostRequest(url, payload, nil); err != nil {
		return fmt.Errorf("error transitioning issue: %v", err)
	}

	return nil
}

func GetAvailableTransitions(issue Issue, makeGetRequest JiraRequestFunc) ([]Transition, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", jiraBaseURL, issue.Key)

	var response TransitionResponse
	if err := makeGetRequest(url, &response); err != nil {
		return nil, fmt.Errorf("error fetching transitions for issue %s: %v", issue.Key, err)
	}

	return response.Transitions, nil
}
