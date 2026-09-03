package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yunkon-kim/agent-parkour/pkg/ir"
)

func TestParseApplyTo(t *testing.T) {
	// 1. Single string
	res1 := parseApplyTo("**/*.go")
	if len(res1) != 1 || res1[0] != "**/*.go" {
		t.Errorf("expected ['**/*.go'], got %v", res1)
	}

	// 2. YAML slice of interfaces (e.g. from YAML array)
	res2 := parseApplyTo([]interface{}{"ui/**", "web/**"})
	if len(res2) != 2 || res2[0] != "ui/**" || res2[1] != "web/**" {
		t.Errorf("expected ['ui/**', 'web/**'], got %v", res2)
	}

	// 3. Empty or nil
	if len(parseApplyTo(nil)) != 0 {
		t.Errorf("expected empty for nil")
	}
	if len(parseApplyTo("")) != 0 {
		t.Errorf("expected empty for empty string")
	}
}

func TestInferGlobFromID(t *testing.T) {
	if inferGlobFromID("go")[0] != "**/*.go" {
		t.Errorf("expected **/*.go for 'go'")
	}
	if inferGlobFromID("ui")[0] != "ui/**" {
		t.Errorf("expected ui/** for 'ui'")
	}
	if inferGlobFromID("analyzer")[0] != "analyzer/**" {
		t.Errorf("expected analyzer/** for 'analyzer'")
	}
}

func TestCopilotParser_DirectoryFirst(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "parkour-copilot-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create structure:
	// .github/copilot-instructions.md
	// .github/instructions/ui.instructions.md (with YAML array applyTo)
	// .github/instructions/no-frontmatter.instructions.md
	// .github/prompts/deploy.prompt.md
	ghDir := filepath.Join(tempDir, ".github")
	instDir := filepath.Join(ghDir, "instructions")
	promptsDir := filepath.Join(ghDir, "prompts")
	os.MkdirAll(instDir, 0755)
	os.MkdirAll(promptsDir, 0755)

	os.WriteFile(filepath.Join(ghDir, "copilot-instructions.md"), []byte("# Global SSOT\nRules here."), 0644)
	os.WriteFile(filepath.Join(instDir, "ui.instructions.md"), []byte("---\napplyTo:\n  - ui/**\n---\n# UI Rules"), 0644)
	os.WriteFile(filepath.Join(instDir, "backend.instructions.md"), []byte("# Backend Rules without frontmatter"), 0644)
	os.WriteFile(filepath.Join(promptsDir, "deploy.prompt.md"), []byte("---\ndescription: Deploy app\n---\nDeploy steps"), 0644)

	// Test parsing by passing root tempDir (not .github directly)
	docs, err := ParseCopilotDirectory(tempDir)
	if err != nil {
		t.Fatalf("ParseCopilotDirectory failed on repo root: %v", err)
	}

	if len(docs) != 4 {
		t.Fatalf("expected 4 documents, got %d", len(docs))
	}

	var foundGlobal, foundUI, foundBackend, foundDeploy bool
	for _, doc := range docs {
		switch doc.Metadata.ID {
		case "instruction-global":
			foundGlobal = true
			if doc.Metadata.Type != ir.TypeInstruction || doc.Activation.Mode != ir.ModeAlwaysOn {
				t.Errorf("instruction-global should be Instruction & AlwaysOn, got %s / %s", doc.Metadata.Type, doc.Activation.Mode)
			}
		case "rule-ui":
			foundUI = true
			if doc.Metadata.Type != ir.TypeRule {
				t.Errorf("rule-ui should be Rule, got %s", doc.Metadata.Type)
			}
			if doc.Activation.Mode != ir.ModeGlob {
				t.Errorf("rule-ui should have ModeGlob, got %s", doc.Activation.Mode)
			}
			if len(doc.Activation.Globs) == 0 || doc.Activation.Globs[0] != "ui/**" {
				t.Errorf("rule-ui should have glob 'ui/**', got %v", doc.Activation.Globs)
			}
		case "rule-backend":
			foundBackend = true
			if doc.Metadata.Type != ir.TypeRule || doc.Activation.Mode != ir.ModeGlob {
				t.Errorf("rule-backend should be Rule & ModeGlob, got %s / %s", doc.Metadata.Type, doc.Activation.Mode)
			}
			if len(doc.Activation.Globs) == 0 || doc.Activation.Globs[0] != "backend/**" {
				t.Errorf("rule-backend should infer glob 'backend/**', got %v", doc.Activation.Globs)
			}
		case "prompt-deploy":
			foundDeploy = true
			if doc.Metadata.Type != ir.TypePrompt {
				t.Errorf("prompt-deploy should be Prompt, got %s", doc.Metadata.Type)
			}
			if doc.Activation.Mode != ir.ModeOnDemand {
				t.Errorf("prompt-deploy should be OnDemand, got %s", doc.Activation.Mode)
			}
			if doc.Activation.SlashCommand != "deploy" {
				t.Errorf("prompt-deploy command should be 'deploy', got %s", doc.Activation.SlashCommand)
			}
		}
	}

	if !foundGlobal || !foundUI || !foundBackend || !foundDeploy {
		t.Errorf("missing expected documents: global=%v, ui=%v, backend=%v, deploy=%v", foundGlobal, foundUI, foundBackend, foundDeploy)
	}
}

func TestAntigravityParser_AgentsSubdir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "parkour-ag-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	agDir := filepath.Join(tempDir, ".agents")
	rulesDir := filepath.Join(agDir, "rules")
	workflowsDir := filepath.Join(agDir, "workflows")
	os.MkdirAll(rulesDir, 0755)
	os.MkdirAll(workflowsDir, 0755)

	// .agents/AGENTS.md
	os.WriteFile(filepath.Join(agDir, "AGENTS.md"), []byte("# Project SSOT"), 0644)
	// .agents/rules/api.md
	os.WriteFile(filepath.Join(rulesDir, "api.md"), []byte("# API Guidelines"), 0644)
	// .agents/workflows/git-commit.md
	os.WriteFile(filepath.Join(workflowsDir, "git-commit.md"), []byte("---\ndescription: Commit\n---\nSteps"), 0644)

	docs, err := ParseAntigravityDirectory(tempDir)
	if err != nil {
		t.Fatalf("ParseAntigravityDirectory failed: %v", err)
	}

	if len(docs) != 3 {
		t.Fatalf("expected 3 documents, got %d", len(docs))
	}

	for _, doc := range docs {
		if doc.Metadata.ID == "rule-api" {
			if doc.Activation.Mode != ir.ModeGlob {
				t.Errorf("expected rule-api to be ModeGlob (inferred), got %s", doc.Activation.Mode)
			}
		}
	}
}
