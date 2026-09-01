package ai

import (
	"fmt"
)

// BuildDecomposePrompt constructs the system/user instruction to decompose a large rule into modular parts.
func BuildDecomposePrompt(title, content string, maxTokens int) string {
	return fmt.Sprintf(`You are an expert AI prompt and rule modularization compiler.
Analyze the following large AI rule/workflow titled "%s" and decompose it into smaller, modular sub-rules or JIT skills.

Rules for decomposition:
1. Extract distinct concerns (e.g. Models, Coding Rules, Git Workflow, API Sync) into separate sub-rules.
2. If a section is large and only needed during specific multi-step tasks, mark "is_skill": true (JIT on-demand skill).
3. Generate concise globs (e.g. ["**/*.go"]) if the rule is file-specific.
4. Each sub-rule must be fully self-contained and clear.
5. Return ONLY a valid JSON object matching the schema below, without any markdown fences.

JSON Schema:
{
  "original_title": "%s",
  "summary": "Brief explanation of how the rule was modularized",
  "sub_rules": [
    {
      "id": "kebab-case-id",
      "title": "Clear Sub-Rule Title",
      "description": "Concise 1-sentence purpose",
      "globs": ["**/*.go"],
      "is_skill": false,
      "content": "Full markdown content of this modular sub-rule..."
    }
  ]
}

---
Content to decompose:
%s
`, title, title, content)
}

// BuildDescriptionPrompt constructs the instruction to generate a 1-sentence activation description.
func BuildDescriptionPrompt(content string) string {
	return fmt.Sprintf(`You are an AI context optimization assistant.
Generate a concise, 1-sentence description (under 120 characters) explaining what this AI rule or workflow does and when it should be activated.
Respond ONLY with the 1-sentence description text, without quotation marks or explanations.

Content:
%s
`, content)
}
