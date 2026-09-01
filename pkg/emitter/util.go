package emitter

import (
	"fmt"
	"strings"
)

// SanitizeFileName cleans up an ID to make it safe for file names
func SanitizeFileName(id string) string {
	id = strings.ToLower(id)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "/", "-")
	id = strings.ReplaceAll(id, "\\", "-")
	return id
}

// RenderMarkdownWithTitle ensures title is rendered cleanly without duplicate # headings
func RenderMarkdownWithTitle(title, body string) string {
	trimmed := strings.TrimSpace(body)
	if strings.HasPrefix(trimmed, "# ") {
		return trimmed + "\n"
	}
	return fmt.Sprintf("# %s\n\n%s\n", title, trimmed)
}
