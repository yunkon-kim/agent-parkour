package describer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/audit"
	"github.com/yunkon-kim/token-hop/pkg/emitter"
	"github.com/yunkon-kim/token-hop/pkg/ir"
)

// MappingAnalyzer analyzes UA-IR documents and generates transformation plans
type MappingAnalyzer struct {
	MaxTokens int
}

// NewMappingAnalyzer creates a new MappingAnalyzer
func NewMappingAnalyzer(maxTokens int) *MappingAnalyzer {
	if maxTokens <= 0 {
		maxTokens = 400
	}
	return &MappingAnalyzer{MaxTokens: maxTokens}
}

// Analyze generates a MappingReport for transforming docs from fromPlatform to toPlatform
func (a *MappingAnalyzer) Analyze(docs []*ir.UADocument, fromPlatform, toPlatform, sourceDir, outputDir string) *MappingReport {
	fromKey := strings.ToLower(strings.TrimSpace(fromPlatform))
	toKey := strings.ToLower(strings.TrimSpace(toPlatform))

	if outputDir == "" {
		outputDir = "."
	}

	report := &MappingReport{
		FromPlatform:     fromKey,
		ToPlatform:       toKey,
		SourceDir:        sourceDir,
		OutputDir:        outputDir,
		TotalSourceFiles: len(docs),
		Notes: []string{
			"Estimated token counts (~Tokens) are calculated locally without external API calls or cost, indicating baseline prompt context overhead per conversation turn.",
		},
	}

	for _, doc := range docs {
		chars := len(doc.Payload.MarkdownBody)
		tokens := audit.EstimateTokens(doc.Payload.MarkdownBody)
		report.TotalCharacters += chars
		report.TotalTokens += tokens

		isOversized := tokens > a.MaxTokens
		if isOversized {
			report.OversizedCount++
		}

		sourcePath := doc.Payload.RawSource
		if sourcePath == "" {
			sourcePath = fmt.Sprintf("%s (%s)", doc.Metadata.ID, doc.Metadata.Type)
		}

		sourceTrigger := formatTrigger(fromKey, doc)
		targetPath, targetType, targetTrigger, action, isAlwaysOn := resolveTargetMapping(toKey, doc, outputDir)

		if isAlwaysOn {
			report.AlwaysOnTokens += tokens
		} else {
			report.OnDemandTokens += tokens
		}

		rec := ""
		if isOversized {
			if doc.Metadata.Type == ir.TypeRule {
				rec = fmt.Sprintf("Rule exceeds %d tokens. Convert to JIT Skill or narrow Glob pattern.", a.MaxTokens)
			} else {
				rec = fmt.Sprintf("Document exceeds %d tokens. Consider decomposing into smaller modules.", a.MaxTokens)
			}
		}

		report.Items = append(report.Items, PlanItem{
			SourcePath:     sourcePath,
			SourceType:     doc.Metadata.Type,
			SourceTrigger:  sourceTrigger,
			TargetPath:     targetPath,
			TargetType:     targetType,
			TargetTrigger:  targetTrigger,
			Characters:     chars,
			Tokens:         tokens,
			IsOversized:    isOversized,
			Action:         action,
			Recommendation: rec,
		})
	}

	report.TotalTargetFiles = len(report.Items)
	return report
}

func formatTrigger(platform string, doc *ir.UADocument) string {
	if doc.Activation.Mode == ir.ModeAlwaysOn {
		return "Always-On"
	}
	if len(doc.Activation.Globs) > 0 {
		switch platform {
		case "copilot":
			return fmt.Sprintf("applyTo: %s", strings.Join(doc.Activation.Globs, ", "))
		case "cursor":
			return fmt.Sprintf("globs: [%s]", strings.Join(doc.Activation.Globs, ", "))
		default:
			return fmt.Sprintf("globs: [%s]", strings.Join(doc.Activation.Globs, ", "))
		}
	}
	if doc.Activation.SlashCommand != "" {
		return fmt.Sprintf("/%s", doc.Activation.SlashCommand)
	}
	if doc.Activation.Mode == ir.ModeOnDemand {
		return "On-Demand (JIT)"
	}
	if doc.Activation.Mode == ir.ModeModelDecision {
		return "Model Decision"
	}
	return "Always-On"
}

func resolveTargetMapping(targetPlatform string, doc *ir.UADocument, outputDir string) (targetPath string, targetType ir.EntityType, targetTrigger string, action string, isAlwaysOn bool) {
	targetType = doc.Metadata.Type
	sanitizedID := emitter.SanitizeFileName(doc.Metadata.ID)

	switch targetPlatform {
	case "antigravity", "gemini":
		switch doc.Metadata.Type {
		case ir.TypeInstruction:
			targetPath = filepath.Join(outputDir, ".agents", "rules", sanitizedID+".md")
			targetTrigger = "Always-On"
			action = "Extract Master SSOT -> Global Rule"
			isAlwaysOn = true
		case ir.TypeRule:
			targetPath = filepath.Join(outputDir, ".agents", "rules", sanitizedID+".md")
			if len(doc.Activation.Globs) > 0 {
				targetTrigger = fmt.Sprintf("globs: [%s]", strings.Join(doc.Activation.Globs, ", "))
				action = "YAML Frontmatter Transpilation"
				isAlwaysOn = false
			} else {
				targetTrigger = "Always-On (alwaysApply: true)"
				action = "Global Rule with alwaysApply"
				isAlwaysOn = true
			}
		case ir.TypeWorkflow, ir.TypePrompt:
			cmdName := doc.Activation.SlashCommand
			if cmdName == "" {
				cmdName = sanitizedID
			}
			cmdName = strings.TrimPrefix(cmdName, "workflow-")
			cmdName = strings.TrimPrefix(cmdName, "prompt-")
			targetPath = filepath.Join(outputDir, ".agents", "workflows", cmdName+".md")
			targetTrigger = fmt.Sprintf("/%s (Slash Command)", cmdName)
			action = "Slash Command Workflow Binding"
			targetType = ir.TypeWorkflow
			isAlwaysOn = false
		case ir.TypeSkill, ir.TypeAgent:
			skillName := strings.TrimPrefix(sanitizedID, "skill-")
			skillName = strings.TrimPrefix(skillName, "agent-")
			targetPath = filepath.Join(outputDir, ".agents", "skills", skillName, "SKILL.md")
			targetTrigger = "JIT On-Demand (Model Decision)"
			action = "Convert to JIT Skill Package"
			targetType = ir.TypeSkill
			isAlwaysOn = false
		default:
			targetPath = filepath.Join(outputDir, ".agents", "rules", sanitizedID+".md")
			targetTrigger = "Always-On"
			action = "Transpile to Rule"
			isAlwaysOn = true
		}

	case "cursor":
		fileName := sanitizedID + ".mdc"
		targetPath = filepath.Join(outputDir, ".cursor", "rules", fileName)
		if doc.Activation.Mode == ir.ModeAlwaysOn || doc.Metadata.Type == ir.TypeInstruction {
			targetTrigger = "alwaysApply: true (globs: [*])"
			action = "Transpile to Always-On MDC Rule"
			isAlwaysOn = true
		} else if len(doc.Activation.Globs) > 0 {
			targetTrigger = fmt.Sprintf("globs: [%s]", strings.Join(doc.Activation.Globs, ", "))
			action = "Transpile to Contextual MDC Rule"
			isAlwaysOn = false
		} else {
			targetTrigger = "Apply Intelligently (Model Decision)"
			action = "Transpile to Intelligent MDC Rule"
			isAlwaysOn = false
		}

	case "copilot", "github":
		switch doc.Metadata.Type {
		case ir.TypeInstruction:
			targetPath = filepath.Join(outputDir, ".github", "copilot-instructions.md")
			targetTrigger = "Always-On"
			action = "Compile to Repository Instructions"
			isAlwaysOn = true
		case ir.TypeRule:
			id := strings.TrimPrefix(sanitizedID, "rule-")
			targetPath = filepath.Join(outputDir, ".github", "instructions", id+".instructions.md")
			if len(doc.Activation.Globs) > 0 {
				targetTrigger = fmt.Sprintf("applyTo: \"%s\"", doc.Activation.Globs[0])
			} else {
				targetTrigger = "applyTo: \"**/*\""
			}
			action = "Transpile to Copilot Instruction"
			isAlwaysOn = false
		case ir.TypeWorkflow, ir.TypePrompt:
			id := strings.TrimPrefix(sanitizedID, "workflow-")
			id = strings.TrimPrefix(id, "prompt-")
			targetPath = filepath.Join(outputDir, ".github", "prompts", id+".prompt.md")
			targetTrigger = fmt.Sprintf("/%s Prompt Command", id)
			action = "Transpile to Copilot Prompt Template"
			isAlwaysOn = false
		case ir.TypeSkill, ir.TypeAgent:
			id := strings.TrimPrefix(sanitizedID, "skill-")
			id = strings.TrimPrefix(id, "agent-")
			targetPath = filepath.Join(outputDir, ".github", "agents", id+".agent.md")
			targetTrigger = "Custom Agent Mode"
			action = "Compile to Copilot Agent Definition"
			isAlwaysOn = false
		default:
			targetPath = filepath.Join(outputDir, ".github", "instructions", sanitizedID+".instructions.md")
			targetTrigger = "Always-On"
			action = "Transpile to Copilot Rule"
			isAlwaysOn = true
		}

	case "claude":
		switch doc.Metadata.Type {
		case ir.TypeInstruction, ir.TypeRule:
			targetPath = filepath.Join(outputDir, "CLAUDE.md")
			targetTrigger = "Always-On"
			action = "Append to CLAUDE.md Guidelines"
			isAlwaysOn = true
		case ir.TypeSkill, ir.TypeAgent:
			skillName := strings.TrimPrefix(sanitizedID, "skill-")
			skillName = strings.TrimPrefix(skillName, "agent-")
			targetPath = filepath.Join(outputDir, ".claude", "skills", skillName, "SKILL.md")
			targetTrigger = "On-Demand Skill Invocation"
			action = "Convert to Claude Skill Package"
			isAlwaysOn = false
		case ir.TypeWorkflow, ir.TypePrompt:
			cmdName := strings.TrimPrefix(sanitizedID, "workflow-")
			cmdName = strings.TrimPrefix(cmdName, "prompt-")
			targetPath = filepath.Join(outputDir, ".claude", "workflows", cmdName+".md")
			targetTrigger = "CLI Workflow Step"
			action = "Convert to Claude Workflow Step"
			isAlwaysOn = false
		default:
			targetPath = filepath.Join(outputDir, "CLAUDE.md")
			targetTrigger = "Always-On"
			action = "Append to CLAUDE.md"
			isAlwaysOn = true
		}

	default:
		targetPath = filepath.Join(outputDir, sanitizedID+".md")
		targetTrigger = "Always-On"
		action = fmt.Sprintf("Transpile to %s Format", targetPlatform)
		isAlwaysOn = true
	}

	return
}
