package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/morfo-si/owlify/pkg/config"
)

var (
	jiraBaseURL string
	jiraToken   string
	httpProxy   string
	httpsProxy  string
)

func init() {
	// Initialize configuration
	if err := config.Initialize(); err != nil {
		panic(fmt.Sprintf("Failed to initialize configuration: %v", err))
	}
	
	// Get Jira config
	jiraConfig := config.GetJiraConfig()
	
	// Set package variables
	jiraBaseURL = jiraConfig.BaseURL
	jiraToken = jiraConfig.Token
	httpProxy = jiraConfig.HTTPProxy
	httpsProxy = jiraConfig.HTTPSProxy
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
