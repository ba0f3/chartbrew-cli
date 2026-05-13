# Claude Code Guidance

## Identity & Philosophy

- You are an expert Go developer building an agent-first CLI for Chartbrew.
- Prioritize predictable command behavior, structured output, and scriptability.
- Produce JSON by default and never mix debug logs into stdout.
- Do not add interactive prompts.

## Build and Run

- Build: `make build`
- Run tests: `make test`
- Vet: `make vet`
- Format: `make fmt`
- Tidy modules: `make tidy`
- Releases: tag `v*` to trigger GoReleaser on GitHub Actions; see `.goreleaser.yaml` and `README.md` for the Homebrew tap.

Use `GOCACHE=$PWD/.cache/go-build` when the default Go build cache is read-only.

## Code Guidelines

- Keep `main.go` minimal.
- Use `github.com/spf13/cobra` for command structure.
- Keep shared behavior in focused internal packages:
  - `internal/config`
  - `internal/client`
  - `internal/output`
  - `internal/body`
- Add resource commands through the shared route factory instead of duplicating HTTP plumbing.
- All errors should be handled centrally and printed as JSON envelopes.
- Never leak secrets in logs, stdout, stderr, or test failures.

## Chartbrew CLI Behavior

- Config priority is flags, env vars, `.env`, then `~/.config/chartbrew/config.json`.
- Supported env vars are `CHARTBREW_API_URL` and `CHARTBREW_TOKEN`.
- Mutating commands accept JSON through `--data`, `--data-file <path>`, or `--data-file -`.
- Delete commands are available for documented Chartbrew delete endpoints only: dashboards, datasets, and charts.
- Delete execution requires `allow_delete: true` in the JSON config file. CLI flags, env vars, and `.env` must not enable deletes.

## Documentation

If you modify command behavior, flags, architecture, or workflows, update `README.md`, `AGENTS.md`, and `skill/SKILL.md` as applicable.
