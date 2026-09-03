package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/agent-parkour/pkg/ir"
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

	// Resolve agents directory if baseDir contains .agents subdirectory
	agentsDir := baseDir
	if child := filepath.Join(baseDir, ".agents"); func() bool {
		childFi, childErr := os.Stat(child)
		return childErr == nil && childFi.IsDir()
	}() {
		agentsDir = child
	}

	findSubDir := func(name string) string {
		p1 := filepath.Join(agentsDir, name)
		if sfi, serr := os.Stat(p1); serr == nil && sfi.IsDir() {
			return p1
		}
		p2 := filepath.Join(baseDir, name)
		if sfi, serr := os.Stat(p2); serr == nil && sfi.IsDir() {
			return p2
		}
		return ""
	}

	// 1. Check AGENTS.md / GEMINI.md in baseDir or agentsDir
	agentsMdCandidates := []string{
		filepath.Join(baseDir, "AGENTS.md"),
		filepath.Join(agentsDir, "AGENTS.md"),
		filepath.Join(baseDir, "GEMINI.md"),
		filepath.Join(agentsDir, "GEMINI.md"),
	}
	for _, candPath := range agentsMdCandidates {
		if data, readErr := os.ReadFile(candPath); readErr == nil {
			doc := ir.NewDocument("instruction-agents-ssot", ir.TypeInstruction, "Master Project SSOT")
			doc.Metadata.Description = "Single Source of Truth instructions"
			doc.Activation.Mode = ir.ModeAlwaysOn
			doc.Payload.MarkdownBody = string(data)
			doc.Payload.RawSource = candPath
			docs = append(docs, doc)
			break
		}
	}

	// 2. Parse rules/ (*.md)
	if rulesDir := findSubDir("rules"); rulesDir != "" {
		if entries, readDirErr := os.ReadDir(rulesDir); readDirErr == nil {
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
					continue
				}
				filePath := filepath.Join(rulesDir, entry.Name())
				data, readErr := os.ReadFile(filePath)
				if readErr != nil {
					continue
				}

				id := strings.TrimSuffix(entry.Name(), ".md")
				var fm AntigravityRuleFrontmatter
				body, _ := ExtractFrontmatterAndUnmarshal(string(data), &fm)

				doc := ir.NewDocument("rule-"+id, ir.TypeRule, formatTitle(id))
				doc.Metadata.Description = fm.Description
				if fm.AlwaysApply {
					doc.Activation.Mode = ir.ModeAlwaysOn
				} else if len(fm.Globs) > 0 {
					doc.Activation.Mode = ir.ModeGlob
					doc.Activation.Globs = fm.Globs
				} else {
					doc.Activation.Mode = ir.ModeGlob
					doc.Activation.Globs = inferGlobFromID(id)
				}
				doc.Payload.MarkdownBody = body
				doc.Payload.RawSource = filePath
				docs = append(docs, doc)
			}
		}
	}

	// 3. Parse workflows/ (*.md)
	if workflowsDir := findSubDir("workflows"); workflowsDir != "" {
		if entries, readDirErr := os.ReadDir(workflowsDir); readDirErr == nil {
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
					continue
				}
				filePath := filepath.Join(workflowsDir, entry.Name())
				data, readErr := os.ReadFile(filePath)
				if readErr != nil {
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
	}

	// 4. Parse skills/
	var skillsDirs []string
	if s := findSubDir("skills"); s != "" {
		skillsDirs = append(skillsDirs, s)
	}
	if agentSkills := filepath.Join(baseDir, ".agent", "skills"); func() bool {
		fi, err := os.Stat(agentSkills)
		return err == nil && fi.IsDir()
	}() {
		skillsDirs = append(skillsDirs, agentSkills)
	}

	seenSkills := map[string]bool{}
	for _, sDir := range skillsDirs {
		if entries, readDirErr := os.ReadDir(sDir); readDirErr == nil {
			for _, entry := range entries {
				if !entry.IsDir() || seenSkills[entry.Name()] {
					continue
				}
				skillFile := filepath.Join(sDir, entry.Name(), "SKILL.md")
				data, readErr := os.ReadFile(skillFile)
				if readErr != nil {
					continue
				}
				seenSkills[entry.Name()] = true

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
