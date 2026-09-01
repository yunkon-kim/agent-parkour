package parser

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

// SplitFrontmatter splits markdown content into YAML frontmatter string and markdown body
func SplitFrontmatter(content string) (map[string]interface{}, string, error) {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "---") {
		return nil, content, nil
	}

	// Find closing ---
	rest := trimmed[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		// No closing frontmatter delimiter
		return nil, content, nil
	}

	fmStr := rest[:idx]
	body := strings.TrimPrefix(rest[idx+4:], "\n")

	var fmData map[string]interface{}
	err := yaml.Unmarshal([]byte(fmStr), &fmData)
	if err != nil {
		return nil, content, err
	}

	return fmData, strings.TrimSpace(body), nil
}

// ExtractFrontmatterAndUnmarshal unmarshals frontmatter into a target struct
func ExtractFrontmatterAndUnmarshal(content string, target interface{}) (string, error) {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "---") {
		return content, nil
	}

	rest := trimmed[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return content, nil
	}

	fmStr := rest[:idx]
	body := strings.TrimPrefix(rest[idx+4:], "\n")

	decoder := yaml.NewDecoder(bytes.NewReader([]byte(fmStr)))
	err := decoder.Decode(target)
	if err != nil {
		return content, err
	}

	return strings.TrimSpace(body), nil
}
