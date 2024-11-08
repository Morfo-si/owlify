package jira

import (
	"bytes"
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
	httpProxy    string
	httpsProxy   string
)

func init() {
	// Try to load .env file from the current directory
	if err := godotenv.Load(); err != nil {
		// If .env doesn't exist, that's fine - we'll use environment variables
		fmt.Fprintf(os.Stderr, "Note: .env file not found, using environment variables\n")
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

func makeGetRequest(reqUrl string, target interface{}) error {
	// Create transport with proxy support
	transport := &http.Transport{}
	if httpProxy != "" || httpsProxy != "" {
		transport.Proxy = http.ProxyFromEnvironment
	} else {
		proxy, err := url.Parse(httpProxy)
		if err != nil {
			return err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
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
	// Create transport with proxy support
	transport := &http.Transport{}
	if httpProxy != "" || httpsProxy != "" {
		transport.Proxy = http.ProxyFromEnvironment
	} else {
		proxy, err := url.Parse(httpProxy)
		if err != nil {
			return err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

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
