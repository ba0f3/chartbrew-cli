# Chartbrew CLI

`chartbrew` is an agent-first CLI for the Chartbrew API. It is built for scripts and AI agents that need structured, non-interactive access to teams, dashboards, connections, datasets, data requests, and charts.

## Build and Test

```bash
make build    # Builds bin/chartbrew
make test     # Runs go test -v ./...
make vet      # Runs go vet ./...
make fmt      # Formats Go source
make tidy     # Updates module metadata
```

## Releases and Homebrew

GitHub Actions runs `go test` and `go vet` on every push and pull request. Pushing a version tag matching `v*` (for example `v0.1.0`) runs [GoReleaser](https://goreleaser.com/), which creates a GitHub release with archives and updates the Homebrew tap repository.

One-time setup:

1. Create an empty GitHub repository named `homebrew-tap` under the same owner as this project (for example `https://github.com/ba0f3/homebrew-tap`), default branch `main`.
2. In this repository’s GitHub **Settings → Secrets and variables → Actions**, add `HOMEBREW_TAP_GITHUB_TOKEN`: a personal access token with **Contents: Read and write** on that tap repository (the default `GITHUB_TOKEN` cannot push to another repo).

Then tag and push:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Install from the tap. GoReleaser publishes prebuilt binaries via Homebrew’s **cask** mechanism (replacing the old `brews` path); the CLI is still command-line only:

```bash
brew tap ba0f3/tap
brew install --cask chartbrew
```

To dry-run a release locally:

```bash
goreleaser release --snapshot --clean
```

## Configuration

Configuration resolves in this priority order:

1. CLI flags
2. Environment variables
3. `.env` in the current working directory
4. `~/.config/chartbrew/config.json`

Supported environment variables:

```bash
export CHARTBREW_API_URL="https://api.chartbrew.com"
export CHARTBREW_TOKEN="your-api-token"
```

Config file example:

```json
{
  "base_url": "https://api.chartbrew.com",
  "token": "your-api-token",
  "output": "json"
}
```

Global flags:

```bash
chartbrew --base-url https://api.chartbrew.com --token "$CHARTBREW_TOKEN" --output json teams list
chartbrew --config ./chartbrew.json teams list
```

Prefer `CHARTBREW_TOKEN`, `.env`, or the config file for tokens. The `--token` flag exists for automation override, but command-line secrets can be stored in shell history.

## Commands

All commands write structured JSON by default. Use `--output markdown` for a fenced JSON block or `--output raw` for compact JSON.

Teams:

```bash
chartbrew teams list
chartbrew teams get --team-id 123
chartbrew teams create --data '{"name":"Demo"}'
chartbrew teams update --team-id 123 --data-file team.json
```

Dashboards use Chartbrew `project` API endpoints:

```bash
chartbrew dashboards list --team-id 123
chartbrew dashboards get --dashboard-id 456
chartbrew dashboards create --data-file dashboard.json
chartbrew dashboards update --dashboard-id 456 --data-file dashboard.json
```

Connections:

```bash
chartbrew connections list --team-id 123
chartbrew connections get --team-id 123 --connection-id 456
chartbrew connections create --team-id 123 --data-file connection.json
chartbrew connections update --team-id 123 --connection-id 456 --data-file connection.json
```

Datasets:

```bash
chartbrew datasets list --team-id 123
chartbrew datasets get --team-id 123 --dataset-id 456
chartbrew datasets create --team-id 123 --data-file dataset.json
chartbrew datasets update --team-id 123 --dataset-id 456 --data-file dataset.json
```

Data requests:

```bash
chartbrew data-requests list --team-id 123 --dataset-id 456
chartbrew data-requests get --team-id 123 --dataset-id 456 --request-id 789
chartbrew data-requests create --team-id 123 --dataset-id 456 --data-file request.json
chartbrew data-requests update --team-id 123 --dataset-id 456 --request-id 789 --data-file request.json
```

Charts:

```bash
chartbrew charts list --dashboard-id 456
chartbrew charts get --dashboard-id 456 --chart-id 789
chartbrew charts create --dashboard-id 456 --data-file chart.json
chartbrew charts update --dashboard-id 456 --chart-id 789 --data-file chart.json
```

## JSON Bodies

Create and update commands accept exactly one body source:

```bash
chartbrew teams create --data '{"name":"Demo"}'
chartbrew datasets create --team-id 123 --data-file dataset.json
printf '%s' '{"name":"Demo"}' | chartbrew teams create --data-file -
```

The CLI validates JSON before sending the request.

## V1 Scope

V1 includes list, get, create, and update commands only. Delete commands are intentionally excluded to keep agent-driven automation safer by default.
