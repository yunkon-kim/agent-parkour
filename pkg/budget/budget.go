package budget

import (
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/ir"
)

// AuditReport contains token and character metrics for documents
type AuditItem struct {
	ID            string        `json:"id"`
	Type          ir.EntityType `json:"type"`
	Characters    int           `json:"characters"`
	Tokens        int           `json:"tokens"`
	ExceedsBudget bool          `json:"exceeds_budget"`
	Recommendation string       `json:"recommendation,omitempty"`
}

type AuditReport struct {
	TotalDocuments int         `json:"total_documents"`
	TotalTokens    int         `json:"total_tokens"`
	TotalCharacters int        `json:"total_characters"`
	Items          []AuditItem `json:"items"`
}

// EstimateTokens calculates estimated token count from text (approx 4 chars/token or word split)
func EstimateTokens(text string) int {
	charCount := len(text)
	wordCount := len(strings.Fields(text))
	
	// Blended estimation: (chars / 3.8 + words * 1.3) / 2
	est := int((float64(charCount)/3.8 + float64(wordCount)*1.3) / 2.0)
	if est < 1 && charCount > 0 {
		return 1
	}
	return est
}

// AuditDocuments audits a slice of UA-IR documents
func AuditDocuments(docs []*ir.UADocument, maxTokensPerRule int) *AuditReport {
	if maxTokensPerRule <= 0 {
		maxTokensPerRule = 400
	}

	report := &AuditReport{
		TotalDocuments: len(docs),
	}

	for _, doc := range docs {
		chars := len(doc.Payload.MarkdownBody)
		tokens := EstimateTokens(doc.Payload.MarkdownBody)

		report.TotalCharacters += chars
		report.TotalTokens += tokens

		exceeds := tokens > maxTokensPerRule
		rec := ""
		if exceeds {
			if doc.Metadata.Type == ir.TypeRule {
				rec = "Split into on-demand JIT Skill (.agent/skills/) or apply stricter Glob pattern"
			} else if doc.Metadata.Type == ir.TypeWorkflow {
				rec = "Decompose into chained step artifacts"
			} else {
				rec = "Move detailed reference specs into documentation artifacts"
			}
		}

		report.Items = append(report.Items, AuditItem{
			ID:             doc.Metadata.ID,
			Type:           doc.Metadata.Type,
			Characters:     chars,
			Tokens:         tokens,
			ExceedsBudget:  exceeds,
			Recommendation: rec,
		})
	}

	return report
}
