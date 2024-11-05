package jira

// JiraIssue represents a JIRA issue
type JiraIssue struct {
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
	Issues []JiraIssue `json:"issues"`
}
