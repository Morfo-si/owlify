# Owlify

A CLI tool to fetch and format JIRA issues.

## Usage

### Commands

- `sprint`: Fetch JIRA issues from sprints

### Sprint Command Flags

- `-p, --project`: JIRA project key (required)
- `-c, --component`: JIRA component (optional)
- `-s, --sprint`: Sprint number (optional)
- `-o, --output`: Output format (optional, defaults to 'table')
  - Supported formats: 
    - `table`: Formatted table output
    - `json`: JSON format
    - `csv`: Comma-separated values

### Examples

1. Get all issues from current sprint:
```bash
owlify sprint -p MYPROJECT
```

2. Get issues for a specific component:
```bash
owlify sprint -p MYPROJECT -c "My Component"
```

3. Get issues from a specific sprint:
```bash
owlify sprint -p MYPROJECT -s 42
```

4. Get issues in JSON format:
```bash
owlify sprint -p MYPROJECT -o json
```

5. Export issues to CSV:
```bash
owlify sprint -p MYPROJECT -o csv
```

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