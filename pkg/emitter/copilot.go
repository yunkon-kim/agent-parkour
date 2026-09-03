package emitter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/agent-parkour/pkg/ir"
	"gopkg.in/yaml.v3"
)

// CopilotEmitter emits UA-IR documents to GitHub Copilot (.github/) format
type CopilotEmitter struct {
	BaseDir string
}

// NewCopilotEmitter creates a new Copilot emitter
func NewCopilotEmitter(baseDir string) *CopilotEmitter {
	return &CopilotEmitter{BaseDir: baseDir}
}

// Emit writes documents to the .github directory structure with automatic backup
func (e *CopilotEmitter) Emit(docs []*ir.UADocument) ([]string, error) {
	var writtenFiles []string

	githubDir := filepath.Join(e.BaseDir, ".github")
	instDir := filepath.Join(githubDir, "instructions")
	promptsDir := filepath.Join(githubDir, "prompts")

	for _, doc := range docs {
		switch doc.Metadata.Type {
		case ir.TypeInstruction:
			filePath := filepath.Join(githubDir, "copilot-instructions.md")
			content := fmt.Sprintf("# %s\n\n%s\n", doc.Metadata.Name, doc.Payload.MarkdownBody)
			backup, err := SafeWriteFile(filePath, []byte(content))
			if err != nil {
				return writtenFiles, err
			}
			if backup != "" {
				fmt.Printf("      📦 Backed up existing file -> %s\n", backup)
			}
			writtenFiles = append(writtenFiles, filePath)

		case ir.TypeRule:
			id := SanitizeFileName(doc.Metadata.ID)
			id = strings.TrimPrefix(id, "rule-")
			filePath := filepath.Join(instDir, id+".instructions.md")

			fm := map[string]interface{}{}
			if len(doc.Activation.Globs) > 0 {
				fm["applyTo"] = doc.Activation.Globs[0]
			} else {
				fm["applyTo"] = "**/*"
			}
			if doc.Metadata.Description != "" {
				fm["description"] = doc.Metadata.Description
			}

			fmBytes, _ := yaml.Marshal(fm)
			var out strings.Builder
			out.WriteString("---\n")
			out.WriteString(string(fmBytes))
			out.WriteString("---\n\n")
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
			id := SanitizeFileName(doc.Metadata.ID)
			id = strings.TrimPrefix(id, "workflow-")
			id = strings.TrimPrefix(id, "prompt-")
			filePath := filepath.Join(promptsDir, id+".prompt.md")

			fm := map[string]interface{}{
				"name":        id,
				"description": doc.Metadata.Description,
				"agent":       "agent",
			}
			fmBytes, _ := yaml.Marshal(fm)
			var out strings.Builder
			out.WriteString("---\n")
			out.WriteString(string(fmBytes))
			out.WriteString("---\n\n")
			out.WriteString(doc.Payload.MarkdownBody)
			out.WriteString("\n")

			backup, err := SafeWriteFile(filePath, []byte(out.String()))
			if err != nil {
				return writtenFiles, err
			}
			if backup != "" {
				fmt.Printf("      📦 Backed up existing file -> %s\n", backup)
			}
			writtenFiles = append(writtenFiles, filePath)
		}
	}

	return writtenFiles, nil
}
