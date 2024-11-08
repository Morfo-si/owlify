package jira

import (
	"fmt"
	"strings"
)

func GetIssue(issueKey string) ([]Issue, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", jiraBaseURL, issueKey)

	var issueData Issue
	if err := makeGetRequest(url, &issueData); err != nil {
		return nil, fmt.Errorf("error fetching issue %s: %v", issueKey, err)
	}

	return []Issue{issueData}, nil
}

func UpdateIssueStatus(issueKey string, newStatus string) error {
	issue, err := GetIssue(issueKey)
	if err != nil {
		return fmt.Errorf("error getting issue: %v", err)
	}

	// First get available transitions
	transitions, err := GetAvailableTransitions(issue[0])
	if err != nil {
		return fmt.Errorf("error getting available transitions: %v", err)
	}

	// Find the transition ID for the desired status
	var transitionID string
	for _, t := range transitions {
		if strings.EqualFold(t.Name, newStatus) {
			transitionID = t.ID
			break
		}
	}

	if transitionID == "" {
		return fmt.Errorf("no transition found for status: %s", newStatus)
	}

	// Create the transition payload
	payload := UpdateTransition{
		Transition: struct {
			ID string `json:"id"`
		}{ID: transitionID},
	}

	url := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", jiraBaseURL, issue[0].ID)
	if err := makePostRequest(url, payload); err != nil {
		return fmt.Errorf("error transitioning issue: %v", err)
	}

	return nil
}

func GetAvailableTransitions(issue Issue) ([]Transition, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", jiraBaseURL, issue.ID)

	var response TransitionResponse
	if err := makeGetRequest(url, &response); err != nil {
		return nil, fmt.Errorf("error fetching transitions for issue %s: %v", issue.ID, err)
	}

	return response.Transitions, nil
}
