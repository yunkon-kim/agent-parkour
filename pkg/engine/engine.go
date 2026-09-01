package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/budget"
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
		Name           string `yaml:"name"`
		OutputDir      string `yaml:"output_dir,omitempty"`
		OutputFile     string `yaml:"output_file,omitempty"`
		SkillsDir      string `yaml:"skills_dir,omitempty"`
		InstructionsDir string `yaml:"instructions_dir,omitempty"`
		PromptsDir     string `yaml:"prompts_dir,omitempty"`
		EnableSkills   bool   `yaml:"enable_skills,omitempty"`
	} `yaml:"targets"`
	ContextBudget struct {
		MaxTokensPerRule     int  `yaml:"max_tokens_per_rule"`
		MaxCharactersPerRule int  `yaml:"max_characters_per_rule"`
		AutoDecompose        bool `yaml:"auto_decompose"`
	} `yaml:"context_budget"`
}

// Engine coordinates parsing and emitting
type Engine struct {
	Config *Config
}

// NewEngine creates a new Engine
func NewEngine(cfg *Config) *Engine {
	return &Engine{Config: cfg}
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
		// Auto-detect based on directory contents
		if _, err := os.Stat(filepath.Join(sourcePath, "copilot-instructions.md")); err == nil {
			return parser.ParseCopilotDirectory(sourcePath)
		}
		if _, err := os.Stat(filepath.Join(sourcePath, "instructions")); err == nil {
			return parser.ParseCopilotDirectory(sourcePath)
		}
		if _, err := os.Stat(filepath.Join(sourcePath, ".agents")); err == nil {
			return parser.ParseAntigravityDirectory(sourcePath)
		}
		if _, err := os.Stat(filepath.Join(sourcePath, ".cursor")); err == nil {
			return parser.ParseCursorDirectory(filepath.Join(sourcePath, ".cursor", "rules"))
		}
		return nil, fmt.Errorf("unknown source format or unable to auto-detect: %s", from)
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
func (e *Engine) Convert(from string, to string, inputPath string, outputDir string) ([]string, *budget.AuditReport, error) {
	docs, err := e.ParseSource(from, inputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse source (%s): %w", from, err)
	}

	if len(docs) == 0 {
		return nil, nil, fmt.Errorf("no valid instruction documents found in %s", inputPath)
	}

	written, err := e.EmitTarget(to, docs, outputDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to emit target (%s): %w", to, err)
	}

	auditReport := budget.AuditDocuments(docs, 400)
	return written, auditReport, nil
}
