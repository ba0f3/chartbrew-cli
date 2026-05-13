# AGENTS.md

This file provides guidance to AI coding agents when working with this repository.

> [!IMPORTANT]
> Whenever you change behavior, add a feature, or modify a command interface, update every affected document in the Documentation Checklist.

## Build, Test, Run

```bash
make build        # build to bin/chartbrew
make test         # go test -v ./...
make vet          # go vet ./...
make fmt          # go fmt ./...
make tidy         # go mod tidy
make lint         # golangci-lint run ./... (optional)

./bin/chartbrew --help
```

If the default Go build cache is not writable in the current sandbox, run commands with:

```bash
GOCACHE=$PWD/.cache/go-build go test ./...
```

`main.go` is a thin wrapper: `cmd.Execute()` then `os.Exit(cmd.HandleError(err))`.

## Architecture

`cmd/root.go` builds an injectable Cobra root command. `PersistentPreRunE` resolves config in this priority order:

1. CLI flags
2. Environment variables: `CHARTBREW_API_URL`, `CHARTBREW_TOKEN`
3. `.env` in the current working directory
4. `~/.config/chartbrew/config.json`

Internal packages:

- `internal/config`: config file, `.env`, env var, flag resolution and validation.
- `internal/client`: authenticated HTTP requests to Chartbrew.
- `internal/output`: JSON, markdown, raw, and error envelope writers.
- `internal/body`: JSON body loading from `--data`, `--data-file`, or stdin.

Resource commands are registered from `cmd/resources.go` through the shared route factory in `cmd/resource.go`.

## Key Patterns

- Output JSON by default and keep stdout parseable.
- Write errors to stderr as `{"error":true,"code":N,"message":"..."}`.
- Avoid interactive prompts. Every input must work through flags, env vars, config files, or stdin.
- Do not log or print tokens.
- Prefer env vars or config files for credentials. `--token` exists for automation override.
- V1 intentionally has no delete commands.

## Project Layout

```text
main.go
cmd/                 Cobra commands and route registration
internal/body/       JSON request body readers
internal/client/     Chartbrew HTTP client
internal/config/     Config resolution and validation
internal/output/     stdout/stderr formatters
docs/                Planning and future user-facing docs
skill/SKILL.md       Agent skill definition
AGENTS.md            Guidance for coding agents
CLAUDE.md            Guidance for Claude Code
```

## Documentation Checklist

| Document | Update when... |
|---|---|
| `README.md` | Any user-facing feature, flag, command, or workflow changes |
| `CLAUDE.md` | Architecture changes, new patterns, pre-run behavior changes |
| `AGENTS.md` | Architecture changes, new patterns, any guidance for agents |
| `skill/SKILL.md` | Any command added/removed/changed, new security notes, new workflows |
