package jira

import (
	"fmt"
)

func FetchBoards(project string) ([]Board, error) {
	var boardResp BoardResponse
	// JQL to find boards for the project and component
	boardSearchURL := fmt.Sprintf("%s/rest/agile/1.0/board?projectKeyOrId=%s", jiraBaseURL, project)

	if err := makeGetRequest(boardSearchURL, &boardResp); err != nil {
		return nil, err
	}
	return boardResp.Values, nil
}
