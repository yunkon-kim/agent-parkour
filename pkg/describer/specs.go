package describer

import (
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/ir"
)

var platformSpecs = map[string]PlatformSpec{
	"antigravity": {
		ID:          "antigravity",
		Name:        "Google Antigravity",
		RootDir:     ".agents/",
		SSOTFile:    "AGENTS.md",
		Description: "Google Antigravity multi-agent workspace with rules, workflows, and JIT skills",
		Entities: []EntitySpec{
			{
				EntityType:  ir.TypeInstruction,
				Location:    "AGENTS.md (Root) / ~/.gemini/GEMINI.md",
				Syntax:      "Plain Markdown (SSOT Guidelines)",
				TriggerMode: "Always-On",
				Description: "Baseline project instructions loaded on session start",
			},
			{
				EntityType:  ir.TypeRule,
				Location:    ".agents/rules/*.md",
				Syntax:      "YAML Frontmatter: globs: [\"**/*.go\"], alwaysApply: true",
				TriggerMode: "Glob Match / Always-On",
				Description: "Coding constraints applied automatically when editing matching files",
			},
			{
				EntityType:  ir.TypeWorkflow,
				Location:    ".agents/workflows/*.md",
				Syntax:      "YAML Frontmatter: description, slash_command",
				TriggerMode: "Slash Command (/<name>)",
				Description: "Chained multi-step task workflows with planning and verification loops",
			},
			{
				EntityType:  ir.TypeSkill,
				Location:    ".agents/skills/<name>/SKILL.md",
				Syntax:      "YAML Frontmatter: name, description + Markdown body",
				TriggerMode: "JIT On-Demand (Model Decision / Tool call)",
				Description: "Domain-specific capability packages loaded Just-In-Time to save context tokens",
			},
			{
				EntityType:  ir.TypePrompt,
				Location:    ".agents/workflows/*.md",
				Syntax:      "Parametric prompt template inside workflows",
				TriggerMode: "On-Demand / Slash Command",
				Description: "Reusable prompt patterns and task scaffolding",
			},
			{
				EntityType:  ir.TypeAgent,
				Location:    ".agents/skills/<name>/SKILL.md / Subagents",
				Syntax:      "YAML Frontmatter + Tool/MCP Bindings",
				TriggerMode: "Subagent Delegation / JIT Skill",
				Description: "Specialized persona with isolated tools, instructions, and hooks",
			},
		},
	},
	"copilot": {
		ID:          "copilot",
		Name:        "GitHub Copilot",
		RootDir:     ".github/",
		SSOTFile:    ".github/copilot-instructions.md",
		Description: "GitHub Copilot repository instructions, prompts, and agent definitions",
		Entities: []EntitySpec{
			{
				EntityType:  ir.TypeInstruction,
				Location:    ".github/copilot-instructions.md",
				Syntax:      "Plain Markdown",
				TriggerMode: "Always-On",
				Description: "Repository-wide instructions included in every Copilot chat request",
			},
			{
				EntityType:  ir.TypeRule,
				Location:    ".github/instructions/*.instructions.md",
				Syntax:      "YAML Frontmatter: applyTo: \"**/*.ts\", description",
				TriggerMode: "File Pattern (applyTo)",
				Description: "Language and directory-specific coding guidelines",
			},
			{
				EntityType:  ir.TypeWorkflow,
				Location:    ".github/prompts/*.prompt.md",
				Syntax:      "YAML Frontmatter: name, description, argument-hint, agent",
				TriggerMode: "Prompt Slash Command (/<name>)",
				Description: "Reusable task prompts and code generation templates",
			},
			{
				EntityType:  ir.TypeSkill,
				Location:    ".github/skills/<name>/SKILL.md",
				Syntax:      "Directory with SKILL.md",
				TriggerMode: "On-Demand",
				Description: "Specialized skill packages for GitHub Copilot Workspace / Agent",
			},
			{
				EntityType:  ir.TypePrompt,
				Location:    ".github/prompts/*.prompt.md",
				Syntax:      "YAML Frontmatter + Prompt Template",
				TriggerMode: "Prompt Command",
				Description: "Task prompts with optional argument placeholders",
			},
			{
				EntityType:  ir.TypeAgent,
				Location:    ".github/agents/*.agent.md",
				Syntax:      "YAML Frontmatter: description, tools, model",
				TriggerMode: "Custom Agent Mode",
				Description: "Custom agent personas with specialized instructions",
			},
		},
	},
	"cursor": {
		ID:          "cursor",
		Name:        "Cursor AI",
		RootDir:     ".cursor/rules/",
		SSOTFile:    ".cursor/rules/base.mdc",
		Description: "Cursor AI MDC rule files with glob filters and intelligent activation",
		Entities: []EntitySpec{
			{
				EntityType:  ir.TypeInstruction,
				Location:    ".cursor/rules/base.mdc (or .cursorrules)",
				Syntax:      "YAML Frontmatter: alwaysApply: true, globs: [\"*\"]",
				TriggerMode: "Always-On",
				Description: "Global workspace instructions applied to all Cursor requests",
			},
			{
				EntityType:  ir.TypeRule,
				Location:    ".cursor/rules/*.mdc",
				Syntax:      "YAML Frontmatter: description, globs: [\"*.go\"], alwaysApply: false",
				TriggerMode: "Glob Filter / Apply Intelligently",
				Description: "Contextual rules attached when relevant files are open or edited",
			},
			{
				EntityType:  ir.TypeWorkflow,
				Location:    ".cursor/rules/*.mdc",
				Syntax:      "Rule with description for agent step-by-step guidance",
				TriggerMode: "@Rule Mention / Apply Intelligently",
				Description: "Procedural instructions invoked via rule mentions or AI decision",
			},
			{
				EntityType:  ir.TypeSkill,
				Location:    ".cursor/rules/<name>.mdc",
				Syntax:      "YAML Frontmatter: description (no alwaysApply)",
				TriggerMode: "Apply Intelligently (Model Decision)",
				Description: "Specialized knowledge rules loaded dynamically when matching keywords",
			},
			{
				EntityType:  ir.TypePrompt,
				Location:    ".cursor/rules/*.mdc",
				Syntax:      "Rule template for task scaffolding",
				TriggerMode: "@Rule Mention",
				Description: "Reusable prompt contexts attached via chat mention",
			},
			{
				EntityType:  ir.TypeAgent,
				Location:    ".cursor/rules/*.mdc + MCP Config",
				Syntax:      "MDC Rules + .cursor/mcp.json",
				TriggerMode: "Agent Mode + MCP",
				Description: "Agent personas configured via custom instructions and MCP tool access",
			},
		},
	},
	"claude": {
		ID:          "claude",
		Name:        "Claude Code",
		RootDir:     ".claude/",
		SSOTFile:    "CLAUDE.md",
		Description: "Anthropic Claude Code CLI workspace instructions and skills",
		Entities: []EntitySpec{
			{
				EntityType:  ir.TypeInstruction,
				Location:    "CLAUDE.md / ~/.claude/CLAUDE.md",
				Syntax:      "Concise Markdown (< 200 lines)",
				TriggerMode: "Always-On",
				Description: "Core project commands, architecture, and conventions for Claude CLI",
			},
			{
				EntityType:  ir.TypeRule,
				Location:    "CLAUDE.md (Directory Cascading)",
				Syntax:      "Subdirectory CLAUDE.md files",
				TriggerMode: "Directory Scope",
				Description: "Hierarchical guidelines loaded when Claude navigates subdirectories",
			},
			{
				EntityType:  ir.TypeWorkflow,
				Location:    ".claude/workflows/*.md / Slash Commands",
				Syntax:      "Numbered steps and verification checks",
				TriggerMode: "Interactive CLI Execution",
				Description: "Step-by-step development workflows executed via Claude commands",
			},
			{
				EntityType:  ir.TypeSkill,
				Location:    ".claude/skills/<name>/SKILL.md",
				Syntax:      "YAML Frontmatter + Skill Instructions",
				TriggerMode: "On-Demand Skill Invocation",
				Description: "Modular capabilities loaded dynamically during task execution",
			},
			{
				EntityType:  ir.TypePrompt,
				Location:    ".claude/prompts/*.md",
				Syntax:      "Task template markdown",
				TriggerMode: "Prompt file inclusion",
				Description: "Reusable prompt instructions",
			},
			{
				EntityType:  ir.TypeAgent,
				Location:    ".claude/subagents/",
				Syntax:      "Agent permission and role definitions",
				TriggerMode: "Subagent spawning",
				Description: "Isolated subagents for specialized long-running tasks",
			},
		},
	},
}

// GetPlatformSpec returns the spec for a platform or a default normalized name
func GetPlatformSpec(platform string) (PlatformSpec, bool) {
	key := strings.ToLower(strings.TrimSpace(platform))
	switch key {
	case "antigravity", "gemini-agent", "google":
		key = "antigravity"
	case "copilot", "github", "gh-copilot":
		key = "copilot"
	case "cursor", "cursor-ai":
		key = "cursor"
	case "claude", "claude-code", "anthropic":
		key = "claude"
	}

	spec, ok := platformSpecs[key]
	return spec, ok
}

// BuildSpecMatrix creates a conceptual comparison matrix between two platforms
func BuildSpecMatrix(fromPlatform, toPlatform string) *SpecMatrixReport {
	fromSpec, fromOk := GetPlatformSpec(fromPlatform)
	if !fromOk {
		fromSpec = PlatformSpec{ID: fromPlatform, Name: strings.ToUpper(fromPlatform)}
	}
	toSpec, toOk := GetPlatformSpec(toPlatform)
	if !toOk {
		toSpec = PlatformSpec{ID: toPlatform, Name: strings.ToUpper(toPlatform)}
	}

	entityTypes := []ir.EntityType{
		ir.TypeInstruction,
		ir.TypeRule,
		ir.TypeWorkflow,
		ir.TypeSkill,
		ir.TypePrompt,
		ir.TypeAgent,
	}

	var items []SpecMatrixItem

	for _, entityType := range entityTypes {
		var fromEntity, toEntity EntitySpec
		for _, e := range fromSpec.Entities {
			if e.EntityType == entityType {
				fromEntity = e
				break
			}
		}
		for _, e := range toSpec.Entities {
			if e.EntityType == entityType {
				toEntity = e
				break
			}
		}

		item := SpecMatrixItem{
			EntityType:     entityType,
			SourceLocation: fromEntity.Location,
			SourceSyntax:   fromEntity.Syntax,
			TargetLocation: toEntity.Location,
			TargetSyntax:   toEntity.Syntax,
			TargetBehavior: toEntity.Description,
		}
		if item.SourceLocation == "" {
			item.SourceLocation = "(Not standard)"
		}
		if item.TargetLocation == "" {
			item.TargetLocation = "(Mapped to Rule/Instruction)"
		}
		items = append(items, item)
	}

	return &SpecMatrixReport{
		FromPlatform: fromSpec.Name,
		ToPlatform:   toSpec.Name,
		Items:        items,
	}
}
