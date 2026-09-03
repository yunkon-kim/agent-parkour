# Agent-Parkour Documentation

[English](README.md) | [한국어](README_KO.md)

This directory contains technical architecture, system design specifications, and reference documentation for the **Agent-Parkour Engine** (Cross-Agent Prompt Compiler).

## Table of Contents

### 1. Design & Architecture
- [docs/design/agent-parkour-design-plan.md](design/agent-parkour-design-plan.md): 
  - Key terminology & acronyms (SSOT, IR, AST, JIT, MDC, MCP, 3-Way Diff)
  - Multi-agent ecosystem fragmentation analysis (Antigravity, Copilot, Claude Code, Cursor, Gemini CLI, Roo Code)
  - 6 core entities (Instruction, Rule, Prompt, Workflow, Skill, Agent) dimensional translation model
  - Universal Agent IR (UA-IR v1.0.0) AST schema specification
  - Compiler pipeline (Parser $\rightarrow$ Context Budgeting & JIT Decomposer $\rightarrow$ Target Emitter)
  - 3-Way Diff bi-directional synchronization engine
  - CLI command specifications (`parkour init/convert/describe/prompt/audit`) and implementation roadmap

### 2. Maintenance & Spec Evolution
- [docs/maintenance/SPEC_EVOLUTION_PROMPT.md](maintenance/SPEC_EVOLUTION_PROMPT.md):
  - Systematic 6-step prompt template and checklist for updating `agent-parkour` when IDE/LLM tools change directory structures, file extensions, or frontmatter schemas
  - Google Antigravity workflow: `.agents/workflows/update-agent-specs.md` (`/update-agent-specs`)

### 3. Specification Verification
- [docs/references/cross-agent-spec-verification.md](references/cross-agent-spec-verification.md):
  - Official vendor citations, hard character limits (250 / 12,000 chars), and live verification reports across Antigravity, Copilot, Claude Code, and Cursor
