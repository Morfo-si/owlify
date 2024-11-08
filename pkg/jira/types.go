package jira

// Issue represents a JIRA issue
type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Summary  string `json:"summary"`
		DueDate  string `json:"duedate"`
		Priority struct {
			Name string `json:"name"`
		} `json:"priority"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
}

type JiraResponse struct {
	Issues []Issue `json:"issues"`
}

type Sprint struct {
	ID            int    `json:"id"`
	Self          string `json:"self"`
	State         string `json:"state"`
	Name          string `json:"name"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	ActivatedDate string `json:"activatedDate"`
	OriginBoardId int    `json:"originBoardId"`
	// Goal          string `json:"goal"`
	Synced        bool `json:"synced"`
	AutoStartStop bool `json:"autoStartStop"`
}

type SprintResponse struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	IsLast     bool     `json:"isLast"`
	Values     []Sprint `json:"values"`
}

type Board struct {
	ID   int    `json:"id"`
	Self string `json:"self"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type BoardResponse struct {
	MaxResults int     `json:"maxResults"`
	StartAt    int     `json:"startAt"`
	Total      int     `json:"total"`
	IsLast     bool    `json:"isLast"`
	Values     []Board `json:"values"`
}

type TransitionResponse struct {
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateTransition struct {
	Transition struct {
		ID string `json:"id"`
	} `json:"transition"`
}
