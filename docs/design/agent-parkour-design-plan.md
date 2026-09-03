# Agent-Parkour Engine: Cross-Agent Prompt Compiler Architecture & Plan
**Cross-Agent Prompt Compiler Design and Architecture for Resolving Configuration Fragmentation Across AI Development Environments**

[English](agent-parkour-design-plan.md) | [한국어](agent-parkour-design-plan_KO.md)

- **Project Name**: `agent-parkour`
- **CLI Binary**: `parkour` (Alias: `pk`)
- **Version**: v1.0.0
- **Author**: Yunkon Kim

---

## 1. Background & Problem Definition

### 1.1 Multi-Agent IDE Context Fragmentation
Modern software engineering increasingly employs multiple AI coding agents in parallel—such as **Google Antigravity, GitHub Copilot, Claude Code, and Cursor AI**. Developers switch between these tools frequently based on model strengths, task domains (UI design, backend logic, CLI scripting), and vendor rate/token limits.

However, each environment suffers from severe configuration fragmentation across rules, workflows, skills, instructions, and prompt templates:

```
+-------------------------------------------------------------------------------+
|                       Multi-Agent Ecosystem Fragmentation                     |
+-------------------------------------------------------------------------------+
| [Google Antigravity] | .agents/rules/*.md, .agents/workflows/, .agents/skills/ |
| [GitHub Copilot]     | .github/copilot-instructions.md, prompts/, skills/     |
| [Claude Code]        | ./CLAUDE.md, ~/.claude/CLAUDE.md, .claude/skills/      |
| [Cursor AI]          | .cursor/rules/*.mdc (globs, alwaysApply, description)  |
+-------------------------------------------------------------------------------+
```

### 1.2 Core Issues
1. **Context Rot & Drift**: High-quality constraints optimized in one tool (e.g. Cursor `.mdc`) are lost or manually mistranslated when switching to another agent (e.g. Antigravity or Claude Code).
2. **Context Window Overflow**: Oversized instruction files (>500 lines or >6,000 characters) bloat Turn-0 overhead and severely degrade model reasoning compliance.
3. **Absence of Bidirectional Sync**: Rule changes made in one IDE diverge from project standards without an automated synchronization bridge.

### 1.3 Glossary & Key Concepts
- **SSOT (Single Source of Truth)**: A centralized master document (`AGENTS.md`) from which all target configurations are deterministically derived.
- **UA-IR (Universal Agent Intermediate Representation)**: A platform-neutral abstract data format acting as the compiler's internal AST representation.
- **AST (Abstract Syntax Tree)**: Hierarchical data structure modeling frontmatter metadata, markdown payloads, and glob rules.
- **JIT (Just-In-Time Dynamic Loading)**: On-demand loading of modular `SKILL.md` packages only when relevant tasks are invoked, minimizing turn-0 token overhead.
- **Directory-First Classification**: Classifying entity lifecycles and types by directory paths (`instructions/`, `rules/`, `prompts/`, `workflows/`, `skills/`) rather than ad-hoc metadata.

---

## 2. System Architecture Overview

`agent-parkour` operates on a deterministic 3-stage compiler pipeline: **Parser $\rightarrow$ UA-IR AST $\rightarrow$ Target Emitter**, accompanied by an optional AI-assisted decomposition module.

```
                       [ Single Source of Truth: AGENTS.md / UA-IR ]
                                              │
                                              ▼
                        ┌───────────────────────────────────────────┐
                        │        agent-parkour Lexer & Parser       │
                        │    (Frontmatter + AST Tree Builder)       │
                        └─────────────────────┬─────────────────────┘
                                              │
                                              ▼
                        ┌───────────────────────────────────────────┐
                        │      Universal Agent IR (UA-IR AST)       │
                        │ (Instruction | Rule | Prompt | Skill ...) │
                        └─────────────────────┬─────────────────────┘
                                              │
                   ┌──────────────────────────┴──────────────────────────┐
                   ▼                                                     ▼
    ┌─────────────────────────────┐                       ┌─────────────────────────────┐
    │  Context Budgeting Engine   │                       │   Bi-Directional 3-Way      │
    │ (Token/Char Slicing & JIT)  │                       │   Diff & Drift Sync Engine  │
    └──────────────┬──────────────┘                       └──────────────┬──────────────┘
                   └──────────────────────────┬──────────────────────────┘
                                              │
                                              ▼
                        ┌───────────────────────────────────────────┐
                        │        Target Emitter & Code Gen          │
                        └──────┬──────────┬──────────┬──────────┬───┘
                               │          │          │          │
         ┌─────────────────────┘          │          │          └─────────────────────┐
         ▼                                ▼          ▼                                ▼
┌──────────────────┐             ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│Google Antigravity│             │  GitHub Copilot  │ │   Claude Code    │ │    Cursor AI     │
│.agents/rules/*.md│             │.github/copilot-  │ │./CLAUDE.md       │ │.cursor/rules/*.md│
│.agents/workflows/│             │ instructions.md  │ │.claude/skills/   │ │(globs/always)    │
│.agents/skills/   │             │.github/prompts/  │ └──────────────────┘ └──────────────────┘
└──────────────────┘             │.github/skills/   │
                                 └──────────────────┘
```

---

## 3. Cross-Agent Dimensional Mapping Matrix

| Entity | Role & Behavior | Google Antigravity | GitHub Copilot | Claude Code | Cursor AI |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Instruction** | Global SSOT master guidelines | `AGENTS.md`<br>*(or `.agents/AGENTS.md`)* | `.github/copilot-instructions.md` | `CLAUDE.md`<br>*(Root global)* | `.cursorrules`<br>*(or `base.mdc`)* |
| **Rule** | Scoped coding standards & constraints | `.agents/rules/*.md`<br>(`globs`) | `.github/instructions/*.md`<br>(`applyTo` string or array) | Subdirectory `CLAUDE.md`<br>(Path-scoped) | `.cursor/rules/*.mdc`<br>(`globs`) |
| **Prompt** | Reusable `/<name>` prompt templates | `.agents/workflows/*.md` | `.github/prompts/*.prompt.md`<br>(`/<cmd>`) | `.claude/workflows/*.md`<br>(`/<cmd>`) | `.cursor/rules/*.mdc` |
| **Workflow** | Multi-step agentic execution & verification | `.agents/workflows/*.md`<br>(`/<cmd>`) | `.github/prompts/*.prompt.md` | `.claude/workflows/*.md` | `.cursor/rules/*.mdc` |
| **Skill** | Dynamic JIT knowledge & tool package | `.agents/skills/<name>/`<br>`SKILL.md` (JIT) | `.github/skills/<name>/`<br>`SKILL.md` (JIT) | `.claude/skills/<name>/`<br>`SKILL.md` (JIT) | `.cursor/rules/*.mdc`<br>(Description-based) |
| **Agent** | Specialized custom sub-agent persona | `.agents/skills/` (or Agent) | `.github/agents/*.agent.md` | `.claude/agents/` | `.cursor/rules/*.mdc` |

---

## 4. Universal Agent IR (UA-IR v1.0.0) Specification

### 4.1 Schema Definition
UA-IR documents contain YAML Frontmatter metadata paired with a standardized Markdown payload:

```yaml
ua_version: "1.0.0"
metadata:
  id: "rule-react-component-standards"
  type: "Rule" # [Instruction | Rule | Prompt | Workflow | Skill | Agent]
  name: "React Component Conventions"
  description: "Enforces strict TypeScript functional patterns and named exports for UI components."
  version: "1.0.0"
  author: "agent-parkour"
  tags: ["frontend", "react", "typescript"]

activation:
  mode: "Glob" # [AlwaysOn | Glob | ModelDecision | OnDemand]
  globs:
    - "src/components/**/*.tsx"
    - "src/ui/**/*.tsx"
  exclude_globs:
    - "**/*.test.tsx"
    - "**/*.stories.tsx"
  slash_command: "scaffold-component"

context_budget:
  priority: "High" # [Critical | High | Medium | Low]
  max_tokens: 400
  max_characters: 1800
  decompose_strategy: "JIT_Skill"

bindings:
  allowed_tools: ["read_file", "write_to_file", "run_command"]
  mcp_servers: []

payload:
  markdown_body: |
    # React Component Conventions
    - Use functional components with explicit type definitions for Props.
    - Always use named exports (`export const ComponentName = ...`) instead of default exports.
    - Colocate styles and unit tests in the same directory.
```

---

## 5. CLI Tool (`agent-parkour` / `pk`) Specification

### 5.1 Commands
```bash
# 1. Inspect configuration mapping and token impact
parkour describe -i /path/to/repo --from copilot --to antigravity [--out plan.md | --spec]

# 2. Direct conversion
parkour convert --from copilot --to antigravity [--dry-run] [--gen-refine-prompt]

# 3. Context budget audit
parkour audit --input .github/ --max-tokens 400

# 4. Initialize SSOT in a new project
parkour init
```

### 5.2 Project Configuration (`parkour.yaml`)
```yaml
version: "1.0"
ssot: "AGENTS.md"
targets:
  - name: "antigravity"
    output_dir: ".agents"
    enable_skills: true
  - name: "copilot"
    instructions_dir: ".github/instructions"
    prompts_dir: ".github/prompts"
    skills_dir: ".github/skills"
  - name: "claude"
    output_file: "CLAUDE.md"
    skills_dir: ".claude/skills"
  - name: "cursor"
    output_dir: ".cursor/rules"
context_budget:
  max_tokens_per_rule: 400
  max_characters_per_rule: 1800
  auto_decompose: true
ai:
  enabled: false
  provider: "gemini"
  model: "gemini-2.5-pro"
```

---

## 6. Official References

1. **Google Antigravity / Gemini Code Assist**:
   - [Google Cloud Gemini Code Assist Documentation](https://cloud.google.com/gemini/docs/codeassist)
2. **GitHub Copilot**:
   - [GitHub Copilot Custom Instructions](https://docs.github.com/en/copilot/customizing-copilot/adding-custom-instructions-for-github-copilot)
   - [GitHub Copilot Custom Prompts & Agent Skills](https://docs.github.com/en/copilot)
3. **Claude Code (Anthropic)**:
   - [Claude Code Overview & Tooling](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview)
   - [Claude Code Memory & CLAUDE.md Guide](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/memory)
4. **Cursor AI**:
   - [Cursor Rules Official Documentation](https://docs.cursor.com/context/rules-for-ai)
