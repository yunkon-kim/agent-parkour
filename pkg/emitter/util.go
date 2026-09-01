package emitter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// SafeWriteFile writes data to targetPath with automatic timestamped backup if file already exists.
// Returns (backupFilePath, error). If no backup was needed, backupFilePath is empty.
func SafeWriteFile(targetPath string, data []byte) (string, error) {
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	backupCreated := ""
	if existingData, err := os.ReadFile(targetPath); err == nil {
		// If content is already identical, no backup needed
		if bytes.Equal(existingData, data) {
			return "", nil
		}
		// Create timestamped backup file: e.g. copilot-instructions.md.bak_20260901_230546
		timestamp := time.Now().Format("20060102_150405")
		backupPath := fmt.Sprintf("%s.bak_%s", targetPath, timestamp)
		if err := os.WriteFile(backupPath, existingData, 0644); err != nil {
			return "", fmt.Errorf("failed to backup existing file %s: %w", targetPath, err)
		}
		backupCreated = backupPath
	}

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return backupCreated, err
	}

	return backupCreated, nil
}
