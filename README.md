# 🏃 agent-parkour (`parkour` / `pk`)

**[English](README.md)** | [한국어](README_KO.md)

> **"Hit a token wall? Just Parkour across your AI coding agents!"** 🏃💨

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](go.mod)
[![Architecture: UA-IR](https://img.shields.io/badge/Architecture-UA--IR%20v1.0-brightgreen)](docs/design/agent-parkour-design-plan.md)

> ⚡ **Zero-Cost & Local-First**: `agent-parkour` runs 100% locally with **zero API keys and zero cost** by default. Generative AI is purely an optional, experimental feature for advanced rule decomposition.

---

`agent-parkour` (`parkour` / `pk`) is a universal, ultra-fast **cross-agent prompt compiler and context synchronizer** designed to eliminate configuration fragmentation across modern AI coding assistants including **Claude Code, Cursor, GitHub Copilot, Google Antigravity, Gemini CLI, and Roo Code (Cline)**.

No more manual copy-pasting and reformatting rules when you hit token limits, usage caps, or switch between AI tools. With a single master **Single Source of Truth (SSOT)** (`AGENTS.md`), `agent-parkour` compiles and syncs all your rules, skills, workflows, and prompts for every AI IDE in **under 10ms**.

---

## ⚡ Why agent-parkour?

1. 🎯 **Single Source of Truth (SSOT)**  
   Manage all your guidelines in one `AGENTS.md` file and automatically synchronize them to Antigravity (`.agents/`), Cursor (`.cursor/rules/`), Copilot (`.github/`), and Claude (`CLAUDE.md`).

2. 🚀 **Zero-Cost & Ultra-Fast Deterministic Core**  
   100% rule-based AST parsing and code emission running locally in `<10ms` without requiring external LLM API keys or incurring token costs.

3. 🛡️ **Automatic Timestamped Backups (`*.bak_YYYYMMDD_HHMMSS`)**  
   Safely modifies rules with instant timestamped backups whenever existing files are overwritten, ensuring you never lose previous work.

4. 🧩 **Context Optimization & JIT Skill Decomposition**  
   Automatically detects context limits (e.g. Cursor 6,000 chars, Antigravity 12,000 chars) and decomposes oversized rules into on-demand JIT `SKILL.md` packages—**slashing context token consumption by 80%+**.

5. 🔄 **1-Second Multi-Agent Migration (`parkour convert`)**  
   Instantly convert legacy GitHub Copilot setups (`.github/instructions`, `.github/prompts`) into modern Google Antigravity workflows and rules with a single command.

---

## 📊 Cross-Agent Dimensional Mapping Matrix

`agent-parkour` adopts a **Directory-First architecture**, preserving the native directory structure and ergonomics of each platform while performing bidirectional compilation:

| Entity | Role & Behavior | Google Antigravity | GitHub Copilot | Claude Code | Cursor AI |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Instruction** | Global SSOT master guidelines | `AGENTS.md`<br>*(or `.agents/AGENTS.md`)* | `.github/copilot-instructions.md` | `CLAUDE.md`<br>*(Root global)* | `.cursorrules`<br>*(or `base.mdc`)* |
| **Rule** | Scoped coding standards & constraints | `.agents/rules/*.md`<br>(`globs`) | `.github/instructions/*.md`<br>(`applyTo` string or array) | Subdirectory `CLAUDE.md`<br>(Path-scoped) | `.cursor/rules/*.mdc`<br>(`globs`) |
| **Prompt** | Reusable `/<name>` prompt templates | `.agents/workflows/*.md` | `.github/prompts/*.prompt.md`<br>(`/<cmd>`) | `.claude/workflows/*.md`<br>(`/<cmd>`) | `.cursor/rules/*.mdc` |
| **Workflow** | Multi-step agentic execution & verification | `.agents/workflows/*.md`<br>(`/<cmd>`) | `.github/prompts/*.prompt.md` | `.claude/workflows/*.md` | `.cursor/rules/*.mdc` |
| **Skill** | Dynamic JIT knowledge & tool package | `.agents/skills/<name>/`<br>`SKILL.md` (JIT) | `.github/skills/<name>/`<br>`SKILL.md` (JIT) | `.claude/skills/<name>/`<br>`SKILL.md` (JIT) | `.cursor/rules/*.mdc`<br>(Description-based) |
| **Agent** | Specialized custom sub-agent persona | `.agents/skills/` (or Agent) | `.github/agents/*.agent.md` | `.claude/agents/` | `.cursor/rules/*.mdc` |

---

### 📖 Classification Guide

When running `parkour describe`, each entity is clearly annotated so readers immediately understand its purpose:

- **Instruction**: Project-wide master guidelines loaded on every conversation turn
- **Rule**: File/language-scoped coding constraints loaded automatically when editing matching files
- **Prompt**: Reusable instruction templates executed via `/<name>` slash commands
- **Workflow**: Multi-step agentic execution procedures and verification loops executed via `/<name>`
- **Skill**: Specialized knowledge manuals and toolkits loaded dynamically (JIT) when needed
- **Agent**: Dedicated custom sub-agent personas with custom roles and tool bindings

> 💡 **Context Token Impact**:
> - **Always-On (Turn-0)**: Fixed overhead injected into every conversation turn
> - **On-Demand (JIT)**: Dynamically loaded only when matching files or commands are triggered, saving token budget

---

## ⌨️ Quick Start

### 1. Installation

Install `agent-parkour` in seconds:

```bash
# Using Go (Recommended)
go install github.com/yunkon-kim/agent-parkour/cmd/parkour@latest

# Or download pre-built binary via installer script (Linux / macOS)
curl -fsSL https://raw.githubusercontent.com/yunkon-kim/agent-parkour/main/install.sh | bash

# Set alias for ultra-fast typing (optional - automatically configured by installer)
alias pk=parkour
```

### 2. Instant Usage (Deterministic Core)

Run inside any of your projects with zero configuration (you can use `parkour` or the short alias `pk`):

```bash
# 1) Preview transformation mapping plan in tabular format before converting
parkour describe --from copilot --to antigravity

# 2) Preview mapping plan and export directly to Markdown or JSON file
parkour describe --from copilot --to antigravity --out plan.md
parkour describe --from copilot --to antigravity --out plan.json

# 3) View cross-agent specification matrix (locations, syntax, and behaviors)
parkour describe --from copilot --to antigravity --spec

# 4) Convert from GitHub Copilot to Google Antigravity (auto-maps all rules & prompts)
parkour convert --from copilot --to antigravity

# 5) Convert and generate 2nd-stage AI refinement workflow (/refine-context)
parkour convert --from copilot --to antigravity --gen-refine-prompt
# ⚠️ Note: Always review converted files (`git diff`) before triggering /refine-context in chat!

# 6) Convert with dry-run preview (simulates without writing files)
parkour convert --from copilot --to antigravity --dry-run

# 7) Convert from GitHub Copilot to all supported targets (Antigravity, Cursor)
parkour convert --from copilot --to all

# 8) Convert to a specific target (auto-detects source in repo)
parkour convert --to cursor

# 9) Audit token counts & context limits across all rules
parkour audit

# 10) Initialize a new project with SSOT scaffolding (AGENTS.md & parkour.yaml)
parkour init
```

---

## 🧪 Optional & Experimental: Generative AI Augmentation

> ⚠️ **Note**: Generative AI integration is **strictly optional and experimental**. `agent-parkour`'s core compiler runs 100% deterministically without requiring any API keys.

If you wish to enable advanced semantic rule decomposition (analyzing complex rules and splitting them into modular JIT skills), you can selectively enable Generative AI:

### 1. Configure `.env` (Optional)
```bash
cp .env.example .env
# Edit .env to set your preferred provider:
# PARKOUR_AI_ENABLED=true
# PARKOUR_AI_PROVIDER=gemini  # [gemini | claude | openai | ollama]
# GEMINI_API_KEY=your_key_here
```

### 2. Run Semantic Decomposition
```bash
# Decompose oversized rules using Google Gemini
parkour convert --from copilot --to antigravity --decompose --ai --provider gemini

# Or using local offline LLM via Ollama (Zero API cost)
parkour convert --from copilot --to antigravity --decompose --ai --provider ollama
```

---

### 🛠️ Building from Source (Contributors)

```bash
git clone https://github.com/yunkon-kim/agent-parkour.git
cd agent-parkour
make build && make test
```

---

## 📚 Documentation

- [Table of Contents (docs/README.md)](docs/README.md)
- [Architecture & Design Plan Specification (docs/design/agent-parkour-design-plan.md)](docs/design/agent-parkour-design-plan.md)
- [Maintenance & Spec Evolution Prompt (docs/maintenance/SPEC_EVOLUTION_PROMPT.md)](docs/maintenance/SPEC_EVOLUTION_PROMPT.md)
- [Test Suite & Real-world Fixtures Guide (test/README.md)](test/README.md)

---

## 📄 License

[MIT License](LICENSE) © 2026 Yunkon Kim