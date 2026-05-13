# AGENTS.md

This file provides guidance to AI coding agents (Gemini CLI, Codex, GPT-Engineer, Claude Code, etc.)
when working with code in this repository.

> [!IMPORTANT]
> **Whenever you change behaviour, add a feature, or modify a command's interface,
> you MUST update every document listed in the [Documentation Checklist](#documentation-checklist)
> that is affected by the change.** Do not leave docs stale.

---

## Build, Test, Run

```bash
make build        # build to bin/cli (injects Version via ldflags from git describe)
make test         # go test -v ./...
make vet          # go vet ./...
make fmt          # go fmt ./...
make tidy         # go mod tidy
make lint         # golangci-lint run ./... (optional)

./bin/cli --help  # verify help text
```

`main.go` is a thin wrapper: `cmd.Execute()` then `os.Exit(cmd.HandleError(err))`.
All logic lives in `cmd/` and `internal/`.

---

## Architecture

### Execution flow
`cmd/root.go` defines a `PersistentPreRunE` that:
1. Parses config and flags into `config.Config`.
2. Resolves config in priority order: **CLI flags → Env Vars → `.env` in CWD → `~/.config/myapp/config.json`**.
3. Validates configuration.
4. Constructs global clients and formatters.

Subcommands read package-level globals or are injected with dependencies. Add a new command by creating `cmd/<resource>.go`, registering in `init()` via `rootCmd.AddCommand(...)`.

### Output formatters
Agents should prefer outputting structured JSON or Markdown so that other agents and tools can easily parse the CLI's output. Never mix stderr debug logs with stdout structured output.

---

## Key Patterns

- **Error contract**: Standardized JSON envelopes `{"error": true, "code": N, "message": "..."}` for errors.
- **Secure by default**: Credentials must be read securely. Avoid command line arguments for secrets (which leak to shell history), prefer reading from stdin or environment variables.
- **No interactive prompts in normal operation**: The tool must be fully scriptable by AI agents. All inputs must be achievable via flags, env vars, or stdin.

---

## Project Layout

```
main.go
cmd/                Cobra commands. root.go owns persistent flags + globals.
internal/config/    Config resolution, validation, and loading.
internal/client/    HTTP/API clients.
internal/output/    json/markdown/raw formatters.
docs/               User-facing docs.
skill/SKILL.md      AI-agent skill definition for Claude Code / Gemini / Cline.
AGENTS.md           This file — guidance for AI coding agents.
CLAUDE.md           Guidance for Claude Code specifically.
```

---

## Documentation Checklist

> [!IMPORTANT]
> **Always update every applicable document below when making a change.** Stale docs
> mislead users and other agents. A pull request that changes behaviour without
> updating docs will be rejected.

| Document | Update when… |
|---|---|
| [`README.md`](README.md) | Any user-facing feature, flag, or workflow changes |
| [`CLAUDE.md`](CLAUDE.md) | Architecture changes, new patterns, pre-run skip list changes |
| [`AGENTS.md`](AGENTS.md) | Architecture changes, new patterns, any guidance for agents |
| [`skill/SKILL.md`](skill/SKILL.md) | Any command added/removed/changed, new security notes, new workflows |

### Quick rule of thumb

- Changed a **flag or command interface**? → `README.md`, `skill/SKILL.md`
- Changed **internal architecture**? → `CLAUDE.md`, `AGENTS.md`
- Added a **new command**? → All of the above.
