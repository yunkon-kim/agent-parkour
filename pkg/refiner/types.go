package refiner

import "github.com/yunkon-kim/token-hop/pkg/ir"

// FileItem represents an individual converted file included in a refinement prompt
type FileItem struct {
	FilePath    string        `json:"file_path"`
	EntityType  ir.EntityType `json:"entity_type"`
	FileContent string        `json:"file_content"`
	Characters  int           `json:"characters"`
	Tokens      int           `json:"tokens"`
	IsOversized bool          `json:"is_oversized"`
}

// PromptContext contains all parameters required to render a target refinement prompt
type PromptContext struct {
	TargetPlatform string     `json:"target_platform"`
	PlatformName   string     `json:"platform_name"`
	RootDir        string     `json:"root_dir"`
	Files          []FileItem `json:"files"`
	TotalTokens    int        `json:"total_tokens"`
	TotalFiles     int        `json:"total_files"`
	CustomGuidance string     `json:"custom_guidance,omitempty"`
}

// RefineOptions controls how refinement prompts are constructed
type RefineOptions struct {
	TargetPlatform string
	FilePath       string
	DirectoryPath  string
	OutputPath     string
	MaxTokens      int
	CustomGuidance string
}
