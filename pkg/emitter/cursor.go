package emitter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/agent-parkour/pkg/ir"
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

// Emit writes documents to the Cursor rules directory with automatic backup
func (e *CursorEmitter) Emit(docs []*ir.UADocument) ([]string, error) {
	var writtenFiles []string

	rulesDir := filepath.Join(e.BaseDir, ".cursor", "rules")

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

		backup, err := SafeWriteFile(filePath, []byte(out.String()))
		if err != nil {
			return writtenFiles, err
		}
		if backup != "" {
			fmt.Printf("      📦 Backed up existing file -> %s\n", backup)
		}
		writtenFiles = append(writtenFiles, filePath)
	}

	return writtenFiles, nil
}
