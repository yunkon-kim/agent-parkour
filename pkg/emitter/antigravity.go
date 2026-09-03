package emitter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/agent-parkour/pkg/ir"
	"gopkg.in/yaml.v3"
)

// AntigravityEmitter emits UA-IR documents to Google Antigravity format
type AntigravityEmitter struct {
	BaseDir string
}

// NewAntigravityEmitter creates a new Antigravity emitter
func NewAntigravityEmitter(baseDir string) *AntigravityEmitter {
	return &AntigravityEmitter{BaseDir: baseDir}
}

// Emit writes documents to the Antigravity directory structure with automatic backup
func (e *AntigravityEmitter) Emit(docs []*ir.UADocument) ([]string, error) {
	var writtenFiles []string

	rulesDir := filepath.Join(e.BaseDir, ".agents", "rules")
	workflowsDir := filepath.Join(e.BaseDir, ".agents", "workflows")
	skillsDir := filepath.Join(e.BaseDir, ".agents", "skills")

	var rootInstructionBody string

	for _, doc := range docs {
		switch doc.Metadata.Type {
		case ir.TypeInstruction:
			if rootInstructionBody == "" {
				rootInstructionBody = doc.Payload.MarkdownBody
			}
			fileName := SanitizeFileName(doc.Metadata.ID) + ".md"
			filePath := filepath.Join(rulesDir, fileName)
			content := RenderMarkdownWithTitle(doc.Metadata.Name, doc.Payload.MarkdownBody)
			backup, err := SafeWriteFile(filePath, []byte(content))
			if err != nil {
				return writtenFiles, err
			}
			if backup != "" {
				fmt.Printf("      📦 Backed up existing file -> %s\n", backup)
			}
			writtenFiles = append(writtenFiles, filePath)

		case ir.TypeRule:
			fileName := SanitizeFileName(doc.Metadata.ID) + ".md"
			filePath := filepath.Join(rulesDir, fileName)

			fm := map[string]interface{}{}
			if len(doc.Activation.Globs) > 0 {
				fm["globs"] = doc.Activation.Globs
			}
			if doc.Metadata.Description != "" {
				fm["description"] = doc.Metadata.Description
			}
			if doc.Activation.Mode == ir.ModeAlwaysOn {
				fm["alwaysApply"] = true
			}

			var out strings.Builder
			if len(fm) > 0 {
				fmBytes, _ := yaml.Marshal(fm)
				out.WriteString("---\n")
				out.WriteString(string(fmBytes))
				out.WriteString("---\n\n")
			}
			out.WriteString(RenderMarkdownWithTitle(doc.Metadata.Name, doc.Payload.MarkdownBody))

			backup, err := SafeWriteFile(filePath, []byte(out.String()))
			if err != nil {
				return writtenFiles, err
			}
			if backup != "" {
				fmt.Printf("      📦 Backed up existing file -> %s\n", backup)
			}
			writtenFiles = append(writtenFiles, filePath)

		case ir.TypeWorkflow, ir.TypePrompt:
			cmdName := doc.Activation.SlashCommand
			if cmdName == "" {
				cmdName = SanitizeFileName(doc.Metadata.ID)
			}
			cmdName = strings.TrimPrefix(cmdName, "workflow-")
			cmdName = strings.TrimPrefix(cmdName, "prompt-")

			fileName := cmdName + ".md"
			filePath := filepath.Join(workflowsDir, fileName)

			fm := map[string]interface{}{}
			if doc.Metadata.Description != "" {
				fm["description"] = doc.Metadata.Description
			}

			var out strings.Builder
			if len(fm) > 0 {
				fmBytes, _ := yaml.Marshal(fm)
				out.WriteString("---\n")
				out.WriteString(string(fmBytes))
				out.WriteString("---\n\n")
			}
			out.WriteString(RenderMarkdownWithTitle("Workflow: "+doc.Metadata.Name, doc.Payload.MarkdownBody))

			backup, err := SafeWriteFile(filePath, []byte(out.String()))
			if err != nil {
				return writtenFiles, err
			}
			if backup != "" {
				fmt.Printf("      📦 Backed up existing file -> %s\n", backup)
			}
			writtenFiles = append(writtenFiles, filePath)

		case ir.TypeSkill, ir.TypeAgent:
			skillName := SanitizeFileName(doc.Metadata.ID)
			skillName = strings.TrimPrefix(skillName, "skill-")
			skillName = strings.TrimPrefix(skillName, "agent-")

			targetSkillDir := filepath.Join(skillsDir, skillName)
			skillFilePath := filepath.Join(targetSkillDir, "SKILL.md")

			fm := map[string]interface{}{
				"name":        skillName,
				"description": doc.Metadata.Description,
			}
			if fm["description"] == "" {
				fm["description"] = fmt.Sprintf("Skill package for %s", doc.Metadata.Name)
			}

			fmBytes, _ := yaml.Marshal(fm)
			var out strings.Builder
			out.WriteString("---\n")
			out.WriteString(string(fmBytes))
			out.WriteString("---\n\n")
			out.WriteString(RenderMarkdownWithTitle(doc.Metadata.Name, doc.Payload.MarkdownBody))

			backup, err := SafeWriteFile(skillFilePath, []byte(out.String()))
			if err != nil {
				return writtenFiles, err
			}
			if backup != "" {
				fmt.Printf("      📦 Backed up existing file -> %s\n", backup)
			}
			writtenFiles = append(writtenFiles, skillFilePath)
		}
	}

	// Write master AGENTS.md in base directory
	agentsMdPath := filepath.Join(e.BaseDir, "AGENTS.md")
	var agentsContent strings.Builder
	agentsContent.WriteString("# Project Agent Guidelines (Single Source of Truth)\n\n")
	if rootInstructionBody != "" {
		agentsContent.WriteString(rootInstructionBody)
		agentsContent.WriteString("\n\n---\n\n")
	}
	agentsContent.WriteString("## Managed Rules & Workflows\n\n")
	for _, doc := range docs {
		agentsContent.WriteString(fmt.Sprintf("- **%s** (`%s`): %s\n", doc.Metadata.Name, doc.Metadata.Type, doc.Metadata.Description))
	}
	agentsContent.WriteString("\n*Compiled and synchronized by [agent-parkour](https://github.com/yunkon-kim/agent-parkour)*\n")

	backup, err := SafeWriteFile(agentsMdPath, []byte(agentsContent.String()))
	if err != nil {
		return writtenFiles, err
	}
	if backup != "" {
		fmt.Printf("      📦 Backed up existing file -> %s\n", backup)
	}
	writtenFiles = append(writtenFiles, agentsMdPath)

	return writtenFiles, nil
}
