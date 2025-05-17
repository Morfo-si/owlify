package jira

import (
    "encoding/json"
    "fmt"
    "testing"
)

func TestFetchBoards(t *testing.T) {
    tests := []struct {
        name           string
        project        string
        mockResponse   BoardResponse
        mockError      error
        expectedBoards []Board
        expectedError  bool
    }{
        {
            name:    "successful fetch",
            project: "TEST",
            mockResponse: BoardResponse{
                Values: []Board{
                    {ID: 1, Name: "Board 1"},
                    {ID: 2, Name: "Board 2"},
                },
            },
            mockError:      nil,
            expectedBoards: []Board{{ID: 1, Name: "Board 1"}, {ID: 2, Name: "Board 2"}},
            expectedError:  false,
        },
        {
            name:           "API error",
            project:        "TEST",
            mockResponse:   BoardResponse{},
            mockError:      fmt.Errorf("API error"),
            expectedBoards: nil,
            expectedError:  true,
        },
        {
            name:    "empty response",
            project: "TEST",
            mockResponse: BoardResponse{
                Values: []Board{},
            },
            mockError:      nil,
            expectedBoards: []Board{},
            expectedError:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock request function
            mockRequest := func(url string, target interface{}) error {
                if tt.mockError != nil {
                    return tt.mockError
                }
                
                // Check URL format
                expectedURL := fmt.Sprintf("%s/rest/agile/1.0/board?projectKeyOrId=%s", jiraBaseURL, tt.project)
                if url != expectedURL {
                    t.Errorf("incorrect URL: got %s, want %s", url, expectedURL)
                }
                
                // Marshal and unmarshal to simulate JSON response
                data, _ := json.Marshal(tt.mockResponse)
                return json.Unmarshal(data, target)
            }

            // Call the function
            boards, err := FetchBoards(tt.project, mockRequest)

            // Check error
            if (err != nil) != tt.expectedError {
                t.Errorf("expected error: %v, got: %v", tt.expectedError, err != nil)
            }

            // Check boards
            if len(boards) != len(tt.expectedBoards) {
                t.Errorf("expected %d boards, got %d", len(tt.expectedBoards), len(boards))
                return
            }

            for i, board := range boards {
                if board.ID != tt.expectedBoards[i].ID || board.Name != tt.expectedBoards[i].Name {
                    t.Errorf("board %d mismatch: expected %+v, got %+v", i, tt.expectedBoards[i], board)
                }
            }
        })
    }
}

func TestFetchBoardByName(t *testing.T) {
    tests := []struct {
        name           string
        boardName      string
        mockResponse   BoardResponse
        mockError      error
        expectedBoard  Board
        expectedErrMsg string
    }{
        {
            name:      "successful board fetch",
            boardName: "Test Board",
            mockResponse: BoardResponse{
                Values: []Board{
                    {ID: 123, Name: "Test Board"},
                },
            },
            mockError:     nil,
            expectedBoard: Board{ID: 123, Name: "Test Board"},
        },
        {
            name:           "no boards found",
            boardName:      "Nonexistent Board",
            mockResponse:   BoardResponse{Values: []Board{}},
            mockError:      nil,
            expectedErrMsg: "no board found with name Nonexistent Board",
        },
        {
            name:           "API error",
            boardName:      "Error Board",
            mockResponse:   BoardResponse{},
            mockError:      fmt.Errorf("API connection error"),
            expectedErrMsg: "API connection error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockGet := func(url string, target interface{}) error {
                if tt.mockError != nil {
                    return tt.mockError
                }
                
                // Simulate the API response
                responseJSON, _ := json.Marshal(tt.mockResponse)
                json.Unmarshal(responseJSON, target)
                return nil
            }

            board, err := FetchBoardByName(tt.boardName, mockGet)

            // Check error cases
            if tt.expectedErrMsg != "" {
                if err == nil {
                    t.Errorf("expected error containing %q, got nil", tt.expectedErrMsg)
                    return
                }
                if err.Error() != tt.expectedErrMsg {
                    t.Errorf("expected error %q, got %q", tt.expectedErrMsg, err.Error())
                }
                return
            }

            // Check success case
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }

            if board.ID != tt.expectedBoard.ID || board.Name != tt.expectedBoard.Name {
                t.Errorf("expected board %+v, got %+v", tt.expectedBoard, board)
            }
        })
    }
}

func TestFetchBoardByID(t *testing.T) {
    tests := []struct {
        name           string
        id             int
        mockResponse   Board
        mockError      error
        expectedBoard  Board
        expectedErrMsg string
    }{
        {
            name: "successful board fetch",
            id:   123,
            mockResponse: Board{
                ID:   123,
                Name: "Test Board",
            },
            mockError:     nil,
            expectedBoard: Board{ID: 123, Name: "Test Board"},
        },
        {
            name:           "api error",
            id:             456,
            mockResponse:   Board{},
            mockError:      fmt.Errorf("API connection error"),
            expectedErrMsg: "API connection error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock request function
            mockRequest := func(url string, target interface{}) error {
                if tt.mockError != nil {
                    return tt.mockError
                }
                
                // Marshal and unmarshal to simulate JSON response handling
                data, _ := json.Marshal(tt.mockResponse)
                return json.Unmarshal(data, target)
            }

            // Call the function being tested
            board, err := FetchBoardByID(tt.id, mockRequest)

            // Check error
            if tt.mockError != nil {
                if err == nil {
                    t.Errorf("Expected error but got nil")
                } else if err.Error() != tt.expectedErrMsg {
                    t.Errorf("Expected error '%s', got '%s'", tt.expectedErrMsg, err.Error())
                }
                return
            }

            // Check success case
            if err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
            
            if board.ID != tt.expectedBoard.ID || board.Name != tt.expectedBoard.Name {
                t.Errorf("Expected board %+v, got %+v", tt.expectedBoard, board)
            }
        })
    }
}
