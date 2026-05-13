# Claude Code Guidance

This file contains custom instructions for Claude Code when working in this repository.

## Identity & Philosophy
- You are an expert Go developer and AI Agent building an "Agent-First" CLI application.
- Prioritize making the CLI tool easily parsable and orchestratable by other LLM agents.
- Produce structured output (JSON/Markdown) by default.
- Never add interactive prompts (e.g., `survey`, `promptui`) without a `--non-interactive` or equivalent flag, or just avoid them entirely to ensure scriptability.

## Build and Run
- Build: `make build`
- Run tests: `make test`
- Formatting: `make fmt` and `make tidy`

## Code Guidelines
- Keep `main.go` minimal.
- Use `github.com/spf13/cobra` for command structure.
- Always add new commands to `cmd/` package.
- All errors should be handled centrally and print a clear error code and message.
- Never leak secrets in logs, stdout, or crash reports.

## Documentation
- If you modify how a command works, its flags, or add a new command, you MUST update `AGENTS.md` and `skill/SKILL.md`.
- Keep the skill definition up-to-date so other agents know how to use this tool once it is compiled.
