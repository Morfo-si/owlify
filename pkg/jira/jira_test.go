package jira

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up mock environment variables for testing
	os.Setenv("JIRA_BASE_URL", "https://example.atlassian.net")
	os.Setenv("JIRA_TOKEN", "mock-token")

	// Run the tests
	exitCode := m.Run()

	// Clean up
	os.Unsetenv("JIRA_BASE_URL")
	os.Unsetenv("JIRA_TOKEN")

	os.Exit(exitCode)
}
