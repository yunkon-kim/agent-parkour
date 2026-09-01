package ir

import "gopkg.in/yaml.v3"

// EntityType represents the 6 core dimensions of AI instructions
type EntityType string

const (
	TypeRule        EntityType = "Rule"
	TypeSkill       EntityType = "Skill"
	TypeWorkflow    EntityType = "Workflow"
	TypeInstruction EntityType = "Instruction"
	TypePrompt      EntityType = "Prompt"
	TypeAgent       EntityType = "Agent"
)

// ActivationMode defines how and when the instruction is activated into context
type ActivationMode string

const (
	ModeAlwaysOn      ActivationMode = "AlwaysOn"
	ModeGlob          ActivationMode = "Glob"
	ModeModelDecision ActivationMode = "ModelDecision"
	ModeOnDemand      ActivationMode = "OnDemand"
	ModeManualMention ActivationMode = "ManualMention"
)

// DecomposeStrategy defines the strategy when context budget is exceeded
type DecomposeStrategy string

const (
	StrategyJITSkill  DecomposeStrategy = "JIT_Skill"
	StrategyTruncate  DecomposeStrategy = "Truncate"
	StrategySplitGlob DecomposeStrategy = "Split_Glob"
)

// UAMetadata represents standard YAML metadata for an instruction
type UAMetadata struct {
	ID          string     `yaml:"id"`
	Type        EntityType `yaml:"type"`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description,omitempty"`
	Version     string     `yaml:"version,omitempty"`
	Author      string     `yaml:"author,omitempty"`
	Tags        []string   `yaml:"tags,omitempty"`
}

// Activation controls condition-based loading into context
type Activation struct {
	Mode         ActivationMode `yaml:"mode"`
	Globs        []string       `yaml:"globs,omitempty"`
	ExcludeGlobs []string       `yaml:"exclude_globs,omitempty"`
	Triggers     []string       `yaml:"triggers,omitempty"`
	SlashCommand string         `yaml:"slash_command,omitempty"`
}

// ContextBudget defines token and character limits
type ContextBudget struct {
	Priority          string            `yaml:"priority,omitempty"` // Critical, High, Medium, Low
	MaxTokens         int               `yaml:"max_tokens,omitempty"`
	MaxCharacters     int               `yaml:"max_characters,omitempty"`
	DecomposeStrategy DecomposeStrategy `yaml:"decompose_strategy,omitempty"`
}

// Bindings contains tool and MCP bindings
type Bindings struct {
	AllowedTools   []string `yaml:"allowed_tools,omitempty"`
	MCPServers     []string `yaml:"mcp_servers,omitempty"`
	SubagentParent string   `yaml:"subagent_parent,omitempty"`
}

// Payload contains the Markdown content body
type Payload struct {
	MarkdownBody string `yaml:"markdown_body"`
	RawSource    string `yaml:"-"`
}

// UADocument represents a single Universal Agent IR AST Node
type UADocument struct {
	UAVersion     string        `yaml:"ua_version"`
	Metadata      UAMetadata    `yaml:"metadata"`
	Activation    Activation    `yaml:"activation"`
	ContextBudget ContextBudget `yaml:"context_budget,omitempty"`
	Bindings      Bindings      `yaml:"bindings,omitempty"`
	Payload       Payload       `yaml:"payload"`
}

// ProjectManifest contains all UADocuments parsed from a project (SSOT)
type ProjectManifest struct {
	Name        string        `yaml:"name"`
	Version     string        `yaml:"version"`
	Description string        `yaml:"description,omitempty"`
	Documents   []*UADocument `yaml:"documents"`
}

// NewDocument creates a new UADocument with defaults
func NewDocument(id string, entityType EntityType, name string) *UADocument {
	return &UADocument{
		UAVersion: "1.0.0",
		Metadata: UAMetadata{
			ID:   id,
			Type: entityType,
			Name: name,
		},
		Activation: Activation{
			Mode: ModeAlwaysOn,
		},
		ContextBudget: ContextBudget{
			Priority:          "High",
			MaxTokens:         400,
			MaxCharacters:     1800,
			DecomposeStrategy: StrategyJITSkill,
		},
	}
}

// ToYAML serializes UADocument into YAML
func (d *UADocument) ToYAML() (string, error) {
	bytes, err := yaml.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
