# Agent-Parkour Test Suite & Fixtures

This directory contains unit tests and integration test fixtures verifying cross-agent conversions against real-world repositories for `agent-parkour` (`parkour` / `pk`).

## Directory Structure

- **`fixtures/cm-beetle/.github/`**:
  - Real GitHub Copilot configuration files from the open-source multi-cloud project [`cloud-barista/cm-beetle`](https://github.com/cloud-barista/cm-beetle).
  - `copilot-instructions.md` (Repository-wide global instructions)
  - `instructions/` (6 rule files: `go`, `analyzer`, `markdown`, `tb-sync`, `transx`, `ui`)
  - `prompts/` (5 prompt files: `api-guide`, `git-commit`, `release-staging`, `sync-tb`, `sync-tb-model`)
- **`cm_beetle_conversion_test.go`**:
  - Comprehensive integration test: Parse `cm-beetle` Copilot configuration $\rightarrow$ compile to Google Antigravity (`.agents/rules/`, `.agents/workflows/`, `AGENTS.md`) $\rightarrow$ compile to Cursor (`.cursor/rules/*.mdc`) $\rightarrow$ verify Token & Context Window audit metrics.

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run the cm-beetle conversion integration test
go test -v -run TestCmBeetleCopilotToAntigravityConversion ./test

# Direct CLI binary conversion test
./parkour convert --from copilot --to antigravity --input test/fixtures/cm-beetle/.github --output /tmp/test-antigravity
./parkour audit --input test/fixtures/cm-beetle/.github
```
