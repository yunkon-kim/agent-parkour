package refiner

import (
	"strings"
	"text/template"
)

const antigravityTemplate = `---
description: Optimize and refine converted rules and workflows to strictly adhere to Antigravity architectural guidelines.
---

# Google Antigravity Configuration Optimizer & Refiner (/refine-context)

> [!IMPORTANT]
> **Review Required Before Execution**:
> Always inspect the converted files (e.g., git status, git diff) and verify changes before applying AI optimization.

You are an expert Google Antigravity agent configuration optimizer.
Please refactor the converted file(s) below to strictly adhere to Google Antigravity's native architecture, UI limits, and cognitive model.

## 📏 Antigravity Hard Constraints & Specifications (Must Satisfy)

1. **Frontmatter ` + "`" + `description:` + "`" + ` Limit**:
   - **Strict Maximum: 250 characters** (e.g., ` + "`" + `description: Generate conventional commits by analyzing staged changes` + "`" + `).
   - Must be a high-density, concise summary explaining *what* it does and *when* the agent should invoke it (required for Antigravity Progressive Disclosure).

2. **Content Body Length & Density**:
   - **Hard UI Limit: 12,000 characters**.
   - **Recommended Target: < 3,500 characters (~700 tokens)** for Workflows, **< 1,800 characters (~400 tokens)** for Rules.
   - Aggressively trim verbose text, duplicated examples, and conversational fluff.

3. **Cognitive Model Alignment (Zero Roleplay Preamble)**:
   - **DO NOT** include roleplay blocks like ` + "`" + `## Role\n You are an expert...` + "`" + ` or ` + "`" + `You are specialized in...` + "`" + `. Antigravity agents already possess context; start directly with the title and actionable steps.

4. **Native Tooling & ` + "`" + `// turbo` + "`" + ` Execution Annotations**:
   - Replace foreign tool bindings (e.g., ` + "`" + `get_changed_files` + "`" + `, ` + "`" + `run_in_terminal` + "`" + `, ` + "`" + `read_file` + "`" + `) with standard Antigravity terminal actions.
   - Add ` + "`" + `// turbo` + "`" + ` immediately on the line preceding automated shell commands to enable one-click execution:
     ` + "```markdown" + `
     // turbo
     git diff --cached --name-status
     ` + "```" + `

5. **Antigravity Formatting Standards**:
   - Structure workflows with numbered sequential steps (e.g., ` + "`" + `1. **Step Name**` + "`" + `).
   - Use GitHub-style alerts (` + "`" + `> [!NOTE]` + "`" + `, ` + "`" + `> [!TIP]` + "`" + `, ` + "`" + `> [!IMPORTANT]` + "`" + `, ` + "`" + `> [!WARNING]` + "`" + `).
   - Format file references as clickable links (` + "`" + `[filename.go](file:///path/to/filename.go)` + "`" + `).

6. **Progressive Skill Modularization**:
   - If a rule or workflow contains extensive reference manuals or multi-stage sub-procedures, isolate them into on-demand JIT Skills (` + "`" + `.agents/skills/<name>/SKILL.md` + "`" + `).

{{if .CustomGuidance}}
## 💡 Additional Project Guidance
{{.CustomGuidance}}
{{end}}

## 📋 Converted Input to Refine (Total: {{.TotalFiles}} file(s), ~{{.TotalTokens}} tokens)

{{range .Files}}
### 📄 Target File: ` + "`" + `{{.FilePath}}` + "`" + ` (Type: **{{.EntityType}}**, Tokens: ~{{.Tokens}}, Characters: {{.Characters}})
{{if .IsOversized}}⚠️ *Warning: Exceeds recommended token threshold (>400 tokens). Aggressively condense into high-density Antigravity format.*{{end}}

` + "```markdown" + `
{{.FileContent}}
` + "```" + `

---
{{end}}

## 📤 Output Instructions (Copy-Paste Ready)
For each target file above, produce a complete drop-in replacement markdown file:
1. Include valid YAML frontmatter with ` + "`" + `description:` + "`" + ` ($\le 250$ chars).
2. Follow with the refined content body ($\le 12,000$ chars, target $< 3,500$ chars for workflows, $< 1,800$ chars for rules).
3. Ensure ` + "`" + `// turbo` + "`" + ` annotations are present before terminal commands.
4. Output the result in a clean markdown code block ready for one-click copy & paste.
5. At the conclusion of your response, provide the user with a concise Next Steps checklist:
   - [Review] Check diff with git diff .agents/
   - [Verify] Audit token reduction with pk audit -i .agents/
   - [Clean]  Remove temporary workflow: rm .agents/workflows/refine-context.md
   - [Commit] Commit changes to repository

> [!TIP]
> After applying optimizations, you can safely remove this temporary workflow:
> // turbo
> rm .agents/workflows/refine-context.md
`

const cursorTemplate = `# [Prompt] Cursor AI MDC Rule Optimizer & Refiner

You are an expert in configuring Cursor AI (.cursor/rules/*.mdc) environments.
Please review and refine the converted rule file(s) to adhere to modern Cursor MDC standards and optimize context relevance.

## 📏 Cursor MDC Hard Constraints & Guidelines

1. **MDC Frontmatter Schema**:
   - ` + "`" + `description:` + "`" + `: Short summary ($\le 200$ chars).
   - ` + "`" + `globs:` + "`" + `: Precise file patterns (e.g., ` + "`" + `["**/*.go", "pkg/**/*.go"]` + "`" + `).
   - ` + "`" + `alwaysApply:` + "`" + `: Set to ` + "`" + `false` + "`" + ` for language/framework rules; set to ` + "`" + `true` + "`" + ` only for global architectural rules.

2. **Context Referencing & Modularity**:
   - Utilize ` + "`" + `@file` + "`" + ` or ` + "`" + `@rule` + "`" + ` referencing syntax to avoid redundant duplication across rules.
   - Keep rules concise, imperative, and under **6,000 characters**.

3. **Tool & Agent Compatibility**:
   - Replace foreign agent tool names with Cursor Composer / Agent native instructions.

{{if .CustomGuidance}}
## 💡 Additional Project Guidance
{{.CustomGuidance}}
{{end}}

## 📋 Converted Input to Refine (Total: {{.TotalFiles}} file(s), ~{{.TotalTokens}} tokens)

{{range .Files}}
### 📄 Target File: ` + "`" + `{{.FilePath}}` + "`" + ` (Type: **{{.EntityType}}**, Tokens: ~{{.Tokens}}, Characters: {{.Characters}})

` + "```markdown" + `
{{.FileContent}}
` + "```" + `

---
{{end}}

## 📤 Output Instructions (Copy-Paste Ready)
Output the complete drop-in replacement ` + "`" + `.mdc` + "`" + ` file with valid frontmatter in a clean code block ready to copy and paste.
`

const claudeTemplate = `# [Prompt] Claude Code (CLAUDE.md) Optimizer & Refiner

You are an expert in configuring Anthropic Claude Code CLI environments.
Please review and consolidate project instructions into a streamlined, high-density ` + "`" + `CLAUDE.md` + "`" + `.

## 📏 Claude Code Hard Constraints & Guidelines

1. **Conciseness & High Density**:
   - Keep ` + "`" + `CLAUDE.md` + "`" + ` strictly under **200 lines** total. Claude Code reads this on every turn.
   - Structure sections cleanly: Build/Lint/Test Commands, Code Style, Architecture, and Critical Invariants.

2. **Actionable CLI Commands**:
   - Provide exact, runnable terminal commands (e.g., ` + "`" + `go build ./...` + "`" + `, ` + "`" + `npm test` + "`" + `).

3. **Sub-Agent Skills Extraction**:
   - Extract large, specialized task procedures into modular skill files under ` + "`" + `.claude/skills/<name>/SKILL.md` + "`" + `.

{{if .CustomGuidance}}
## 💡 Additional Project Guidance
{{.CustomGuidance}}
{{end}}

## 📋 Converted Input to Refine (Total: {{.TotalFiles}} file(s), ~{{.TotalTokens}} tokens)

{{range .Files}}
### 📄 Target File: ` + "`" + `{{.FilePath}}` + "`" + ` (Type: **{{.EntityType}}**, Tokens: ~{{.Tokens}}, Characters: {{.Characters}})

` + "```markdown" + `
{{.FileContent}}
` + "```" + `

---
{{end}}

## 📤 Output Instructions (Copy-Paste Ready)
Output the complete, optimized ` + "`" + `CLAUDE.md` + "`" + ` or skill file in a clean code block ready to copy and paste.
`

const copilotTemplate = `# [Prompt] GitHub Copilot Instruction & Prompt Optimizer

You are an expert in configuring GitHub Copilot repository guidelines (.github/) environments.
Please review and refine the converted configuration files.

## 🎯 Target Optimization Objectives

1. **Repository Instructions (` + "`" + `.github/copilot-instructions.md` + "`" + `)**:
   - Provide overarching repository architecture and coding conventions.
2. **Targeted File Instructions (` + "`" + `.github/instructions/*.instructions.md` + "`" + `)**:
   - Ensure ` + "`" + `applyTo: "<glob>"` + "`" + ` frontmatter is correctly scoped.
3. **Reusable Prompts (` + "`" + `.github/prompts/*.prompt.md` + "`" + `)**:
   - Utilize parameterized placeholders and clear argument hints.

{{if .CustomGuidance}}
## 💡 Additional Project Guidance
{{.CustomGuidance}}
{{end}}

## 📋 Converted Files to Refine (Total: {{.TotalFiles}} file(s), ~{{.TotalTokens}} tokens)

{{range .Files}}
### 📄 Target: ` + "`" + `{{.FilePath}}` + "`" + ` (Type: **{{.EntityType}}**, ~{{.Tokens}} tokens)

` + "```markdown" + `
{{.FileContent}}
` + "```" + `

---
{{end}}

## 📤 Output Instructions
Please output the refined GitHub Copilot configuration files.
`

var parsedTemplates = map[string]*template.Template{}

func init() {
	parsedTemplates["antigravity"] = template.Must(template.New("antigravity").Parse(antigravityTemplate))
	parsedTemplates["cursor"] = template.Must(template.New("cursor").Parse(cursorTemplate))
	parsedTemplates["claude"] = template.Must(template.New("claude").Parse(claudeTemplate))
	parsedTemplates["copilot"] = template.Must(template.New("copilot").Parse(copilotTemplate))
}

// GetTemplate returns the parsed template for a target platform
func GetTemplate(targetPlatform string) *template.Template {
	key := strings.ToLower(strings.TrimSpace(targetPlatform))
	switch key {
	case "antigravity", "gemini", "google":
		return parsedTemplates["antigravity"]
	case "cursor", "cursor-ai":
		return parsedTemplates["cursor"]
	case "claude", "claude-code", "anthropic":
		return parsedTemplates["claude"]
	case "copilot", "github", "gh-copilot":
		return parsedTemplates["copilot"]
	default:
		return parsedTemplates["antigravity"]
	}
}
