package config

import (
    "fmt"
    "os"
    "strings"

    "github.com/joho/godotenv"
)

// JiraConfig holds all Jira-related configuration
type JiraConfig struct {
    BaseURL    string
    Token      string
    HTTPProxy  string
    HTTPSProxy string
}

var (
    // Global configuration instance
    jiraConfig JiraConfig
)

// Initialize loads configuration from environment variables or .env file
func Initialize() error {
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

    // Try to load .env from current directory first
    var envPath string

    if err := godotenv.Load(); err != nil {
        // If .env doesn't exist, try to create one in XDG config dir
        configDir, err := os.UserConfigDir()
        if err != nil {
            return fmt.Errorf("error getting config directory: %v", err)
        }
        
        appConfigDir := fmt.Sprintf("%s/owlify", configDir)
        if err := os.MkdirAll(appConfigDir, 0755); err != nil {
            return fmt.Errorf("error creating config directory: %v", err)
        }
        
        envPath = fmt.Sprintf("%s/.env", appConfigDir)
        
        // Check if .env already exists in config dir
        if _, err := os.Stat(envPath); os.IsNotExist(err) {
            // Create new .env file from example
            exampleEnv, err := os.ReadFile(".env.example")
            if err != nil {
                return fmt.Errorf("error reading .env.example: %v", err)
            }

            if err := os.WriteFile(envPath, exampleEnv, 0600); err != nil {
                return fmt.Errorf("error creating .env file: %v", err)
            }

            fmt.Printf("Created new .env file at %s. Please edit it with your credentials.\n", envPath)
        }
        
        // Try to load the config file we just created/found
        if err := godotenv.Load(envPath); err != nil {
            return fmt.Errorf("error loading .env file: %v", err)
        }
    }

    // Load environment variables
    jiraConfig.BaseURL = getEnvOrDefault("JIRA_BASE_URL", "")
    jiraConfig.Token = getEnvOrDefault("JIRA_TOKEN", "")
    jiraConfig.HTTPProxy = os.Getenv("HTTP_PROXY")
    jiraConfig.HTTPSProxy = os.Getenv("HTTPS_PROXY")
    
    // Validate required config
    if jiraConfig.BaseURL == "" {
        return fmt.Errorf("JIRA_BASE_URL must be set in environment or .env file")
    }
    if jiraConfig.Token == "" {
        return fmt.Errorf("JIRA_TOKEN must be set in environment or .env file")
    }
    
    return nil
}

// GetJiraConfig returns the current Jira configuration
func GetJiraConfig() JiraConfig {
    return jiraConfig
}

// getEnvOrDefault retrieves an environment variable or returns the default if not set
func getEnvOrDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}