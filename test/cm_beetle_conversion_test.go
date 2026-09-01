package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yunkon-kim/token-hop/pkg/budget"
	"github.com/yunkon-kim/token-hop/pkg/emitter"
	"github.com/yunkon-kim/token-hop/pkg/parser"
)

func TestCmBeetleCopilotToAntigravityConversion(t *testing.T) {
	fixtureDir := filepath.Join("fixtures", "cm-beetle", ".github")
	if _, err := os.Stat(fixtureDir); os.IsNotExist(err) {
		t.Skipf("Fixture directory %s not found, skipping", fixtureDir)
	}

	// 1. Parse GitHub Copilot configuration
	docs, err := parser.ParseCopilotDirectory(fixtureDir)
	if err != nil {
		t.Fatalf("Failed to parse cm-beetle Copilot directory: %v", err)
	}

	if len(docs) == 0 {
		t.Fatalf("Expected documents to be parsed from cm-beetle, got 0")
	}

	t.Logf("Parsed %d documents from cm-beetle .github directory", len(docs))

	// Verify key documents exist
	var foundGlobal, foundGoRule, foundApiGuide, foundSyncTb bool
	for _, doc := range docs {
		t.Logf("  - [%s] %s (%s, Mode: %s)", doc.Metadata.Type, doc.Metadata.ID, doc.Metadata.Name, doc.Activation.Mode)
		if doc.Metadata.ID == "instruction-global" {
			foundGlobal = true
		}
		if doc.Metadata.ID == "rule-go" {
			foundGoRule = true
			if len(doc.Activation.Globs) == 0 || doc.Activation.Globs[0] != "**/*.go" {
				t.Errorf("Expected rule-go to have glob '**/*.go', got %v", doc.Activation.Globs)
			}
		}
		if doc.Metadata.ID == "workflow-api-guide" {
			foundApiGuide = true
		}
		if doc.Metadata.ID == "workflow-sync-tb" {
			foundSyncTb = true
		}
	}

	if !foundGlobal {
		t.Errorf("Expected instruction-global to be parsed")
	}
	if !foundGoRule {
		t.Errorf("Expected rule-go to be parsed")
	}
	if !foundApiGuide {
		t.Errorf("Expected workflow-api-guide to be parsed")
	}
	if !foundSyncTb {
		t.Errorf("Expected workflow-sync-tb to be parsed")
	}

	// 2. Emit to Antigravity format in temporary directory
	tempOutDir, err := os.MkdirTemp("", "token-hop-antigravity-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempOutDir)

	agEmitter := emitter.NewAntigravityEmitter(tempOutDir)
	writtenFiles, err := agEmitter.Emit(docs)
	if err != nil {
		t.Fatalf("Failed to emit Antigravity files: %v", err)
	}

	t.Logf("Emitted %d Antigravity files to %s", len(writtenFiles), tempOutDir)

	// Verify generated Antigravity structure
	expectedAgentsMd := filepath.Join(tempOutDir, "AGENTS.md")
	if _, err := os.Stat(expectedAgentsMd); os.IsNotExist(err) {
		t.Errorf("Expected AGENTS.md to be created")
	}

	expectedGoRule := filepath.Join(tempOutDir, ".agents", "rules", "rule-go.md")
	data, err := os.ReadFile(expectedGoRule)
	if err != nil {
		t.Errorf("Expected .agents/rules/rule-go.md to exist: %v", err)
	} else {
		content := string(data)
		if !strings.Contains(content, "globs:") || !strings.Contains(content, "**/*.go") {
			t.Errorf(".agents/rules/rule-go.md missing glob frontmatter: %s", content[:min(200, len(content))])
		}
	}

	expectedWorkflow := filepath.Join(tempOutDir, ".agents", "workflows", "api-guide.md")
	if _, err := os.Stat(expectedWorkflow); os.IsNotExist(err) {
		t.Errorf("Expected .agents/workflows/api-guide.md to exist")
	}

	// 3. Context Budget Audit Test
	report := budget.AuditDocuments(docs, 400)
	t.Logf("Audit Report: Total Docs: %d, Total Tokens: ~%d, Total Chars: %d",
		report.TotalDocuments, report.TotalTokens, report.TotalCharacters)

	for _, item := range report.Items {
		if item.ExceedsBudget {
			t.Logf("  ⚠️  Large item [%s] %s: ~%d tokens (Chars: %d) -> %s",
				item.Type, item.ID, item.Tokens, item.Characters, item.Recommendation)
		}
	}
}

func TestCmBeetleCopilotToCursorConversion(t *testing.T) {
	fixtureDir := filepath.Join("fixtures", "cm-beetle", ".github")
	docs, err := parser.ParseCopilotDirectory(fixtureDir)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	tempOutDir, err := os.MkdirTemp("", "token-hop-cursor-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempOutDir)

	curEmitter := emitter.NewCursorEmitter(tempOutDir)
	writtenFiles, err := curEmitter.Emit(docs)
	if err != nil {
		t.Fatalf("Failed to emit Cursor files: %v", err)
	}

	t.Logf("Emitted %d Cursor rules (.mdc) to %s", len(writtenFiles), tempOutDir)

	expectedMdc := filepath.Join(tempOutDir, ".cursor", "rules", "rule-go.mdc")
	data, err := os.ReadFile(expectedMdc)
	if err != nil {
		t.Errorf("Expected rule-go.mdc to exist: %v", err)
	} else {
		content := string(data)
		if !strings.Contains(content, "globs:") || !strings.Contains(content, "**/*.go") {
			t.Errorf("rule-go.mdc missing globs: %s", content[:min(200, len(content))])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
