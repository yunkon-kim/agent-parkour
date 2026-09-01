# 🦘 token-hop (`thop`)

**[English](README.md)** | [한국어](README_KO.md)

> **"Hit a token limit? Just Hop across your AI coding agents!"** 💨  
> *"Token 제한에 걸리셨나요? 가볍게 Hop!"*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](file:///home/ubuntu/dev/yunkon-kim/token-hop/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](file:///home/ubuntu/dev/yunkon-kim/token-hop/go.mod)
[![Architecture: UA-IR](https://img.shields.io/badge/Architecture-UA--IR%20v1.0-brightgreen)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/design/token-hop-design-plan.md)

> ⚡ **Zero-Cost & Local-First**: `token-hop` runs 100% locally with **zero API keys and zero cost** by default. Generative AI is purely an optional, experimental feature for advanced rule decomposition.

---

`token-hop` (`thop`) is a universal, ultra-fast **cross-agent prompt compiler and context synchronizer** designed to eliminate configuration fragmentation across modern AI coding assistants including **Claude Code, Cursor, GitHub Copilot, Google Antigravity, Gemini CLI, and Roo Code (Cline)**.

No more manual copy-pasting and reformatting rules when you hit token limits, usage caps, or switch between AI tools. With a single master **Single Source of Truth (SSOT)** (`AGENTS.md`), `token-hop` compiles and syncs all your rules, skills, workflows, and prompts for every AI IDE in **under 10ms**.

---

## ⚡ Why token-hop?

1. 🎯 **Single Source of Truth (SSOT)**  
   Manage all your guidelines in one `AGENTS.md` file and automatically synchronize them to Antigravity (`.agents/`), Cursor (`.cursor/rules/`), Copilot (`.github/`), and Claude (`CLAUDE.md`).

2. 🚀 **Zero-Cost & Ultra-Fast Deterministic Core**  
   100% rule-based AST parsing and code emission running locally in `<10ms` without requiring external LLM API keys or incurring token costs.

3. 🛡️ **Automatic Timestamped Backups (`*.bak_YYYYMMDD_HHMMSS`)**  
   Safely modifies rules with instant timestamped backups whenever existing files are overwritten, ensuring you never lose previous work.

4. 🧩 **Context Optimization & JIT Skill Decomposition**  
   Automatically detects context limits (e.g. Cursor 6,000 chars, Antigravity 12,000 chars) and decomposes oversized rules into on-demand JIT `SKILL.md` packages—**slashing context token consumption by 80%+**.

5. 🔄 **1-Second Multi-Agent Migration (`thop convert`)**  
   Instantly convert legacy GitHub Copilot setups (`.github/instructions`, `.github/prompts`) into modern Google Antigravity workflows and rules with a single command.

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

### 2. Instant Usage (Deterministic Core)

Run inside any of your projects with zero configuration:

```bash
# 1) Convert from GitHub Copilot to Google Antigravity (auto-maps all rules & prompts)
thop convert --from copilot --to antigravity

# 2) Convert from GitHub Copilot to all supported targets (Antigravity, Cursor)
thop convert --from copilot --to all

# 3) Convert to a specific target (auto-detects source in repo)
thop convert --to cursor

# 4) Audit token counts & context limits across all rules
thop audit

# 5) Initialize a new project with SSOT scaffolding (AGENTS.md)
thop init
```

---

## 🧪 Optional & Experimental: Generative AI Augmentation

> ⚠️ **Note**: Generative AI integration is **strictly optional and experimental**. `token-hop`'s core compiler runs 100% deterministically without requiring any API keys.

If you wish to enable advanced semantic rule decomposition (analyzing complex rules and splitting them into modular JIT skills), you can selectively enable Generative AI:

### 1. Configure `.env` (Optional)
```bash
cp .env.example .env
# Edit .env to set your preferred provider:
# TOKEN_HOP_AI_ENABLED=true
# TOKEN_HOP_AI_PROVIDER=gemini  # [gemini | claude | openai | ollama]
# GEMINI_API_KEY=your_key_here
```

### 2. Run Semantic Decomposition
```bash
# Decompose oversized rules using Google Gemini
thop convert --from copilot --to antigravity --decompose --ai --provider gemini

# Or using local offline LLM via Ollama (Zero API cost)
thop convert --from copilot --to antigravity --decompose --ai --provider ollama
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
| **Skill** | On-demand JIT dynamic package | `.agents/skills/<name>/`<br>`SKILL.md` (JIT Load) | `.cursor/rules/*.mdc`<br>(Apply Intelligently) | `.github/skills/<name>/`<br>`SKILL.md` | `.claude/skills/<name>/`<br>`SKILL.md` |
| **Workflow** | Multi-step chained prompt loop | `.agents/workflows/*.md`<br>(`/workflow-name`) | `@rule` chained prompt | `.github/prompts/*.prompt.md`<br>(`/scaffold`) | CLI lifecycle hooks |
| **Instruction** | Global persona & system guidelines | `AGENTS.md` | `.cursor/rules/base.mdc`<br>(`alwaysApply: true`) | `.github/copilot-instructions.md` | `~/.claude/CLAUDE.md`<br>+ `./CLAUDE.md` |

---

## 📚 Documentation

- [Table of Contents (docs/README.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/README.md)
- [Architecture & Design Plan Specification (docs/design/token-hop-design-plan.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/design/token-hop-design-plan.md)
- [Maintenance & Spec Evolution Prompt (docs/maintenance/SPEC_EVOLUTION_PROMPT.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/maintenance/SPEC_EVOLUTION_PROMPT.md)
- [Test Suite & Real-world Fixtures Guide (test/README.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/test/README.md)

---

## 📄 License

[MIT License](file:///home/ubuntu/dev/yunkon-kim/token-hop/LICENSE) © 2026 Yunkon Kim