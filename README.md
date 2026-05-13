# Agent-First CLI Template

This is a GitHub project template for creating "Agent-First" CLI tools in Go. 
An "Agent-First" CLI is designed specifically to be executed, orchestrated, and understood by AI agents (like Claude Code, Gemini, GitHub Copilot Workspaces, etc.) as effectively as human users.

## Features

- **Cobra-based CLI**: Standard, robust command parsing and help generation.
- **Agent-Oriented Documentation**: Includes `AGENTS.md`, `CLAUDE.md`, and a predefined `skill/SKILL.md` to instantly teach AI agents how to use your CLI.
- **Structured Output**: Native support for returning data in JSON or Markdown to avoid messy stdout parsing.
- **Secure Secret Handling**: Standardized ways to consume secrets via stdin or env vars to prevent leaks into `.bash_history`.
- **Non-Interactive First**: Guaranteed to be fully automatable without hanging on `[y/N]` prompts.

## Quick Start

1. Click **Use this template** on GitHub.
2. Clone your new repository.
3. Replace `github.com/username/agent-cli-template` with your own module name in `go.mod` and all imports.
4. Run `make build` to build the binary.
5. Extend the tool by adding commands to the `cmd/` package.

## Building and Testing

```bash
make build    # Builds the binary into bin/cli
make test     # Runs tests
make lint     # Runs golangci-lint
make tidy     # Updates go.mod dependencies
```

## How to Make It "Agent-First"

- **Keep Help Text Descriptive**: Agents read `--help` output just like humans. Describe the exact expected formats for flags.
- **Maintain `skill/SKILL.md`**: Whenever you add features, document them in `skill/SKILL.md` so users can import your tool directly into their agent's skill library.
- **Avoid TTY Assumptions**: Never assume the tool is running in a fully interactive TTY. Rely on exit codes for success/failure, not just colored console output.
