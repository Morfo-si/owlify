package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	jiraBaseURL string
	jiraToken   string
	httpProxy   string
	httpsProxy  string
)

func init() {
	// Check if we're in a test environment
	if strings.HasSuffix(os.Args[0], ".test") || strings.HasSuffix(os.Args[0], "_test") {
		// In test environment, use mock values if not set
		if os.Getenv("JIRA_BASE_URL") == "" {
			os.Setenv("JIRA_BASE_URL", "https://example.atlassian.net")
		}
		if os.Getenv("JIRA_TOKEN") == "" {
			os.Setenv("JIRA_TOKEN", "mock-token")
		}
	}

	// Try to load .env file from the current directory
	// Try to load .env from current directory first
	var envPath string

	if err := godotenv.Load(); err != nil {
		// If .env doesn't exist, try to create one in XDG config dir
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
			appConfigDir := fmt.Sprintf("%s/owlify", configDir)
			if err := os.MkdirAll(appConfigDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
				return
			}
			exampleEnv, err := os.ReadFile(".env.example")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading .env.example: %v\n", err)
				return
			}

			if err := os.WriteFile(envPath, exampleEnv, 0600); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating .env file: %v\n", err)
				return
			}

			fmt.Printf("Created new .env file at %s. Please edit it with your credentials.\n", envPath)

		}
	}

	// Load environment variables
	jiraBaseURL = getEnvOrPanic("JIRA_BASE_URL")
	jiraToken = getEnvOrPanic("JIRA_TOKEN")
	httpProxy = os.Getenv("HTTP_PROXY")
	httpsProxy = os.Getenv("HTTPS_PROXY")
}

// getEnvOrPanic retrieves an environment variable or panics if it's not set
func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s must be set in environment or .env file", key))
	}
	return value
}

// createHTTPClient creates an http.Client with proxy support if configured
func createHTTPClient() *http.Client {
	transport := &http.Transport{}
	if httpProxy != "" || httpsProxy != "" {
		transport.Proxy = http.ProxyFromEnvironment
	}

	if transport.Proxy != nil {
		return &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		}
	}

	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

type JiraRequestFunc func(string, any) error
type JiraPostRequestFunc func(string, any, any) error

func JIRAGetRequest(reqUrl string, target any) error {
	// Replace client creation with new function
	client := createHTTPClient()

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jiraToken))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func JIRAPostRequest(reqUrl string, payload any, target any) error {
	// Replace client creation with new function
	client := createHTTPClient()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jiraToken))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}

	return nil
}
