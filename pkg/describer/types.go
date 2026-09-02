package describer

import (
	"github.com/yunkon-kim/token-hop/pkg/ir"
)

// OutputFormat defines the format for rendering description output
type OutputFormat string

const (
	FormatTable    OutputFormat = "table"
	FormatMarkdown OutputFormat = "markdown"
	FormatJSON     OutputFormat = "json"
)

// PlanItem represents a single file mapping transformation plan
type PlanItem struct {
	SourcePath     string        `json:"source_path"`
	SourceType     ir.EntityType `json:"source_type"`
	SourceTrigger  string        `json:"source_trigger"`
	TargetPath     string        `json:"target_path"`
	TargetType     ir.EntityType `json:"target_type"`
	TargetTrigger  string        `json:"target_trigger"`
	Characters     int           `json:"characters"`
	Tokens         int           `json:"tokens"`
	IsOversized    bool          `json:"is_oversized"`
	Action         string        `json:"action"`
	Recommendation string        `json:"recommendation,omitempty"`
}

// MappingReport represents the complete preview plan for converting between platforms
type MappingReport struct {
	FromPlatform     string     `json:"from_platform"`
	ToPlatform       string     `json:"to_platform"`
	SourceDir        string     `json:"source_dir"`
	OutputDir        string     `json:"output_dir"`
	TotalSourceFiles int        `json:"total_source_files"`
	TotalTargetFiles int        `json:"total_target_files"`
	TotalTokens      int        `json:"total_tokens"`
	TotalCharacters  int        `json:"total_characters"`
	AlwaysOnTokens   int        `json:"always_on_tokens"`
	OnDemandTokens   int        `json:"on_demand_tokens"`
	OversizedCount   int        `json:"oversized_count"`
	Items            []PlanItem `json:"items"`
	Notes            []string   `json:"notes,omitempty"`
}

// EntitySpec defines how a specific entity is configured in a platform
type EntitySpec struct {
	EntityType  ir.EntityType `json:"entity_type"`
	Location    string        `json:"location"`
	Syntax      string        `json:"syntax"`
	TriggerMode string        `json:"trigger_mode"`
	Description string        `json:"description"`
}

// PlatformSpec defines the overall configuration specification for an AI platform
type PlatformSpec struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	RootDir     string       `json:"root_dir"`
	SSOTFile    string       `json:"ssot_file,omitempty"`
	Entities    []EntitySpec `json:"entities"`
	Description string       `json:"description"`
}

// SpecMatrixItem represents a conceptual mapping row between two platforms
type SpecMatrixItem struct {
	EntityType     ir.EntityType `json:"entity_type"`
	SourceLocation string        `json:"source_location"`
	SourceSyntax   string        `json:"source_syntax"`
	TargetLocation string        `json:"target_location"`
	TargetSyntax   string        `json:"target_syntax"`
	TargetBehavior string        `json:"target_behavior"`
}

// SpecMatrixReport represents the conceptual specification comparison report
type SpecMatrixReport struct {
	FromPlatform string           `json:"from_platform"`
	ToPlatform   string           `json:"to_platform"`
	Items        []SpecMatrixItem `json:"items"`
}
