package jira

import (
	"encoding/json"
	"time"
)

type Fields struct {
	Summary    string   `json:"summary"`
	Assignee   Assignee `json:"assignee"`
	StoryPoint float64  `json:"customfield_12310243"`
	DueDate    string   `json:"duedate"`
	Priority   Priority `json:"priority"`
	Status     Status   `json:"status"`
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

// Issue represents a JIRA issue
type Issue struct {
	ID     string `json:"id"`
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
	if temp.StartDate != "" {
		if t, err := time.Parse("2006-01-02T15:04:05.999-0700", temp.StartDate); err == nil {
			s.StartDate = &t
		}
	}
	if temp.EndDate != "" {
		if t, err := time.Parse("2006-01-02T15:04:05.999-0700", temp.EndDate); err == nil {
			s.EndDate = &t
		}
	}
	if temp.ActivatedDate != "" {
		if t, err := time.Parse("2006-01-02T15:04:05.999-0700", temp.ActivatedDate); err == nil {
			s.ActivatedDate = &t
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
