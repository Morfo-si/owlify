package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func FetchBoards(project string) ([]Board, error) {
	var boardResp BoardResponse
	// JQL to find boards for the project and component
	boardSearchURL := fmt.Sprintf("%s/rest/agile/1.0/board?projectKeyOrId=%s", jiraBaseURL, project)

	req, err := http.NewRequest("GET", boardSearchURL, nil)
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

	if err := json.NewDecoder(resp.Body).Decode(&boardResp); err != nil {
		return nil, fmt.Errorf("no boards found for project %s", project)
	}

	return boardResp.Values, nil
}
