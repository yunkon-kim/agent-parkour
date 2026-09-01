package emitter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/ir"
	"gopkg.in/yaml.v3"
)

// CursorEmitter emits UA-IR documents to Cursor AI (.cursor/rules/*.mdc) format
type CursorEmitter struct {
	BaseDir string
}

// NewCursorEmitter creates a new Cursor emitter
func NewCursorEmitter(baseDir string) *CursorEmitter {
	return &CursorEmitter{BaseDir: baseDir}
}

// Emit writes documents to the Cursor rules directory
func (e *CursorEmitter) Emit(docs []*ir.UADocument) ([]string, error) {
	var writtenFiles []string

	rulesDir := filepath.Join(e.BaseDir, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return nil, err
	}

	for _, doc := range docs {
		fileName := SanitizeFileName(doc.Metadata.ID) + ".mdc"
		filePath := filepath.Join(rulesDir, fileName)

		fm := map[string]interface{}{
			"description": doc.Metadata.Description,
		}
		if fm["description"] == "" {
			fm["description"] = doc.Metadata.Name
		}

		if doc.Activation.Mode == ir.ModeAlwaysOn || doc.Metadata.Type == ir.TypeInstruction {
			fm["alwaysApply"] = true
			fm["globs"] = []string{"*"}
		} else if len(doc.Activation.Globs) > 0 {
			fm["alwaysApply"] = false
			fm["globs"] = doc.Activation.Globs
		} else {
			fm["alwaysApply"] = false
			fm["globs"] = []string{"*"}
		}

		fmBytes, err := yaml.Marshal(fm)
		if err != nil {
			return writtenFiles, err
		}

		var out strings.Builder
		out.WriteString("---\n")
		out.WriteString(string(fmBytes))
		out.WriteString("---\n\n")
		out.WriteString(RenderMarkdownWithTitle(doc.Metadata.Name, doc.Payload.MarkdownBody))

		if err := os.WriteFile(filePath, []byte(out.String()), 0644); err != nil {
			return writtenFiles, err
		}
		writtenFiles = append(writtenFiles, filePath)
	}

	return writtenFiles, nil
}


