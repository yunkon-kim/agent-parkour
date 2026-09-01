package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/ai"
	"github.com/yunkon-kim/token-hop/pkg/audit"
	"github.com/yunkon-kim/token-hop/pkg/emitter"
	"github.com/yunkon-kim/token-hop/pkg/ir"
	"github.com/yunkon-kim/token-hop/pkg/parser"
	"gopkg.in/yaml.v3"
)

// Config represents token-hop.yaml configuration
type Config struct {
	Version string `yaml:"version"`
	SSOT    string `yaml:"ssot"`
	Targets []struct {
		Name            string `yaml:"name"`
		OutputDir       string `yaml:"output_dir,omitempty"`
		OutputFile      string `yaml:"output_file,omitempty"`
		SkillsDir       string `yaml:"skills_dir,omitempty"`
		InstructionsDir string `yaml:"instructions_dir,omitempty"`
		PromptsDir      string `yaml:"prompts_dir,omitempty"`
		EnableSkills    bool   `yaml:"enable_skills,omitempty"`
	} `yaml:"targets"`
	ContextBudget struct {
		MaxTokensPerRule     int  `yaml:"max_tokens_per_rule"`
		MaxCharactersPerRule int  `yaml:"max_characters_per_rule"`
		AutoDecompose        bool `yaml:"auto_decompose"`
	} `yaml:"context_budget"`
}

// Engine coordinates parsing, AI augmentation, and emitting
type Engine struct {
	Config     *Config
	AIProvider ai.AIProvider
}

// NewEngine creates a new Engine
func NewEngine(cfg *Config) *Engine {
	return &Engine{Config: cfg}
}

// SetAIProvider binds an active AI provider (Gemini, Claude, OpenAI, Ollama, Mock)
func (e *Engine) SetAIProvider(provider ai.AIProvider) {
	e.AIProvider = provider
}

// LoadConfig loads token-hop.yaml or returns defaults
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{
			Version: "1.0",
			SSOT:    "AGENTS.md",
		}, nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// AutoDetectSource inspects a directory or repo root and returns the detected format and path
func (e *Engine) AutoDetectSource(root string) (string, string) {
	// 1. Direct Copilot directory check (.github/ or current folder containing copilot files)
	if _, err := os.Stat(filepath.Join(root, "copilot-instructions.md")); err == nil {
		return "copilot", root
	}
	if _, err := os.Stat(filepath.Join(root, "instructions")); err == nil {
		return "copilot", root
	}
	if _, err := os.Stat(filepath.Join(root, "prompts")); err == nil {
		return "copilot", root
	}

	// 2. Child .github directory check
	copilotDir := filepath.Join(root, ".github")
	if _, err := os.Stat(filepath.Join(copilotDir, "copilot-instructions.md")); err == nil {
		return "copilot", copilotDir
	}
	if _, err := os.Stat(filepath.Join(copilotDir, "instructions")); err == nil {
		return "copilot", copilotDir
	}

	// 3. Check for Antigravity (AGENTS.md or .agents/)
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); err == nil {
		return "antigravity", root
	}
	if _, err := os.Stat(filepath.Join(root, ".agents")); err == nil {
		return "antigravity", root
	}

	// 4. Check for Cursor (.cursor/rules/)
	cursorRules := filepath.Join(root, ".cursor", "rules")
	if _, err := os.Stat(cursorRules); err == nil {
		return "cursor", cursorRules
	}
	if _, err := os.Stat(filepath.Join(root, "rules")); err == nil {
		return "cursor", filepath.Join(root, "rules")
	}

	// Fallback to copilot if .github exists, else antigravity
	if _, err := os.Stat(filepath.Join(root, ".github")); err == nil {
		return "copilot", filepath.Join(root, ".github")
	}

	return "antigravity", root
}

// AutoDetectTargets determines target formats based on source format
func (e *Engine) AutoDetectTargets(sourceFormat, targetArg string) []string {
	targetArg = strings.TrimSpace(targetArg)
	if targetArg == "all" {
		return []string{"antigravity", "cursor", "copilot"}
	}
	if targetArg != "" && targetArg != "auto" {
		var list []string
		for _, item := range strings.Split(targetArg, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				list = append(list, item)
			}
		}
		return list
	}

	// Default: convert to other platforms
	switch strings.ToLower(sourceFormat) {
	case "copilot", "github":
		return []string{"antigravity", "cursor"}
	case "antigravity", "gemini":
		return []string{"cursor", "copilot"}
	case "cursor":
		return []string{"antigravity", "copilot"}
	default:
		return []string{"antigravity", "cursor"}
	}
}

// ParseSource parses instructions from a given source format and path
func (e *Engine) ParseSource(from string, sourcePath string) ([]*ir.UADocument, error) {
	from = strings.ToLower(from)
	switch from {
	case "copilot", "github":
		return parser.ParseCopilotDirectory(sourcePath)
	case "antigravity", "gemini":
		return parser.ParseAntigravityDirectory(sourcePath)
	case "cursor":
		return parser.ParseCursorDirectory(sourcePath)
	default:
		detectedFrom, detectedPath := e.AutoDetectSource(sourcePath)
		return e.ParseSource(detectedFrom, detectedPath)
	}
}

// EmitTarget emits documents to target format
func (e *Engine) EmitTarget(target string, docs []*ir.UADocument, outputDir string) ([]string, error) {
	target = strings.ToLower(target)
	switch target {
	case "antigravity", "gemini":
		em := emitter.NewAntigravityEmitter(outputDir)
		return em.Emit(docs)
	case "cursor":
		em := emitter.NewCursorEmitter(outputDir)
		return em.Emit(docs)
	case "copilot", "github":
		em := emitter.NewCopilotEmitter(outputDir)
		return em.Emit(docs)
	default:
		return nil, fmt.Errorf("unsupported target format: %s", target)
	}
}

// Convert converts instructions directly from source to target
func (e *Engine) Convert(from string, to string, inputPath string, outputDir string) ([]string, *audit.Report, error) {
	return e.ConvertWithAI(context.Background(), from, to, inputPath, outputDir, false, 400)
}

// ConvertWithAI converts instructions with optional AI-powered semantic decomposition
func (e *Engine) ConvertWithAI(ctx context.Context, from, to, inputPath, outputDir string, decompose bool, maxTokens int) ([]string, *audit.Report, error) {
	docs, err := e.ParseSource(from, inputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse source (%s): %w", from, err)
	}

	if len(docs) == 0 {
		return nil, nil, fmt.Errorf("no valid instruction documents found in %s", inputPath)
	}

	finalDocs := docs
	if decompose && e.AIProvider != nil {
		var decomposedList []*ir.UADocument
		for _, doc := range docs {
			tokenCount := audit.EstimateTokens(doc.Payload.MarkdownBody)
			if tokenCount > maxTokens {
				res, err := e.AIProvider.DecomposeRule(ctx, doc.Metadata.Name, doc.Payload.MarkdownBody, maxTokens)
				if err == nil && len(res.SubRules) > 0 {
					for _, sub := range res.SubRules {
						entityType := ir.TypeRule
						if sub.IsSkill {
							entityType = ir.TypeSkill
						}
						mode := ir.ModeGlob
						if len(sub.Globs) == 0 {
							mode = ir.ModeAlwaysOn
						}

						subDoc := ir.NewDocument(sub.ID, entityType, sub.Title)
						subDoc.Metadata.Description = sub.Description
						subDoc.Activation.Mode = mode
						subDoc.Activation.Globs = sub.Globs
						subDoc.Payload.MarkdownBody = sub.Content

						decomposedList = append(decomposedList, subDoc)
					}
					continue
				}
			}
			decomposedList = append(decomposedList, doc)
		}
		finalDocs = decomposedList
	}

	written, err := e.EmitTarget(to, finalDocs, outputDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to emit target (%s): %w", to, err)
	}

	auditReport := audit.AuditDocuments(finalDocs, maxTokens)
	return written, auditReport, nil
}
