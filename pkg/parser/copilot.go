package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/agent-parkour/pkg/ir"
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
	ApplyTo      interface{} `yaml:"applyTo"`
	ExcludeAgent []string    `yaml:"excludeAgent"`
	Description  string      `yaml:"description"`
}

func parseApplyTo(raw interface{}) []string {
	if raw == nil {
		return nil
	}
	switch v := raw.(type) {
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			return []string{trimmed}
		}
	case []interface{}:
		var result []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				s = strings.TrimSpace(s)
				if s != "" {
					result = append(result, s)
				}
			}
		}
		return result
	case []string:
		var result []string
		for _, s := range v {
			s = strings.TrimSpace(s)
			if s != "" {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

func inferGlobFromID(id string) []string {
	knownExts := map[string]string{
		"go":       "**/*.go",
		"ts":       "**/*.ts",
		"js":       "**/*.js",
		"py":       "**/*.py",
		"java":     "**/*.java",
		"rust":     "**/*.rs",
		"markdown": "**/*.md",
		"md":       "**/*.md",
		"html":     "**/*.html",
		"css":      "**/*.css",
	}
	lower := strings.ToLower(id)
	if glob, ok := knownExts[lower]; ok {
		return []string{glob}
	}
	return []string{id + "/**"}
}

// ParseCopilotDirectory parses a .github/ directory containing copilot-instructions.md, instructions/, prompts/, agents/, or a single file
func ParseCopilotDirectory(sourcePath string) ([]*ir.UADocument, error) {
	fi, err := os.Stat(sourcePath)
	if err == nil && !fi.IsDir() {
		data, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, err
		}
		baseName := filepath.Base(sourcePath)
		lower := strings.ToLower(sourcePath)

		if strings.HasSuffix(lower, "copilot-instructions.md") {
			doc := ir.NewDocument("instruction-global", ir.TypeInstruction, "Global Project Instructions")
			doc.Metadata.Description = "Global instructions and architecture rules imported from GitHub Copilot"
			doc.Activation.Mode = ir.ModeAlwaysOn
			doc.Payload.MarkdownBody = string(data)
			doc.Payload.RawSource = sourcePath
			return []*ir.UADocument{doc}, nil
		}

		if strings.Contains(lower, "prompts") || strings.HasSuffix(lower, ".prompt.md") {
			id := strings.TrimSuffix(baseName, ".prompt.md")
			id = strings.TrimSuffix(id, ".md")
			var fm CopilotPromptFrontmatter
			body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)
			name := fm.Name
			if name == "" {
				name = formatTitle(id)
			}
			doc := ir.NewDocument("prompt-"+id, ir.TypePrompt, name)
			doc.Metadata.Description = fm.Description
			doc.Activation.Mode = ir.ModeOnDemand
			doc.Activation.SlashCommand = id
			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = sourcePath
			return []*ir.UADocument{doc}, nil
		}

		// Otherwise parse as instruction/rule
		id := strings.TrimSuffix(baseName, ".instructions.md")
		id = strings.TrimSuffix(id, ".md")
		var fm CopilotInstructionFrontmatter
		body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)
		doc := ir.NewDocument("rule-"+id, ir.TypeRule, formatTitle(id))
		doc.Metadata.Description = fm.Description
		globs := parseApplyTo(fm.ApplyTo)
		if len(globs) == 0 {
			globs = inferGlobFromID(id)
		}
		doc.Activation.Mode = ir.ModeGlob
		doc.Activation.Globs = globs
		doc.Payload.MarkdownBody = body
		doc.Payload.RawSource = sourcePath
		return []*ir.UADocument{doc}, nil
	}

	var docs []*ir.UADocument
	githubDir := sourcePath

	// If sourcePath contains .github subdirectory, descend into it
	if child := filepath.Join(sourcePath, ".github"); func() bool {
		childFi, childErr := os.Stat(child)
		return childErr == nil && childFi.IsDir()
	}() {
		githubDir = child
	}

	findSubDir := func(name string) string {
		p1 := filepath.Join(githubDir, name)
		if sfi, serr := os.Stat(p1); serr == nil && sfi.IsDir() {
			return p1
		}
		p2 := filepath.Join(sourcePath, name)
		if sfi, serr := os.Stat(p2); serr == nil && sfi.IsDir() {
			return p2
		}
		return ""
	}

	// 1. Check copilot-instructions.md (Project Root Instruction)
	copilotInstCandidates := []string{
		filepath.Join(githubDir, "copilot-instructions.md"),
		filepath.Join(sourcePath, "copilot-instructions.md"),
	}
	for _, instPath := range copilotInstCandidates {
		if data, readErr := os.ReadFile(instPath); readErr == nil {
			doc := ir.NewDocument("instruction-global", ir.TypeInstruction, "Global Project Instructions")
			doc.Metadata.Description = "Global instructions and architecture rules imported from GitHub Copilot"
			doc.Activation.Mode = ir.ModeAlwaysOn
			doc.Payload.MarkdownBody = string(data)
			doc.Payload.RawSource = instPath
			docs = append(docs, doc)
			break
		}
	}

	// 2. Parse instructions/ (*.instructions.md or *.md)
	if instDir := findSubDir("instructions"); instDir != "" {
		if entries, readDirErr := os.ReadDir(instDir); readDirErr == nil {
			for _, entry := range entries {
				if entry.IsDir() || (!strings.HasSuffix(entry.Name(), ".md") && !strings.HasSuffix(entry.Name(), ".instructions.md")) {
					continue
				}
				filePath := filepath.Join(instDir, entry.Name())
				data, readErr := os.ReadFile(filePath)
				if readErr != nil {
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

				globs := parseApplyTo(fm.ApplyTo)
				if len(globs) == 0 {
					globs = inferGlobFromID(id)
				}
				doc.Activation.Mode = ir.ModeGlob
				doc.Activation.Globs = globs

				doc.Payload.MarkdownBody = body
				doc.Payload.RawSource = filePath
				docs = append(docs, doc)
			}
		}
	}

	// 3. Parse prompts/ (*.prompt.md or *.md)
	if promptsDir := findSubDir("prompts"); promptsDir != "" {
		if entries, readDirErr := os.ReadDir(promptsDir); readDirErr == nil {
			for _, entry := range entries {
				if entry.IsDir() || (!strings.HasSuffix(entry.Name(), ".md") && !strings.HasSuffix(entry.Name(), ".prompt.md")) {
					continue
				}
				filePath := filepath.Join(promptsDir, entry.Name())
				data, readErr := os.ReadFile(filePath)
				if readErr != nil {
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

				doc := ir.NewDocument("prompt-"+id, ir.TypePrompt, name)
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

	// 5. Parse skills/ (<name>/SKILL.md)
	if skillsDir := findSubDir("skills"); skillsDir != "" {
		if entries, readDirErr := os.ReadDir(skillsDir); readDirErr == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				skillFile := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
				data, readErr := os.ReadFile(skillFile)
				if readErr != nil {
					continue
				}

				var fm struct {
					Name        string `yaml:"name"`
					Description string `yaml:"description"`
				}
				body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)
				name := fm.Name
				if name == "" {
					name = entry.Name()
				}

				doc := ir.NewDocument("skill-"+entry.Name(), ir.TypeSkill, name)
				doc.Metadata.Description = fm.Description
				doc.Activation.Mode = ir.ModeOnDemand
				doc.Payload.MarkdownBody = body
				doc.Payload.RawSource = skillFile
				docs = append(docs, doc)
			}
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
