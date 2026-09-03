package refiner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yunkon-kim/agent-parkour/pkg/ir"
)

func TestGenerator_GenerateFromDocs(t *testing.T) {
	doc1 := ir.NewDocument("workflow-git-commit", ir.TypeWorkflow, "Git Commit")
	doc1.Payload.MarkdownBody = "### Step 1: Check diff\nUse get_changed_files to see git status."
	doc1.Payload.RawSource = ".agents/workflows/git-commit.md"

	doc2 := ir.NewDocument("rule-go-style", ir.TypeRule, "Go Style")
	doc2.Payload.MarkdownBody = strings.Repeat("Always use explicit error handling. ", 30)
	doc2.Payload.RawSource = ".agents/rules/go-style.md"

	docs := []*ir.UADocument{doc1, doc2}

	gen := NewGenerator(400)
	prompt, err := gen.GenerateFromDocs("antigravity", docs, "Special guidance for backend services.")
	if err != nil {
		t.Fatalf("GenerateFromDocs failed: %v", err)
	}

	if !strings.Contains(prompt, "Google Antigravity Configuration Optimizer") {
		t.Errorf("prompt missing Antigravity header: %s", prompt)
	}
	if !strings.Contains(prompt, "// turbo") {
		t.Errorf("prompt missing // turbo guideline: %s", prompt)
	}
	if !strings.Contains(prompt, "Special guidance for backend services.") {
		t.Errorf("prompt missing custom guidance: %s", prompt)
	}
	if !strings.Contains(prompt, ".agents/workflows/git-commit.md") {
		t.Errorf("prompt missing target file path: %s", prompt)
	}
}

func TestGenerator_GenerateFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-rule.mdc")
	err := os.WriteFile(testFile, []byte("---\ndescription: Test rule\n---\n# Rule\nDo not use any."), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	gen := NewGenerator(400)
	prompt, err := gen.GenerateFromFile("cursor", testFile, "")
	if err != nil {
		t.Fatalf("GenerateFromFile failed: %v", err)
	}

	if !strings.Contains(prompt, "Cursor AI MDC Rule Optimizer") {
		t.Errorf("prompt missing Cursor header: %s", prompt)
	}
	if !strings.Contains(prompt, "@file") || !strings.Contains(prompt, "@rule") {
		t.Errorf("prompt missing Cursor reference guidelines: %s", prompt)
	}
}

func TestGenerator_ClaudeAndCopilot(t *testing.T) {
	gen := NewGenerator(400)

	doc := ir.NewDocument("instruction-main", ir.TypeInstruction, "Main Instructions")
	doc.Payload.MarkdownBody = "# Project Architecture\nFollow Clean Architecture."
	doc.Payload.RawSource = "CLAUDE.md"

	claudePrompt, err := gen.GenerateFromDocs("claude", []*ir.UADocument{doc}, "")
	if err != nil {
		t.Fatalf("Claude prompt failed: %v", err)
	}
	if !strings.Contains(claudePrompt, "Claude Code (CLAUDE.md) Optimizer") {
		t.Errorf("missing Claude header: %s", claudePrompt)
	}
	if !strings.Contains(claudePrompt, "200 lines") {
		t.Errorf("missing 200 lines guideline for Claude: %s", claudePrompt)
	}

	copilotPrompt, err := gen.GenerateFromDocs("copilot", []*ir.UADocument{doc}, "")
	if err != nil {
		t.Fatalf("Copilot prompt failed: %v", err)
	}
	if !strings.Contains(copilotPrompt, "GitHub Copilot Instruction & Prompt Optimizer") {
		t.Errorf("missing Copilot header: %s", copilotPrompt)
	}
}
