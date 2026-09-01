# 🦘 token-hop (`thop`)

**[English](README.md)** | [한국어](README_KO.md)

> **"Hit a token limit? Just Hop across your AI coding agents!"** 💨  
> *"Token 제한에 걸리셨나요? 가볍게 Hop!"*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](file:///home/ubuntu/dev/yunkon-kim/token-hop/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](file:///home/ubuntu/dev/yunkon-kim/token-hop/go.mod)
[![Architecture: UA-IR](https://img.shields.io/badge/Architecture-UA--IR%20v1.0-brightgreen)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/design/token-hop-design-plan.md)

> ⚡ **Zero-Cost & Local-First**: `token-hop` runs 100% locally with **zero API keys and zero cost** by default. Generative AI (Claude/Gemini/OpenAI/Ollama) is purely optional for smart rule decomposition.

---

`token-hop` (`thop`) is a universal, ultra-fast **cross-agent prompt compiler and context synchronizer** designed to eliminate configuration fragmentation across modern AI coding assistants including **Claude Code, Cursor, GitHub Copilot, Google Antigravity, Gemini CLI, and Roo Code (Cline)**.

No more manual copy-pasting and reformatting rules when you hit token limits, usage caps, or switch between AI tools. With a single master **Single Source of Truth (SSOT)** (`AGENTS.md`), `token-hop` compiles and syncs all your rules, skills, workflows, and prompts for every AI IDE in **under 10ms**.

---

## ⚡ Why token-hop?

1. 🎯 **Single Source of Truth (SSOT)**  
   Manage all your guidelines in one `AGENTS.md` file and automatically compile them to Antigravity (`.agents/`), Cursor (`.cursor/rules/`), Copilot (`.github/`), and Claude (`CLAUDE.md`).

2. 🚀 **Zero-Cost & Ultra-Fast Deterministic Core**  
   100% rule-based AST parsing and code emission running locally in `<10ms` without requiring any external LLM API keys or incurring token costs.

3. 🧩 **Context Budgeting & JIT Skill Decomposition**  
   Automatically detects context window constraints (e.g. Cursor 6,000 chars, Antigravity 12,000 chars) and decomposes oversized rules into on-demand JIT `SKILL.md` packages—**slashing context token consumption by 80%+**.

4. 🔄 **1-Second Project Migration (`thop convert`)**  
   Instantly convert legacy GitHub Copilot setups (`.github/instructions`, `.github/prompts`) into modern Google Antigravity workflows and rules with a single command.

5. 🤖 **Selective Generative AI Augmentation (Hybrid Architecture)**  
   Optionally connect Claude, Gemini, OpenAI, or **local offline LLMs (Ollama)** for intelligent semantic rule decomposition, smart trigger generation, and 3-way merge conflict resolution.

---

## ⌨️ Quick Start

### 1. Installation

Install `token-hop` in seconds:

```bash
# Using Go (Recommended)
go install github.com/yunkon-kim/token-hop/cmd/token-hop@latest

# Or download pre-built binary via installer script (Linux / macOS)
curl -fsSL https://raw.githubusercontent.com/yunkon-kim/token-hop/main/install.sh | bash

# Set alias for ultra-fast typing (optional)
alias thop=token-hop
```

### 2. Instant Usage

Run inside any of your projects:

```bash
# 1) Convert existing GitHub Copilot (.github/) setup to Google Antigravity (.agents/)
thop convert --from copilot --to antigravity

# 2) Semantically decompose oversized rules into JIT Skills using AI (e.g. Gemini / Claude / Ollama)
thop convert --from copilot --to antigravity --decompose --ai --provider gemini

# 3) Audit token budget & character limits across all rules
thop audit

# 4) Compile SSOT (AGENTS.md) to all AI targets (Antigravity, Cursor, Copilot)
thop compile

# 5) Initialize a new project with SSOT scaffolding
thop init
```

---

### 🛠️ Building from Source (Contributors)

```bash
git clone https://github.com/yunkon-kim/token-hop.git
cd token-hop
make build && make test
```

---

## 📊 Cross-Agent Dimensional Mapping Matrix

| Entity | UA-IR Normalized Role | Google Antigravity | Cursor AI | GitHub Copilot | Claude Code |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Rule** | File-scoped static constraints | `.agents/rules/*.md`<br>(Glob / Always On) | `.cursor/rules/*.mdc`<br>(`globs: []`, `alwaysApply`) | `.github/instructions/*.md`<br>(`applyTo: glob`) | `./CLAUDE.md`<br>Directory cascade |
| **Skill** | On-demand JIT dynamic package | `.agent/skills/<name>/`<br>`SKILL.md` (JIT Load) | `.cursor/rules/*.mdc`<br>(Apply Intelligently) | `.github/skills/<name>/`<br>`SKILL.md` | `.claude/skills/<name>/`<br>`SKILL.md` |
| **Workflow** | Multi-step chained prompt loop | `.agents/workflows/*.md`<br>(`/workflow-name`) | `@rule` chained prompt | `.github/prompts/*.prompt.md`<br>(`/scaffold`) | CLI lifecycle hooks |
| **Instruction** | Global persona & system guidelines | `AGENTS.md` | `.cursor/rules/base.mdc`<br>(`alwaysApply: true`) | `.github/copilot-instructions.md` | `~/.claude/CLAUDE.md`<br>+ `./CLAUDE.md` |

---

## 📚 Documentation

- [Table of Contents (docs/README.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/README.md)
- [Architecture & Design Plan Specification (docs/design/token-hop-design-plan.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/design/token-hop-design-plan.md)
- [Test Suite & Real-world Fixtures Guide (test/README.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/test/README.md)

---

## 📄 License

[MIT License](file:///home/ubuntu/dev/yunkon-kim/token-hop/LICENSE) © 2026 Yunkon Kim