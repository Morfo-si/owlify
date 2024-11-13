package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	jiraBaseURL  string
	jiraUsername string
	jiraToken    string
	httpProxy    string
	httpsProxy   string
)

func init() {
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

		envPath := fmt.Sprintf("%s/owlify/.env", configDir)
		if err := godotenv.Load(envPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
			return
		}
	}

	// Get values from either .env file or environment variables
	jiraBaseURL = getEnvOrPanic("JIRA_BASE_URL")
	jiraUsername = getEnvOrPanic("JIRA_USERNAME")
	jiraToken = getEnvOrPanic("JIRA_TOKEN")

	// Get proxy settings from environment
	httpProxy = os.Getenv("HTTP_PROXY")
	if httpProxy == "" {
		httpProxy = os.Getenv("http_proxy")
	}
	httpsProxy = os.Getenv("HTTPS_PROXY")
	if httpsProxy == "" {
		httpsProxy = os.Getenv("https_proxy")
	}
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

func makeGetRequest(reqUrl string, target interface{}) error {
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

func makePostRequest(reqUrl string, payload interface{}, target interface{}) error {
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
