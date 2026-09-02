package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/ir"
)

// AntigravityRuleFrontmatter represents YAML frontmatter in .agents/rules/*.md
type AntigravityRuleFrontmatter struct {
	Globs       []string `yaml:"globs,omitempty"`
	Description string   `yaml:"description,omitempty"`
	AlwaysApply bool     `yaml:"alwaysApply,omitempty"`
}

// AntigravitySkillFrontmatter represents YAML frontmatter in .agent/skills/<name>/SKILL.md
type AntigravitySkillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ParseAntigravityDirectory parses Antigravity directory structure or a single file
func ParseAntigravityDirectory(baseDir string) ([]*ir.UADocument, error) {
	fi, err := os.Stat(baseDir)
	if err == nil && !fi.IsDir() {
		data, err := os.ReadFile(baseDir)
		if err != nil {
			return nil, err
		}
		baseName := filepath.Base(baseDir)
		lower := strings.ToLower(baseDir)

		if strings.HasSuffix(lower, "agents.md") || strings.HasSuffix(lower, "gemini.md") {
			doc := ir.NewDocument("instruction-agents-ssot", ir.TypeInstruction, "Master Project SSOT")
			doc.Metadata.Description = "Single Source of Truth instructions"
			doc.Activation.Mode = ir.ModeAlwaysOn
			doc.Payload.MarkdownBody = string(data)
			doc.Payload.RawSource = baseDir
			return []*ir.UADocument{doc}, nil
		}

		if strings.Contains(lower, "workflows") {
			id := strings.TrimSuffix(baseName, ".md")
			var fm struct {
				Description string `yaml:"description"`
			}
			body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)
			doc := ir.NewDocument("workflow-"+id, ir.TypeWorkflow, formatTitle(id))
			doc.Metadata.Description = fm.Description
			doc.Activation.Mode = ir.ModeOnDemand
			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = baseDir
			return []*ir.UADocument{doc}, nil
		}

		if strings.Contains(lower, "skills") || strings.HasSuffix(lower, "skill.md") {
			var fm AntigravitySkillFrontmatter
			body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)
			skillName := fm.Name
			if skillName == "" {
				skillName = filepath.Base(filepath.Dir(baseDir))
			}
			doc := ir.NewDocument("skill-"+skillName, ir.TypeSkill, formatTitle(skillName))
			doc.Metadata.Description = fm.Description
			doc.Activation.Mode = ir.ModeOnDemand
			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = baseDir
			return []*ir.UADocument{doc}, nil
		}

		// Default parse as rule
		id := strings.TrimSuffix(baseName, ".md")
		var fm AntigravityRuleFrontmatter
		body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)
		doc := ir.NewDocument("rule-"+id, ir.TypeRule, formatTitle(id))
		doc.Metadata.Description = fm.Description
		if len(fm.Globs) > 0 {
			doc.Activation.Mode = ir.ModeGlob
			doc.Activation.Globs = fm.Globs
		} else if fm.AlwaysApply {
			doc.Activation.Mode = ir.ModeAlwaysOn
		} else {
			doc.Activation.Mode = ir.ModeModelDecision
		}
		doc.Payload.MarkdownBody = body
		doc.Payload.RawSource = baseDir
		return []*ir.UADocument{doc}, nil
	}

	var docs []*ir.UADocument

	// 1. Check AGENTS.md in baseDir
	agentsMdPath := filepath.Join(baseDir, "AGENTS.md")
	if data, err := os.ReadFile(agentsMdPath); err == nil {
		doc := ir.NewDocument("instruction-agents-ssot", ir.TypeInstruction, "Master Project SSOT")
		doc.Metadata.Description = "Single Source of Truth instructions"
		doc.Activation.Mode = ir.ModeAlwaysOn
		doc.Payload.MarkdownBody = string(data)
		doc.Payload.RawSource = agentsMdPath
		docs = append(docs, doc)
	}

	// 2. Parse .agents/rules/
	rulesDir := filepath.Join(baseDir, ".agents", "rules")
	if entries, err := os.ReadDir(rulesDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			filePath := filepath.Join(rulesDir, entry.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			id := strings.TrimSuffix(entry.Name(), ".md")
			var fm AntigravityRuleFrontmatter
			body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)

			doc := ir.NewDocument("rule-"+id, ir.TypeRule, formatTitle(id))
			doc.Metadata.Description = fm.Description
			if len(fm.Globs) > 0 {
				doc.Activation.Mode = ir.ModeGlob
				doc.Activation.Globs = fm.Globs
			} else {
				doc.Activation.Mode = ir.ModeAlwaysOn
			}
			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = filePath
			docs = append(docs, doc)
		}
	}

	// 3. Parse .agents/workflows/
	workflowsDir := filepath.Join(baseDir, ".agents", "workflows")
	if entries, err := os.ReadDir(workflowsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			filePath := filepath.Join(workflowsDir, entry.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			id := strings.TrimSuffix(entry.Name(), ".md")
			fmData, body, _ := SplitFrontmatter(string(data))

			doc := ir.NewDocument("workflow-"+id, ir.TypeWorkflow, formatTitle(id))
			if desc, ok := fmData["description"].(string); ok {
				doc.Metadata.Description = desc
			}
			doc.Activation.Mode = ir.ModeOnDemand
			doc.Activation.SlashCommand = id
			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = filePath
			docs = append(docs, doc)
		}
	}

	// 4. Parse .agent/skills/ or .agents/skills/
	skillsDirs := []string{
		filepath.Join(baseDir, ".agent", "skills"),
		filepath.Join(baseDir, ".agents", "skills"),
	}
	for _, sDir := range skillsDirs {
		if entries, err := os.ReadDir(sDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				skillFile := filepath.Join(sDir, entry.Name(), "SKILL.md")
				data, err := os.ReadFile(skillFile)
				if err != nil {
					continue
				}

				var fm AntigravitySkillFrontmatter
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

// ParseCursorDirectory parses .cursor/rules/*.mdc or root .cursorrules
func ParseCursorDirectory(cursorPath string) ([]*ir.UADocument, error) {
	var docs []*ir.UADocument

	// 1. Check if single file (e.g. .cursorrules or specific .mdc)
	fi, err := os.Stat(cursorPath)
	if err == nil && !fi.IsDir() {
		data, err := os.ReadFile(cursorPath)
		if err != nil {
			return nil, err
		}
		doc := ir.NewDocument("instruction-cursor-root", ir.TypeInstruction, "Cursor Root Rules")
		doc.Metadata.Description = "Root instructions imported from Cursor"
		doc.Activation.Mode = ir.ModeAlwaysOn
		doc.Payload.MarkdownBody = string(data)
		doc.Payload.RawSource = cursorPath
		return []*ir.UADocument{doc}, nil
	}

	// 2. Check root .cursorrules if cursorPath is directory
	legacyCursorRules := filepath.Join(cursorPath, ".cursorrules")
	if data, err := os.ReadFile(legacyCursorRules); err == nil {
		doc := ir.NewDocument("instruction-cursorrules", ir.TypeInstruction, "Cursor Legacy Rules")
		doc.Metadata.Description = "Root instructions imported from .cursorrules"
		doc.Activation.Mode = ir.ModeAlwaysOn
		doc.Payload.MarkdownBody = string(data)
		doc.Payload.RawSource = legacyCursorRules
		docs = append(docs, doc)
	}

	// 3. Scan directory (.cursor/rules/ or cursorPath)
	targetDir := cursorPath
	if _, err := os.Stat(filepath.Join(cursorPath, ".cursor", "rules")); err == nil {
		targetDir = filepath.Join(cursorPath, ".cursor", "rules")
	}

	entries, err := os.ReadDir(targetDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || (!strings.HasSuffix(entry.Name(), ".mdc") && !strings.HasSuffix(entry.Name(), ".md")) {
				continue
			}
			filePath := filepath.Join(targetDir, entry.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			id := strings.TrimSuffix(entry.Name(), ".mdc")
			id = strings.TrimSuffix(id, ".md")

			var fm struct {
				Description string   `yaml:"description"`
				Globs       []string `yaml:"globs"`
				AlwaysApply bool     `yaml:"alwaysApply"`
			}
			body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)

			doc := ir.NewDocument("rule-"+id, ir.TypeRule, formatTitle(id))
			doc.Metadata.Description = fm.Description
			if fm.AlwaysApply {
				doc.Activation.Mode = ir.ModeAlwaysOn
			} else if len(fm.Globs) > 0 {
				doc.Activation.Mode = ir.ModeGlob
				doc.Activation.Globs = fm.Globs
			} else {
				doc.Activation.Mode = ir.ModeModelDecision
			}

			doc.Payload.MarkdownBody = body
			doc.Payload.RawSource = filePath
			docs = append(docs, doc)
		}
	}

	return docs, nil
}
