package jira

import (
	"encoding/json"
	"time"
)

// Epic represents a JIRA epic
type Epic struct {
	Key     string `json:"key"`
	Summary string `json:"summary"`
}

// Fields represents the content fields of a JIRA issue
type Fields struct {
	Summary    string     `json:"summary"`
	Assignee   Assignee   `json:"assignee"`
	IssueType  IssueType  `json:"issuetype"`
	StoryPoint float64    `json:"customfield_12310243"`
	Priority   Priority   `json:"priority"`
	Status     Status     `json:"status"`
	Epic       *Epic      `json:"epic,omitempty"`
	DueDate    *time.Time `json:"duedate,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Fields
func (f *Fields) UnmarshalJSON(data []byte) error {
	type FieldsAlias Fields
	type FieldsTemp struct {
		*FieldsAlias
		DueDate string `json:"duedate"`
	}

	temp := &FieldsTemp{FieldsAlias: (*FieldsAlias)(f)}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	// Parse DueDate if it's not empty
	if temp.DueDate != "" {
		// Try different date formats
		formats := []string{
			"2006-01-02",
			"2006-01-02T15:04:05.000-0700",
			"2006-01-02T15:04:05.000Z",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, temp.DueDate); err == nil {
				f.DueDate = &t
				break
			}
		}
	}

	return nil
}

// IsOverdue returns true if the issue is past its due date
func (f Fields) IsOverdue() bool {
	if f.DueDate == nil {
		return false
	}
	return time.Now().After(*f.DueDate)
}

// DaysUntilDue returns the number of days until the issue is due
// Returns negative values for overdue issues
func (f Fields) DaysUntilDue() int {
	if f.DueDate == nil {
		return 0
	}

	now := time.Now()
	due := *f.DueDate

	// Normalize to midnight for consistent day calculations
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	due = time.Date(due.Year(), due.Month(), due.Day(), 0, 0, 0, 0, due.Location())

	days := int(due.Sub(now).Hours() / 24)
	return days
}

type Assignee struct {
	Name string `json:"name"`
}

type Priority struct {
	Name string `json:"name"`
}

type Status struct {
	Name string `json:"name"`
}

type IssueType struct {
	Name string `json:"name"`
}

// Issue represents a JIRA issue
type Issue struct {
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

type JiraResponse struct {
	Issues []Issue `json:"issues"`
}

// Sprint state constants
const (
	SprintStateActive = "active"
	SprintStateClosed = "closed"
	SprintStateFuture = "future"
)

// Sprint represents a JIRA sprint with proper time handling
type Sprint struct {
	ID            int        `json:"id"`
	Self          string     `json:"self"`
	State         string     `json:"state"`
	Name          string     `json:"name"`
	StartDate     *time.Time `json:"startDate,omitempty"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	ActivatedDate *time.Time `json:"activatedDate,omitempty"`
	OriginBoardId int        `json:"originBoardId"`
	// Goal          string `json:"goal"`
	Synced        bool `json:"synced"`
	AutoStartStop bool `json:"autoStartStop"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Sprint
func (s *Sprint) UnmarshalJSON(data []byte) error {
	type SprintAlias Sprint
	type SprintTemp struct {
		*SprintAlias
		StartDate     string `json:"startDate"`
		EndDate       string `json:"endDate"`
		ActivatedDate string `json:"activatedDate"`
	}

	temp := &SprintTemp{SprintAlias: (*SprintAlias)(s)}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	// Parse dates if they're not empty
	// Define formats to try
	formats := []string{
		"2006-01-02T15:04:05.999-0700",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	// Parse StartDate
	if temp.StartDate != "" {
		for _, format := range formats {
			if t, err := time.Parse(format, temp.StartDate); err == nil {
				s.StartDate = &t
				break
			}
		}
	}
	
	// Parse EndDate
	if temp.EndDate != "" {
		for _, format := range formats {
			if t, err := time.Parse(format, temp.EndDate); err == nil {
				s.EndDate = &t
				break
			}
		}
	}
	
	// Parse ActivatedDate
	if temp.ActivatedDate != "" {
		for _, format := range formats {
			if t, err := time.Parse(format, temp.ActivatedDate); err == nil {
				s.ActivatedDate = &t
				break
			}
		}
	}

	return nil
}

// IsActive returns true if the sprint is in active state
func (s Sprint) IsActive() bool {
	return s.State == SprintStateActive
}

// IsCompleted returns true if the sprint has ended
func (s Sprint) IsCompleted() bool {
	return s.EndDate != nil && time.Now().After(*s.EndDate)
}

// Duration returns the planned duration of the sprint
func (s Sprint) Duration() time.Duration {
	if s.StartDate == nil || s.EndDate == nil {
		return 0
	}
	return s.EndDate.Sub(*s.StartDate)
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
