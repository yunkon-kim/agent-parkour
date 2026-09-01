package audit

import (
	"strings"
	"testing"

	"github.com/yunkon-kim/token-hop/pkg/ir"
)

func TestEstimateTokens(t *testing.T) {
	shortText := "Hello world this is a test guideline"
	tokens := EstimateTokens(shortText)
	if tokens < 1 || tokens > 15 {
		t.Fatalf("expected reasonable token count for short text, got %d", tokens)
	}

	emptyText := ""
	tokens = EstimateTokens(emptyText)
	if tokens != 0 {
		t.Fatalf("expected 0 tokens for empty text, got %d", tokens)
	}
}

func TestAuditDocuments(t *testing.T) {
	doc1 := ir.NewDocument("small-rule", ir.TypeRule, "Small Rule")
	doc1.Payload.MarkdownBody = "This is a small rule."

	doc2 := ir.NewDocument("huge-rule", ir.TypeRule, "Huge Rule")
	// Make a long string with repeated sentences > 400 tokens (~3000 chars)
	doc2.Payload.MarkdownBody = strings.Repeat("This is a detailed and long guideline rule containing many tokens and instructions for the agent. ", 50)

	docs := []*ir.UADocument{doc1, doc2}
	report := AuditDocuments(docs, 400)

	if report.TotalDocuments != 2 {
		t.Fatalf("expected 2 documents, got %d", report.TotalDocuments)
	}
	if report.Items[0].ExceedsLimit {
		t.Fatalf("expected small-rule not to exceed limit")
	}
	if !report.Items[1].ExceedsLimit {
		t.Fatalf("expected huge-rule to exceed limit, tokens = %d", report.Items[1].Tokens)
	}
	if report.Items[1].Recommendation == "" {
		t.Fatalf("expected recommendation for huge-rule, got empty")
	}
}
