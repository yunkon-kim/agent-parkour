package describer

import (
	"strings"
	"testing"

	"github.com/yunkon-kim/agent-parkour/pkg/ir"
)

func TestMappingAnalyzer_Analyze(t *testing.T) {
	doc1 := ir.NewDocument("instruction-global", ir.TypeInstruction, "Global Instructions")
	doc1.Payload.MarkdownBody = "# Global Guidelines\nFollow clean architecture and write unit tests."
	doc1.Payload.RawSource = ".github/copilot-instructions.md"
	doc1.Activation.Mode = ir.ModeAlwaysOn

	doc2 := ir.NewDocument("rule-backend", ir.TypeRule, "Backend Rule")
	doc2.Payload.MarkdownBody = "Apply strict typing in Go."
	doc2.Payload.RawSource = ".github/instructions/backend.instructions.md"
	doc2.Activation.Mode = ir.ModeGlob
	doc2.Activation.Globs = []string{"**/*.go"}

	doc3 := ir.NewDocument("workflow-deploy", ir.TypeWorkflow, "Deploy Workflow")
	doc3.Payload.MarkdownBody = "Step 1: Build\nStep 2: Test\nStep 3: Deploy"
	doc3.Payload.RawSource = ".github/prompts/deploy.prompt.md"
	doc3.Activation.SlashCommand = "deploy"

	doc4 := ir.NewDocument("skill-architect", ir.TypeSkill, "Architect Skill")
	doc4.Payload.MarkdownBody = strings.Repeat("Long skill content for architecture design. ", 60)
	doc4.Payload.RawSource = ".github/agents/architect.agent.md"

	docs := []*ir.UADocument{doc1, doc2, doc3, doc4}

	analyzer := NewMappingAnalyzer(400)
	report := analyzer.Analyze(docs, "copilot", "antigravity", ".github", ".")

	if report.TotalSourceFiles != 4 {
		t.Errorf("expected 4 source files, got %d", report.TotalSourceFiles)
	}
	if report.TotalTargetFiles != 4 {
		t.Errorf("expected 4 target files, got %d", report.TotalTargetFiles)
	}
	if report.TotalTokens == 0 {
		t.Errorf("expected positive total tokens, got 0")
	}

	// Verify Note is present
	if len(report.Notes) == 0 || !strings.Contains(report.Notes[0], "Estimated token counts") {
		t.Errorf("expected note about token estimation, got %v", report.Notes)
	}

	// Verify Target Paths
	foundRule := false
	foundWorkflow := false
	foundSkill := false
	for _, item := range report.Items {
		if strings.Contains(item.TargetPath, ".agents/rules/rule-backend.md") {
			foundRule = true
		}
		if strings.Contains(item.TargetPath, ".agents/workflows/deploy.md") {
			foundWorkflow = true
		}
		if strings.Contains(item.TargetPath, ".agents/skills/architect/SKILL.md") {
			foundSkill = true
		}
	}

	if !foundRule {
		t.Errorf("rule target path was not mapped properly to .agents/rules")
	}
	if !foundWorkflow {
		t.Errorf("workflow target path was not mapped properly to .agents/workflows")
	}
	if !foundSkill {
		t.Errorf("skill target path was not mapped properly to .agents/skills")
	}
}

func TestFormatters(t *testing.T) {
	doc := ir.NewDocument("rule-test", ir.TypeRule, "Test Rule")
	doc.Payload.MarkdownBody = "Simple rule content."
	doc.Payload.RawSource = ".cursor/rules/test.mdc"

	analyzer := NewMappingAnalyzer(400)
	report := analyzer.Analyze([]*ir.UADocument{doc}, "cursor", "antigravity", ".cursor/rules", ".")

	// Test Table Formatter
	tableOut := FormatMappingReport(report, FormatTable)
	if !strings.Contains(tableOut, "agent-parkour describe") {
		t.Errorf("table output missing header: %s", tableOut)
	}
	if !strings.Contains(tableOut, "Estimated token counts") {
		t.Errorf("table output missing Note: %s", tableOut)
	}

	// Test Markdown Formatter
	mdOut := FormatMappingReport(report, FormatMarkdown)
	if !strings.Contains(mdOut, "| Source File | Entity Type | Target File |") {
		t.Errorf("markdown output missing table header: %s", mdOut)
	}

	// Test JSON Formatter
	jsonOut := FormatMappingReport(report, FormatJSON)
	if !strings.Contains(jsonOut, `"from_platform": "cursor"`) {
		t.Errorf("json output missing from_platform: %s", jsonOut)
	}
}

func TestBuildSpecMatrix(t *testing.T) {
	matrix := BuildSpecMatrix("copilot", "antigravity")
	if matrix.FromPlatform != "GitHub Copilot" {
		t.Errorf("expected GitHub Copilot, got %s", matrix.FromPlatform)
	}
	if matrix.ToPlatform != "Google Antigravity" {
		t.Errorf("expected Google Antigravity, got %s", matrix.ToPlatform)
	}
	if len(matrix.Items) != 6 {
		t.Errorf("expected 6 core entities, got %d", len(matrix.Items))
	}

	tableOut := FormatSpecMatrix(matrix, FormatTable)
	if !strings.Contains(tableOut, "Specification Matrix") {
		t.Errorf("spec matrix table missing title: %s", tableOut)
	}
}
