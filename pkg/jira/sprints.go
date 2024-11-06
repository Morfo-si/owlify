package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func FetchOpenSprints(boardId int) ([]Sprint, error) {
	// JQL to find boards for the project and component
	// For each board, fetch active sprints
	var allSprints SprintResponse

	sprintURL := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?state=active", jiraBaseURL, boardId)

	req, err := http.NewRequest("GET", sprintURL, nil)
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

	if err := json.NewDecoder(resp.Body).Decode(&allSprints); err != nil {
		return nil, fmt.Errorf("no sprints found for board %d", boardId)
	}

	return allSprints.Values, nil
}
