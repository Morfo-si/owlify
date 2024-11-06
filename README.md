# owlify
Keep yourself informed about your tasks with Owlify.

## Configuration

Configuration can be provided either through environment variables or a `.env` file in the project root. If both exist, the `.env` file takes precedence.

### Required

```bash
export JIRA_BASE_URL="https://your-jira-instance.com" # Your JIRA instance URL
export JIRA_USERNAME="your.username" # Your JIRA username
export JIRA_TOKEN="your-api-token" # Your JIRA API token
```

### Optional

```bash
export JIRA_PROJECT="your-project-key" # Your JIRA project key
export JIRA_COMPONENT="your-component-key" # Your JIRA component key
```

### Using .env file:
Create a `.env` file in the project root:
```env
JIRA_BASE_URL="https://your-jira-instance.com"
JIRA_USERNAME="your.username"
JIRA_TOKEN="your-api-token"
```

Note: Make sure not to commit your `.env` file to version control as it may contain sensitive information.

## Usage

Basic usage with required project flag:

```bash
owlify -p PROJECT
```

All available options:
```bash
owlify -p PROJECT -c COMPONENT -s SPRINT_NUMBER
```

### Examples

1. Get all issues from a specific project:
```bash
owlify -p MYPROJECT
```

2. Get issues from a specific component:
```bash
owlify -p MYPROJECT -c "My Component"
```

3. Get issues from a specific sprint:
```bash
owlify -p MYPROJECT -s 123
```

### Flags

- `-p, --project`: JIRA project key (required)
- `-c, --component`: JIRA component (optional)
- `-s, --sprint`: Sprint number (optional)

## Building from source

```bash
go build -o owlify
```

## License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
```

This README now includes:
1. Required and optional environment variables
2. Basic usage instructions
3. Example commands
4. Available flags and their descriptions
5. Build instructions

Feel free to customize it further based on your specific needs!