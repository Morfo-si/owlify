package jira

import (
	"fmt"
)

// FetchBoards retrieves all Jira boards associated with the specified project.
// It makes a GET request to the Jira API and returns the list of boards.
//
// Parameters:
//   - project: The project key or ID to fetch boards for
//   - makeGetRequest: Function to make the Jira API request
//
// Returns:
//   - []Board: Slice of Board objects if successful
//   - error: Error if the request fails
func FetchBoards(project string, makeGetRequest JiraRequestFunc) ([]Board, error) {
	var boardResp BoardResponse
	// JQL to find boards for the project and component
	boardSearchURL := fmt.Sprintf("%s/rest/agile/1.0/board?projectKeyOrId=%s", jiraBaseURL, project)

	if err := makeGetRequest(boardSearchURL, &boardResp); err != nil {
		return nil, err
	}
	return boardResp.Values, nil
}

// FetchBoardByName retrieves a Jira board with the specified name.
// It makes a GET request to the Jira API and returns the board.
//
// Parameters:
//   - name: The name of the board to fetch
//   - makeGetRequest: Function to make the Jira API request
//
// Returns:
//   - Board: Board object if successful
//   - error: Error if the request fails
func FetchBoardByName(name string, makeGetRequest JiraRequestFunc) (Board, error) {
	var boardResp BoardResponse
	// JQL to find boards for the project and component
	boardSearchURL := fmt.Sprintf("%s/rest/agile/1.0/board?name=%s", jiraBaseURL, name)

	if err := makeGetRequest(boardSearchURL, &boardResp); err != nil {
		return Board{}, err
	}
	if len(boardResp.Values) == 0 {
		return Board{}, fmt.Errorf("no board found with name %s", name)
	}
	return boardResp.Values[0], nil
}

// FetchBoardByID retrieves a Jira board with the specified ID.
// It makes a GET request to the Jira API and returns the board.
//
// Parameters:
//   - id: The ID of the board to fetch
//   - makeGetRequest: Function to make the Jira API request
//
// Returns:
//   - Board: Board object if successful
//   - error: Error if the request fails
func FetchBoardByID(id int, makeGetRequest JiraRequestFunc) (Board, error) {
	var boardResp Board
	// JQL to find boards for the project and component
	boardSearchURL := fmt.Sprintf("%s/rest/agile/1.0/board/%d", jiraBaseURL, id)

	if err := makeGetRequest(boardSearchURL, &boardResp); err != nil {
		return Board{}, err
	}
	return boardResp, nil
}
