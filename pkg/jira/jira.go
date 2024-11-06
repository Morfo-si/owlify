package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	jiraBaseURL  string
	jiraUsername string
	jiraToken    string
)

func init() {
	// Try to load .env file from the current directory
	if err := godotenv.Load(); err != nil {
		// If .env doesn't exist, that's fine - we'll use environment variables
		fmt.Fprintf(os.Stderr, "Note: .env file not found, using environment variables\n")
	}

	// Get values from either .env file or environment variables
	jiraBaseURL = fmt.Sprintf("%s/rest/api/2/search", getEnvOrPanic("JIRA_BASE_URL"))
	jiraUsername = getEnvOrPanic("JIRA_USERNAME")
	jiraToken = getEnvOrPanic("JIRA_TOKEN")
}

// getEnvOrPanic retrieves an environment variable or panics if it's not set
func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s must be set in environment or .env file", key))
	}
	return value
}

func FetchCurrentSprintIssues(project, component string, sprintNumber int) ([]JiraIssue, error) {
	var jql string

	// If no sprint number is provided, fetch issues from all open sprints
	if sprintNumber == 0 {
		jql = "sprint in openSprints()"
	} else {
		jql = fmt.Sprintf("sprint = %d", sprintNumber)
	}

	if component != "" {
		jql += fmt.Sprintf(" AND component = '%s'", component)
	}

	if project != "" {
		jql += fmt.Sprintf(" AND project = '%s'", project)
	}

	url := fmt.Sprintf("%s?jql=%s", jiraBaseURL, url.QueryEscape(jql))

	req, err := http.NewRequest("GET", url, nil)
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

	var jiraResponse JiraResponse
	if err := json.NewDecoder(resp.Body).Decode(&jiraResponse); err != nil {
		return nil, err
	}

	return jiraResponse.Issues, nil
}
