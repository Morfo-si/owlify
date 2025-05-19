# Owlify Release Notes

## Version 0.0.3

- Added new `sprint get` command to fetch information from a specific sprint by its numerical ID
- Added board and sprint APIs with enhanced epic support
- Default display now includes epic information for sprint issues (119a23f)

## Version 0.0.2

- 5401c09 Clean up command line flags.
- 119a23f Default to display sprint issues with epic info.
- e49bbb0 Improvements to sprint related functions.
- cb8d3f4 Removed obsoleted sprint function.
- 7154e86 feat: add board and sprint APIs with epic support and tests

## Version 0.0.1

- c737752 Added Assignee and Story Points.
- aab039c Added JSON output format.
- acdce9e Added PROXY to .env.example
- b0c62db Added default table output.
- 48d687d Added new search command.
- 627b4c0 Adding unittests for reports.go.
- 8de9fea Create XDG dir/.env automatically.
- ccfa6ec Display each field as column.
- e120f62 Fixed csv formatting for fields
- 96ede8d Fixed jira.makePostRequest.
- 957151c Grouped existing functions to new subcommand sprint.
- 12f6ac7 Include PROXY example.
- 17afc02 Initial commit
- a597784 Initial commit.
- d4fc736 Initial supprt for GET issue
- e32b909 Massive refactoring.
- f4e3978 Minor fix for Issues display
- 56a0735 Moved basic report to pkg.
- a16ec0d New createHTTPClient func.
- dabe4b7 Only use Proxy if available.
- 0fcd13c Renamed Text to Table format.
- c67a570 Renamed pkg/jira/sprints.go
- 329d0a7 Support for .env file.
- cd14418 Tweaks to Makefile.
- daf9a36 Unified makeGetRequest func.
- 956613a Update issue's status.
- e82bc41 Use goreleaser.
- e51ab1c fix: add main package and update build target
- 08771f5 fix: handle return values from MarkFlagRequired in issue commands
- 04892c9 fix: update jql tests to use relative dates for consistent test results
- f1be68d refactor: add JIRA URL constants and improve sprint functionality
- c8ed816 refactor: extract JIRA types to improve code organization
- d531a05 refactor: improve Fields date handling with proper time.Time type and helper methods
- eaa7a36 refactor: improve Sprint date handling with proper time.Time types and helper methods
- 69bdb59 refactor: introduce JiraRequestFunc type for better dependency injection
- dc38fc3 refactor: update API calls to use JIRAGetRequest
- 91b58aa refactor: update command layer to use new JIRA request function
