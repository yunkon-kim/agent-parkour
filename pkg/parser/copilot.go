package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/ir"
	"gopkg.in/yaml.v3"
)

// CopilotPromptFrontmatter represents YAML frontmatter in .github/prompts/*.prompt.md
type CopilotPromptFrontmatter struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	ArgumentHint string `yaml:"argument-hint"`
	Agent        string `yaml:"agent"`
	Model        string `yaml:"model"`
}

// CopilotInstructionFrontmatter represents YAML frontmatter in .github/instructions/*.instructions.md
type CopilotInstructionFrontmatter struct {
	ApplyTo      string   `yaml:"applyTo"`
	ExcludeAgent []string `yaml:"excludeAgent"`
	Description  string   `yaml:"description"`
}

// ParseCopilotDirectory parses a .github/ directory containing copilot-instructions.md, instructions/, prompts/, agents/
func ParseCopilotDirectory(githubDir string) ([]*ir.UADocument, error) {
	var docs []*ir.UADocument

	// 1. Check copilot-instructions.md (Project Root Instruction)
	copilotInstPath := filepath.Join(githubDir, "copilot-instructions.md")
	if data, err := os.ReadFile(copilotInstPath); err == nil {
		doc := ir.NewDocument("instruction-global", ir.TypeInstruction, "Global Project Instructions")
		doc.Metadata.Description = "Global instructions and architecture rules imported from GitHub Copilot"
		doc.Activation.Mode = ir.ModeAlwaysOn
		doc.Payload.MarkdownBody = string(data)
		doc.Payload.RawSource = copilotInstPath
		docs = append(docs, doc)
	}

	// 2. Parse instructions/ (*.instructions.md or *.md)
	instDir := filepath.Join(githubDir, "instructions")
	if entries, err := os.ReadDir(instDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || (!strings.HasSuffix(entry.Name(), ".md") && !strings.HasSuffix(entry.Name(), ".instructions.md")) {
				continue
			}
			filePath := filepath.Join(instDir, entry.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			id := strings.TrimSuffix(entry.Name(), ".instructions.md")
			id = strings.TrimSuffix(id, ".md")

			var fm CopilotInstructionFrontmatter
			body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)

			doc := ir.NewDocument("rule-"+id, ir.TypeRule, formatTitle(id))
			doc.Metadata.Description = fm.Description
			if doc.Metadata.Description == "" {
				doc.Metadata.Description = fmt.Sprintf("Coding rules and conventions for %s", id)
			}

			if fm.ApplyTo != "" {
				doc.Activation.Mode = ir.ModeGlob
				doc.Activation.Globs = []string{fm.ApplyTo}
			} else {
				doc.Activation.Mode = ir.ModeAlwaysOn
			}

			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = filePath
			docs = append(docs, doc)
		}
	}

	// 3. Parse prompts/ (*.prompt.md or *.md)
	promptsDir := filepath.Join(githubDir, "prompts")
	if entries, err := os.ReadDir(promptsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || (!strings.HasSuffix(entry.Name(), ".md") && !strings.HasSuffix(entry.Name(), ".prompt.md")) {
				continue
			}
			filePath := filepath.Join(promptsDir, entry.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			id := strings.TrimSuffix(entry.Name(), ".prompt.md")
			id = strings.TrimSuffix(id, ".md")

			var fm CopilotPromptFrontmatter
			body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)

			name := fm.Name
			if name == "" {
				name = formatTitle(id)
			}

			doc := ir.NewDocument("workflow-"+id, ir.TypeWorkflow, name)
			doc.Metadata.Description = fm.Description
			doc.Activation.Mode = ir.ModeOnDemand
			doc.Activation.SlashCommand = id

			if fm.ArgumentHint != "" {
				doc.Metadata.Tags = append(doc.Metadata.Tags, "hint:"+fm.ArgumentHint)
			}

			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = filePath
			docs = append(docs, doc)
		}
	}

	// 4. Parse agents/ (*.agent.md or *.md)
	agentsDir := filepath.Join(githubDir, "agents")
	if entries, err := os.ReadDir(agentsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			filePath := filepath.Join(agentsDir, entry.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			id := strings.TrimSuffix(entry.Name(), ".agent.md")
			id = strings.TrimSuffix(id, ".md")

			fmData, body, _ := SplitFrontmatter(string(data))
			doc := ir.NewDocument("agent-"+id, ir.TypeAgent, formatTitle(id))
			if desc, ok := fmData["description"].(string); ok {
				doc.Metadata.Description = desc
			}
			doc.Activation.Mode = ir.ModeOnDemand
			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = filePath
			docs = append(docs, doc)
		}
	}

	return docs, nil
}

func formatTitle(id string) string {
	parts := strings.Split(id, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// ParseUADocument parses a standard UA-IR YAML or Markdown with UA Frontmatter
func ParseUADocument(content string) (*ir.UADocument, error) {
	trimmed := strings.TrimSpace(content)
	if strings.HasPrefix(trimmed, "ua_version:") {
		var doc ir.UADocument
		err := yaml.Unmarshal([]byte(trimmed), &doc)
		if err != nil {
			return nil, err
		}
		return &doc, nil
	}

	var doc ir.UADocument
	body, err := ExtractFrontmatterAndUnmarshal(content, &doc)
	if err != nil {
		return nil, err
	}
	doc.Payload.MarkdownBody = body
	return &doc, nil
}
