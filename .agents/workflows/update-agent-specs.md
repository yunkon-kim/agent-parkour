---
description: Automated workflow to update token-hop parsers, emitters, and documentation when AI assistant specs or directory paths change.
---

# Workflow: Update Agent Specifications & Add Targets

This workflow automates updating `token-hop` when AI coding environments (Antigravity, Cursor, Copilot, Claude, Windsurf, etc.) update their configuration paths, file formats, or frontmatter schemas.

## Step 1: Analyze Spec Changes
1. Identify the target tool name and new/changed specifications:
   - Config directory path & file extensions
   - YAML frontmatter schema fields
   - Context window / token / character limitations
   - Activation trigger modes (Glob, AlwaysOn, JIT, Slash command)

## Step 2: Update Core Engine
1. **UA-IR AST (`pkg/ir/types.go`)**: Update entity types or metadata fields if needed.
2. **Parser (`pkg/parser/`)**: Update or create parser for the target tool.
3. **Emitter (`pkg/emitter/`)**: Update or create code generator with proper frontmatter formatting.
4. **Engine & CLI (`pkg/engine/engine.go`, `cmd/token-hop/main.go`)**: Register the format in CLI options.

## Step 3: Test Fixtures & Verification
1. Add/update test fixtures in `test/fixtures/<tool_name>/`.
2. Run `go test -v ./...` to verify zero regressions.
3. Rebuild binary with `make build` and test conversion.

## Step 4: Synchronize Documentation
1. Update `docs/design/token-hop-design-plan.md` mapping matrix.
2. Update `README.md` & `README_KO.md`.

## Official Reference Sources
- **Google Antigravity**: https://cloud.google.com/gemini/docs/codeassist
- **Cursor AI Rules**: https://docs.cursor.com/context/rules-for-ai
- **GitHub Copilot Instructions**: https://docs.github.com/en/copilot/customizing-copilot/adding-custom-instructions-for-github-copilot
- **Claude Code**: https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview
- **Windsurf Rules**: https://docs.codeium.com/windsurf/memories
- **Roo Code / Cline**: https://docs.roocode.com

