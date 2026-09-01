package emitter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeWriteFileAndBackup(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "token-hop-backup-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	targetFile := filepath.Join(tmpDir, "copilot-instructions.md")

	// 1. Initial write (no existing file -> no backup)
	backup, err := SafeWriteFile(targetFile, []byte("# Initial Guidelines\n"))
	if err != nil {
		t.Fatalf("unexpected error on initial write: %v", err)
	}
	if backup != "" {
		t.Fatalf("expected no backup on initial write, got %s", backup)
	}

	// 2. Same content write (no backup needed)
	backup, err = SafeWriteFile(targetFile, []byte("# Initial Guidelines\n"))
	if err != nil {
		t.Fatalf("unexpected error on same content write: %v", err)
	}
	if backup != "" {
		t.Fatalf("expected no backup when content is unchanged, got %s", backup)
	}

	// 3. Modified content write (should create backup with timestamp)
	backup, err = SafeWriteFile(targetFile, []byte("# Updated Guidelines with Antigravity improvements\n"))
	if err != nil {
		t.Fatalf("unexpected error on overwrite: %v", err)
	}
	if backup == "" {
		t.Fatalf("expected backup file path, got empty string")
	}
	if !strings.Contains(backup, ".bak_") {
		t.Fatalf("expected backup path to contain '.bak_', got %s", backup)
	}

	// Verify backup file content matches initial content
	backupData, err := os.ReadFile(backup)
	if err != nil {
		t.Fatalf("failed to read backup file: %v", err)
	}
	if string(backupData) != "# Initial Guidelines\n" {
		t.Fatalf("backup content mismatch: got %q", string(backupData))
	}

	// Verify target file content has new updated content
	newTargetData, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("failed to read updated target file: %v", err)
	}
	if string(newTargetData) != "# Updated Guidelines with Antigravity improvements\n" {
		t.Fatalf("target content mismatch: got %q", string(newTargetData))
	}
}
