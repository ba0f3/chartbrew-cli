# Chartbrew CLI V1 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go-based, agent-first `chartbrew` CLI that provides list/get/create/update commands for Chartbrew teams, dashboards, connections, datasets, data requests, and charts.

**Architecture:** Keep Cobra commands in `cmd/`, shared configuration in `internal/config`, HTTP behavior in `internal/client`, and stdout formatting in `internal/output`. Mutating commands accept JSON bodies through `--data`, `--data-file`, or stdin, while command names and required path flags stay typed and discoverable.

**Tech Stack:** Go, Cobra, standard `net/http`, standard `encoding/json`, standard `testing`, `httptest`.

---

## File Structure

- Modify: `go.mod` to rename the module from the template path to this repository path.
- Modify: `main.go` to import the renamed module path.
- Modify: `Makefile` to build `bin/chartbrew` and inject the correct module path for `cmd.Version`.
- Modify: `cmd/root.go` to rename the binary, add persistent config flags, initialize config/client/output dependencies, and expose testable command construction.
- Modify: `cmd/version.go` to use the shared output writer.
- Create: `internal/config/config.go` for flag/env/`.env`/config file resolution.
- Create: `internal/config/config_test.go` for configuration priority and validation tests.
- Create: `internal/client/client.go` for authenticated Chartbrew HTTP requests.
- Create: `internal/client/client_test.go` for request construction, auth, body handling, and HTTP errors.
- Create: `internal/output/output.go` for JSON, markdown, raw, and error envelopes.
- Create: `internal/output/output_test.go` for formatter behavior.
- Create: `internal/body/body.go` for reading JSON payloads from `--data`, `--data-file`, or stdin.
- Create: `internal/body/body_test.go` for input precedence and invalid JSON tests.
- Create: `cmd/resource.go` for shared resource command builders.
- Create: `cmd/resource_test.go` for command path and method coverage.
- Modify: `README.md`, `AGENTS.md`, `CLAUDE.md`, and `skill/SKILL.md` to document the real Chartbrew CLI behavior.

## Command Surface

Global flags:

- `--base-url`: Chartbrew API base URL.
- `--token`: Chartbrew API token. Supported for automation, but docs recommend env or config file.
- `--config`: config file path, defaulting to `~/.config/chartbrew/config.json`.
- `--output`, `-o`: `json`, `markdown`, or `raw`; default `json`.

Config priority:

1. CLI flags
2. Environment variables: `CHARTBREW_API_URL`, `CHARTBREW_TOKEN`
3. `.env` in the current working directory
4. Config file: `~/.config/chartbrew/config.json`

Resource commands:

- `teams list`
- `teams get --team-id <id>`
- `teams create --data <json>|--data-file <path>|--data-file -`
- `teams update --team-id <id> --data <json>|--data-file <path>|--data-file -`
- `dashboards list --team-id <id>`
- `dashboards get --team-id <id> --dashboard-id <id>`
- `dashboards create --team-id <id> --data ...`
- `dashboards update --team-id <id> --dashboard-id <id> --data ...`
- `connections list --team-id <id>`
- `connections get --team-id <id> --connection-id <id>`
- `connections create --team-id <id> --data ...`
- `connections update --team-id <id> --connection-id <id> --data ...`
- `datasets list --team-id <id>`
- `datasets get --team-id <id> --dataset-id <id>`
- `datasets create --team-id <id> --data ...`
- `datasets update --team-id <id> --dataset-id <id> --data ...`
- `data-requests list --team-id <id> --dataset-id <id>`
- `data-requests get --team-id <id> --dataset-id <id> --request-id <id>`
- `data-requests create --team-id <id> --dataset-id <id> --data ...`
- `data-requests update --team-id <id> --dataset-id <id> --request-id <id> --data ...`
- `charts list --team-id <id> --dashboard-id <id>`
- `charts get --team-id <id> --dashboard-id <id> --chart-id <id>`
- `charts create --team-id <id> --dashboard-id <id> --data ...`
- `charts update --team-id <id> --dashboard-id <id> --chart-id <id> --data ...`

No delete commands ship in v1.

## Task 1: Rename Template to Chartbrew CLI

**Files:**
- Modify: `go.mod`
- Modify: `main.go`
- Modify: `Makefile`
- Modify: `cmd/root.go`
- Test: `cmd/root_test.go`

- [ ] **Step 1: Write a root command test**

Create `cmd/root_test.go`:

```go
package cmd

import "testing"

func TestRootCommandMetadata(t *testing.T) {
	root := NewRootCommand()

	if root.Use != "chartbrew" {
		t.Fatalf("Use = %q, want chartbrew", root.Use)
	}
	if root.Short == "" {
		t.Fatal("Short help must not be empty")
	}

	flagNames := []string{"base-url", "token", "config", "output"}
	for _, name := range flagNames {
		if root.PersistentFlags().Lookup(name) == nil {
			t.Fatalf("missing persistent flag %q", name)
		}
	}
}
```

- [ ] **Step 2: Run the failing test**

Run: `go test ./cmd -run TestRootCommandMetadata -v`

Expected: FAIL because `NewRootCommand` does not exist.

- [ ] **Step 3: Rename module and root command**

Update `go.mod`:

```go
module github.com/ba0f3/chartbrew-cli

go 1.26.2

require github.com/spf13/cobra v1.10.2

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)
```

Update `main.go` import:

```go
import (
	"os"

	"github.com/ba0f3/chartbrew-cli/cmd"
)
```

Update `Makefile`:

```make
BINARY_NAME=chartbrew
LDFLAGS := -ldflags "-X 'github.com/ba0f3/chartbrew-cli/cmd.Version=$(VERSION)'"
```

Update `cmd/root.go` so it exposes `NewRootCommand()` and sets `Use: "chartbrew"`.

- [ ] **Step 4: Run the test**

Run: `go test ./cmd -run TestRootCommandMetadata -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go.mod main.go Makefile cmd/root.go cmd/root_test.go
git commit -m "chore: rename template to chartbrew cli"
```

## Task 2: Implement Configuration Resolution

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`
- Modify: `cmd/root.go`

- [ ] **Step 1: Write priority tests**

Create tests covering these cases:

```go
func TestResolvePriorityFlagsEnvDotenvFile(t *testing.T) {
	// Config file values are loaded first.
	// .env overrides config file.
	// Environment variables override .env.
	// Explicit flags override environment variables.
}

func TestValidateRequiresBaseURLAndToken(t *testing.T) {
	// Missing base URL returns a validation error.
	// Missing token returns a validation error.
}
```

Use `t.Setenv`, `t.TempDir`, and isolated fixture files so the tests never read a real user config.

- [ ] **Step 2: Run the failing tests**

Run: `go test ./internal/config -v`

Expected: FAIL because the package does not exist.

- [ ] **Step 3: Implement `internal/config`**

Implement:

```go
type Config struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
	Output  string `json:"output"`
}

type Sources struct {
	FlagBaseURL string
	FlagToken   string
	FlagOutput  string
	Env         map[string]string
	DotenvPath  string
	ConfigPath  string
}

func Resolve(src Sources) (Config, error)
func (c Config) Validate() error
```

`Resolve` must read config file first, then `.env`, then env map, then flags. `.env` supports simple `KEY=value` lines for `CHARTBREW_API_URL` and `CHARTBREW_TOKEN`.

- [ ] **Step 4: Run config tests**

Run: `go test ./internal/config -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/config cmd/root.go
git commit -m "feat: add chartbrew config resolution"
```

## Task 3: Implement Output Formatting and Error Envelopes

**Files:**
- Create: `internal/output/output.go`
- Create: `internal/output/output_test.go`
- Modify: `cmd/version.go`

- [ ] **Step 1: Write formatter tests**

Create tests for:

```go
func TestWriteJSON(t *testing.T)
func TestWriteRaw(t *testing.T)
func TestWriteErrorEnvelope(t *testing.T)
```

JSON output must be valid JSON. Error output must match:

```json
{"error":true,"code":1,"message":"example"}
```

- [ ] **Step 2: Run failing tests**

Run: `go test ./internal/output -v`

Expected: FAIL because the package does not exist.

- [ ] **Step 3: Implement formatter**

Implement:

```go
type Format string

const (
	JSON     Format = "json"
	Markdown Format = "markdown"
	Raw      Format = "raw"
)

func Write(w io.Writer, format Format, value any) error
func WriteError(w io.Writer, code int, message string) error
```

For v1, markdown may render JSON in a fenced `json` block so it remains parseable.

- [ ] **Step 4: Wire `version` through formatter**

Update `cmd/version.go` so it writes `map[string]string{"version": Version}` through `internal/output`.

- [ ] **Step 5: Run tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/output cmd/version.go
git commit -m "feat: add structured output formatting"
```

## Task 4: Implement JSON Body Reader

**Files:**
- Create: `internal/body/body.go`
- Create: `internal/body/body_test.go`

- [ ] **Step 1: Write body reader tests**

Create tests for:

```go
func TestReadInlineJSON(t *testing.T)
func TestReadJSONFile(t *testing.T)
func TestReadJSONFromStdin(t *testing.T)
func TestRejectInvalidJSON(t *testing.T)
func TestRequireExactlyOneBodySource(t *testing.T)
```

- [ ] **Step 2: Run failing tests**

Run: `go test ./internal/body -v`

Expected: FAIL because the package does not exist.

- [ ] **Step 3: Implement body reader**

Implement:

```go
type Source struct {
	Data     string
	DataFile string
	Stdin    io.Reader
}

func ReadJSON(src Source) ([]byte, error)
```

Rules:

- `Data` reads inline JSON.
- `DataFile` reads a file path.
- `DataFile == "-"` reads from stdin.
- Zero body sources returns an error for create/update commands.
- Multiple body sources returns an error.
- Invalid JSON returns an error before any HTTP request is sent.

- [ ] **Step 4: Run tests**

Run: `go test ./internal/body -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/body
git commit -m "feat: add json body reader"
```

## Task 5: Implement Chartbrew HTTP Client

**Files:**
- Create: `internal/client/client.go`
- Create: `internal/client/client_test.go`

- [ ] **Step 1: Write client tests**

Create tests for:

```go
func TestRequestAddsBearerToken(t *testing.T)
func TestRequestEncodesJSONBody(t *testing.T)
func TestRequestReturnsHTTPErrorEnvelope(t *testing.T)
func TestRequestTrimsBaseURLSlash(t *testing.T)
```

Use `httptest.Server` and verify request method, path, headers, and body.

- [ ] **Step 2: Run failing tests**

Run: `go test ./internal/client -v`

Expected: FAIL because the package does not exist.

- [ ] **Step 3: Implement client**

Implement:

```go
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string, httpClient *http.Client) *Client
func (c *Client) Do(ctx context.Context, method, path string, body []byte) (json.RawMessage, error)
```

The client must:

- Send `Authorization: Bearer <token>`.
- Send `Content-Type: application/json` for requests with a body.
- Decode successful responses as raw JSON.
- Return structured errors for non-2xx responses without printing secrets.

- [ ] **Step 4: Run tests**

Run: `go test ./internal/client -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/client
git commit -m "feat: add chartbrew http client"
```

## Task 6: Build Shared Resource Command Factory

**Files:**
- Create: `cmd/resource.go`
- Create: `cmd/resource_test.go`
- Modify: `cmd/root.go`

- [ ] **Step 1: Write command factory tests**

Create tests proving:

```go
func TestListCommandCallsExpectedPath(t *testing.T)
func TestGetCommandRequiresID(t *testing.T)
func TestCreateCommandReadsData(t *testing.T)
func TestUpdateCommandReadsDataAndID(t *testing.T)
```

Use a fake API interface:

```go
type fakeRequester struct {
	method string
	path   string
	body   string
}

func (f *fakeRequester) Do(ctx context.Context, method, path string, body []byte) (json.RawMessage, error)
```

- [ ] **Step 2: Run failing tests**

Run: `go test ./cmd -run 'Test(List|Get|Create|Update)Command' -v`

Expected: FAIL because `cmd/resource.go` does not exist.

- [ ] **Step 3: Implement factory**

Implement reusable builders that accept:

```go
type Route struct {
	Use      string
	Short    string
	Method   string
	Path     func(values map[string]string) string
	IDFlags  []string
	NeedsBody bool
}
```

Generated commands must validate required flags before calling the client and must write the response through `internal/output`.

- [ ] **Step 4: Run command tests**

Run: `go test ./cmd -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/resource.go cmd/resource_test.go cmd/root.go
git commit -m "feat: add reusable resource command factory"
```

## Task 7: Add Resource Commands

**Files:**
- Create: `cmd/teams.go`
- Create: `cmd/dashboards.go`
- Create: `cmd/connections.go`
- Create: `cmd/datasets.go`
- Create: `cmd/data_requests.go`
- Create: `cmd/charts.go`
- Modify: `cmd/root.go`
- Test: `cmd/resources_test.go`

- [ ] **Step 1: Write route registration tests**

Create tests that execute `--help` for each resource and assert the four v1 verbs are present:

```go
func TestResourceCommandsRegistered(t *testing.T) {
	resources := []string{"teams", "dashboards", "connections", "datasets", "data-requests", "charts"}
	for _, resource := range resources {
		// Execute: chartbrew <resource> --help
		// Assert: list, get, create, update are shown.
	}
}
```

- [ ] **Step 2: Run failing tests**

Run: `go test ./cmd -run TestResourceCommandsRegistered -v`

Expected: FAIL because resource command files do not exist.

- [ ] **Step 3: Implement route files**

Use the route definitions from the command surface section. Each resource file should register only routes for that resource. Use flags named `--team-id`, `--dashboard-id`, `--connection-id`, `--dataset-id`, `--request-id`, and `--chart-id`.

- [ ] **Step 4: Run tests**

Run: `go test ./cmd -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/teams.go cmd/dashboards.go cmd/connections.go cmd/datasets.go cmd/data_requests.go cmd/charts.go cmd/root.go cmd/resources_test.go
git commit -m "feat: add chartbrew resource commands"
```

## Task 8: Wire Root Execution and Central Error Handling

**Files:**
- Modify: `cmd/root.go`
- Modify: `main.go`
- Test: `cmd/root_execution_test.go`

- [ ] **Step 1: Write execution tests**

Create tests for:

```go
func TestMissingConfigPrintsJSONError(t *testing.T)
func TestInvalidOutputFormatFailsBeforeRequest(t *testing.T)
func TestCommandErrorDoesNotWritePartialSuccessToStdout(t *testing.T)
```

- [ ] **Step 2: Run failing tests**

Run: `go test ./cmd -run 'Test(Missing|Invalid|Command)' -v`

Expected: FAIL until root execution uses centralized error handling.

- [ ] **Step 3: Implement error handling**

Expose:

```go
func Execute() error
func HandleError(err error) int
```

`HandleError` writes JSON error envelopes to stderr and returns a non-zero exit code. `main.go` should remain a thin wrapper around `cmd.Execute()` and `cmd.HandleError(err)`.

- [ ] **Step 4: Run tests**

Run: `go test ./cmd -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add cmd/root.go main.go cmd/root_execution_test.go
git commit -m "feat: centralize cli error handling"
```

## Task 9: Update Documentation and Agent Skill

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `CLAUDE.md`
- Modify: `skill/SKILL.md`

- [ ] **Step 1: Update README**

Document:

- Install/build commands.
- Config priority order.
- Auth examples using env vars and config file.
- Resource command examples.
- JSON body input examples.
- No delete commands in v1.

- [ ] **Step 2: Update agent guidance**

Update `AGENTS.md` and `CLAUDE.md` so their architecture descriptions match `internal/config`, `internal/client`, `internal/output`, and `internal/body`.

- [ ] **Step 3: Update skill**

Update `skill/SKILL.md` so agents know to:

- Prefer `--output=json`.
- Set `CHARTBREW_API_URL` and `CHARTBREW_TOKEN`.
- Use `--data-file -` for secrets or long JSON bodies.
- Avoid looking for delete commands in v1.

- [ ] **Step 4: Verify docs mention real binary name**

Run: `rg 'agent-cli-template|\\bcli\\b|github.com/username' README.md AGENTS.md CLAUDE.md skill/SKILL.md main.go Makefile go.mod cmd`

Expected: no stale template references, except explanatory prose that intentionally mentions generic CLI concepts.

- [ ] **Step 5: Commit**

```bash
git add README.md AGENTS.md CLAUDE.md skill/SKILL.md
git commit -m "docs: document chartbrew cli usage"
```

## Task 10: Final Verification

**Files:**
- Modify only files needed to fix verification failures.

- [ ] **Step 1: Format**

Run: `make fmt`

Expected: Go files are formatted.

- [ ] **Step 2: Tidy modules**

Run: `make tidy`

Expected: `go.mod` and `go.sum` are consistent.

- [ ] **Step 3: Run all tests**

Run: `make test`

Expected: PASS for all packages.

- [ ] **Step 4: Run vet**

Run: `make vet`

Expected: PASS.

- [ ] **Step 5: Build binary**

Run: `make build`

Expected: `bin/chartbrew` exists.

- [ ] **Step 6: Check help output**

Run: `./bin/chartbrew --help`

Expected: help mentions Chartbrew, config flags, and `--output`.

- [ ] **Step 7: Commit verification fixes**

```bash
git add .
git commit -m "test: verify chartbrew cli v1"
```

## Self-Review Notes

- Spec coverage: The plan covers binary rename, config priority, secure token handling, structured output, shared HTTP client, JSON body input, list/get/create/update commands, documentation, and final verification.
- Placeholder scan: Tasks contain concrete file paths, commands, validation rules, and expected outcomes.
- Scope check: The plan deliberately excludes delete commands and higher-level workflow shortcuts from v1.
- Ambiguity check: Config priority, command names, body input rules, and output behavior are explicit.
