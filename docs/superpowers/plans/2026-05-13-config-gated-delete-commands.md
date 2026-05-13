# Config-Gated Delete Commands Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add delete commands for Chartbrew resources, guarded by a config-file-only `allow_delete` opt-in.

**Architecture:** Extend `internal/config.Config` with `AllowDelete`, populated only from JSON config files. Carry that setting through `cmd/root.go` into route execution, and require it for routes marked as destructive before any HTTP request is made.

**Tech Stack:** Go, Cobra, standard `net/http`, standard `testing`.

---

## Task 1: Config-File Delete Opt-In

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`

- [ ] Add a failing test proving `allow_delete: true` is read from config JSON.
- [ ] Add a failing test proving env vars, `.env`, and CLI flags cannot enable `AllowDelete`.
- [ ] Add `AllowDelete bool` to `Config` and parse it only in `readConfigFile`.
- [ ] Run `GOCACHE=$PWD/.cache/go-build go test ./internal/config -v`.

## Task 2: Destructive Route Gate

**Files:**
- Modify: `cmd/root.go`
- Modify: `cmd/resource.go`
- Modify: `cmd/resource_test.go`

- [ ] Add a failing command test proving `datasets delete` refuses to call the API when `allow_delete` is absent.
- [ ] Add a failing command test proving `datasets delete` sends `DELETE /team/{team_id}/datasets/{dataset_id}` when config has `allow_delete: true`.
- [ ] Store `allowDelete` on `appState`.
- [ ] Add `Destructive bool` to `Route` and block destructive routes unless `state.allowDelete` is true.
- [ ] Run `GOCACHE=$PWD/.cache/go-build go test ./cmd -v`.

## Task 3: Add Documented Delete Routes

**Files:**
- Modify: `cmd/resources.go`
- Modify: `cmd/resources_test.go`

- [ ] Add delete routes for documented endpoints:
  - `dashboards delete --dashboard-id <id>` -> `DELETE /project/{id}`
  - `datasets delete --team-id <id> --dataset-id <id>` -> `DELETE /team/{team_id}/datasets/{dataset_id}`
  - `charts delete --dashboard-id <id> --chart-id <id>` -> `DELETE /project/{project_id}/chart/{id}`
- [ ] Update help tests so these resources include delete.
- [ ] Run `GOCACHE=$PWD/.cache/go-build go test ./cmd -v`.

## Task 4: Docs and Verification

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `CLAUDE.md`
- Modify: `skill/SKILL.md`

- [ ] Document `allow_delete`.
- [ ] Document that delete opt-in must come from the JSON config file.
- [ ] Run:
  - `GOCACHE=$PWD/.cache/go-build make fmt`
  - `GOCACHE=$PWD/.cache/go-build make test`
  - `GOCACHE=$PWD/.cache/go-build make vet`
  - `GOCACHE=$PWD/.cache/go-build GOFLAGS=-buildvcs=false make build`

## Self-Review

- Spec coverage: Covers config-file-only opt-in, command gate, documented delete endpoints, docs, and verification.
- Placeholder scan: No unresolved placeholders.
- Scope check: Only documented Chartbrew delete endpoints are added.
