---
name: chartbrew-cli
description: Agent-first CLI for managing Chartbrew teams, dashboards, connections, datasets, data requests, and charts.
---

# Chartbrew CLI Skill

Use `chartbrew` to interact with the Chartbrew API from scripts and AI agents.

## Setup

Prefer environment variables or config files for credentials:

```bash
export CHARTBREW_API_URL="https://api.chartbrew.com"
export CHARTBREW_TOKEN="your-api-token"
```

Config resolves in this order:

1. CLI flags
2. `CHARTBREW_API_URL` and `CHARTBREW_TOKEN`
3. `.env` in the current working directory
4. `~/.config/chartbrew/config.json`

Config file example:

```json
{
  "base_url": "https://api.chartbrew.com",
  "token": "your-api-token",
  "output": "json"
}
```

## Operating Rules

- Use `--output=json` when another tool or agent will parse the result.
- Do not place sensitive connection passwords in shell history. Put long or sensitive JSON in a file or pipe it through `--data-file -`.
- Do not expect delete commands in v1. The CLI supports list, get, create, and update.
- Do not use interactive prompt workarounds such as `yes`; this CLI is designed to be non-interactive.

## Commands

```bash
chartbrew teams list
chartbrew teams get --team-id 123
chartbrew teams create --data '{"name":"Demo"}'
chartbrew teams update --team-id 123 --data-file team.json

chartbrew dashboards list --team-id 123
chartbrew dashboards get --dashboard-id 456
chartbrew dashboards create --data-file dashboard.json
chartbrew dashboards update --dashboard-id 456 --data-file dashboard.json

chartbrew connections list --team-id 123
chartbrew connections get --team-id 123 --connection-id 456
chartbrew connections create --team-id 123 --data-file connection.json
chartbrew connections update --team-id 123 --connection-id 456 --data-file connection.json

chartbrew datasets list --team-id 123
chartbrew datasets get --team-id 123 --dataset-id 456
chartbrew datasets create --team-id 123 --data-file dataset.json
chartbrew datasets update --team-id 123 --dataset-id 456 --data-file dataset.json

chartbrew data-requests list --team-id 123 --dataset-id 456
chartbrew data-requests get --team-id 123 --dataset-id 456 --request-id 789
chartbrew data-requests create --team-id 123 --dataset-id 456 --data-file request.json
chartbrew data-requests update --team-id 123 --dataset-id 456 --request-id 789 --data-file request.json

chartbrew charts list --dashboard-id 456
chartbrew charts get --dashboard-id 456 --chart-id 789
chartbrew charts create --dashboard-id 456 --data-file chart.json
chartbrew charts update --dashboard-id 456 --chart-id 789 --data-file chart.json
```

## JSON Body Input

Create and update commands require exactly one JSON body source:

```bash
chartbrew teams create --data '{"name":"Demo"}'
chartbrew datasets create --team-id 123 --data-file dataset.json
printf '%s' '{"name":"Demo"}' | chartbrew teams create --data-file -
```
