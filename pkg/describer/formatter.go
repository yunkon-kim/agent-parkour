package describer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatMappingReport renders a MappingReport in the specified OutputFormat
func FormatMappingReport(report *MappingReport, format OutputFormat) string {
	switch format {
	case FormatMarkdown:
		return formatMarkdownTable(report)
	case FormatJSON:
		bytes, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Sprintf("{\"error\": %q}", err.Error())
		}
		return string(bytes)
	case FormatTable:
		fallthrough
	default:
		return formatCLITable(report)
	}
}

// FormatSpecMatrix renders a SpecMatrixReport in the specified OutputFormat
func FormatSpecMatrix(matrix *SpecMatrixReport, format OutputFormat) string {
	switch format {
	case FormatMarkdown:
		return formatSpecMatrixMarkdown(matrix)
	case FormatJSON:
		bytes, err := json.MarshalIndent(matrix, "", "  ")
		if err != nil {
			return fmt.Sprintf("{\"error\": %q}", err.Error())
		}
		return string(bytes)
	case FormatTable:
		fallthrough
	default:
		return formatSpecMatrixCLITable(matrix)
	}
}

func formatCLITable(report *MappingReport) string {
	var sb strings.Builder

	sb.WriteString("🔍 [token-hop describe] Cross-Agent Configuration Mapping Plan\n")
	sb.WriteString(fmt.Sprintf("   • Source Platform : %s (%s)\n", strings.ToUpper(report.FromPlatform), report.SourceDir))
	sb.WriteString(fmt.Sprintf("   • Target Platform : %s (%s)\n", strings.ToUpper(report.ToPlatform), report.OutputDir))
	sb.WriteString(fmt.Sprintf("   • Detected Files  : %d document(s)\n\n", report.TotalSourceFiles))

	if len(report.Items) == 0 {
		sb.WriteString("   (No instruction documents found to map)\n")
		return sb.String()
	}

	headers := []string{"Source File", "Type", "Target File", "Est. Tokens", "Action"}
	colWidths := []int{24, 10, 26, 12, 26}

	// Calculate dynamic widths based on content (capped for terminal readability)
	for _, item := range report.Items {
		if len(item.SourcePath) > colWidths[0] {
			colWidths[0] = min(len(item.SourcePath), 34)
		}
		if len(string(item.SourceType)) > colWidths[1] {
			colWidths[1] = min(len(string(item.SourceType)), 12)
		}
		if len(item.TargetPath) > colWidths[2] {
			colWidths[2] = min(len(item.TargetPath), 36)
		}
		tokStr := fmt.Sprintf("~%d tok", item.Tokens)
		if item.IsOversized {
			tokStr += " [!]"
		}
		if len(tokStr) > colWidths[3] {
			colWidths[3] = len(tokStr)
		}
		if len(item.Action) > colWidths[4] {
			colWidths[4] = min(len(item.Action), 34)
		}
	}

	topBorder := buildBorder("┌", "┬", "┐", "─", colWidths)
	midBorder := buildBorder("├", "┼", "┤", "─", colWidths)
	botBorder := buildBorder("└", "┴", "┘", "─", colWidths)

	sb.WriteString(topBorder)
	sb.WriteByte('\n')
	sb.WriteString(buildRow(headers, colWidths))
	sb.WriteByte('\n')
	sb.WriteString(midBorder)
	sb.WriteByte('\n')

	for _, item := range report.Items {
		tokStr := fmt.Sprintf("~%d tok", item.Tokens)
		if item.IsOversized {
			tokStr += " [!]"
		}
		row := []string{
			truncateString(item.SourcePath, colWidths[0]),
			string(item.SourceType),
			truncateString(item.TargetPath, colWidths[2]),
			tokStr,
			truncateString(item.Action, colWidths[4]),
		}
		sb.WriteString(buildRow(row, colWidths))
		sb.WriteByte('\n')
	}
	sb.WriteString(botBorder)
	sb.WriteString("\n\n")

	// Summary statistics
	sb.WriteString("📊 Summary & Context Metrics:\n")
	sb.WriteString(fmt.Sprintf("   • Total Source Files : %d files\n", report.TotalSourceFiles))
	sb.WriteString(fmt.Sprintf("   • Total Target Files : %d files\n", report.TotalTargetFiles))
	sb.WriteString(fmt.Sprintf("   • Total Context Size : ~%d tokens (%d characters)\n", report.TotalTokens, report.TotalCharacters))
	sb.WriteString(fmt.Sprintf("   • Always-On (Turn-0) : ~%d tokens (Injected in every request)\n", report.AlwaysOnTokens))
	sb.WriteString(fmt.Sprintf("   • On-Demand (JIT)    : ~%d tokens (Loaded only when needed)\n", report.OnDemandTokens))

	if report.OversizedCount > 0 {
		sb.WriteString(fmt.Sprintf("   ⚠️  Oversized Rules  : %d file(s) exceed recommended limit (>400 tokens)\n", report.OversizedCount))
		for _, item := range report.Items {
			if item.IsOversized && item.Recommendation != "" {
				sb.WriteString(fmt.Sprintf("      └─ [%s] %s\n", item.SourcePath, item.Recommendation))
			}
		}
	}

	// Notes & Next Steps
	sb.WriteString("\n💡 Note: Estimated token counts (~Tokens) are calculated locally without external API calls or cost, indicating baseline prompt context overhead per conversation turn.\n\n")
	sb.WriteString(fmt.Sprintf("👉 Next Step:\n   Run 'thop convert --from %s --to %s' to execute this transformation.\n", report.FromPlatform, report.ToPlatform))

	return sb.String()
}

func formatMarkdownTable(report *MappingReport) string {
	var sb strings.Builder

	sb.WriteString("### 🔍 Cross-Agent Configuration Mapping Plan\n\n")
	sb.WriteString(fmt.Sprintf("- **Source Platform**: `%s` (`%s`)\n", report.FromPlatform, report.SourceDir))
	sb.WriteString(fmt.Sprintf("- **Target Platform**: `%s` (`%s`)\n", report.ToPlatform, report.OutputDir))
	sb.WriteString(fmt.Sprintf("- **Detected Documents**: %d\n\n", report.TotalSourceFiles))

	sb.WriteString("| Source File | Entity Type | Target File | Est. Tokens | Action |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")

	for _, item := range report.Items {
		tokStr := fmt.Sprintf("~%d tokens", item.Tokens)
		if item.IsOversized {
			tokStr += " ⚠️ `[OVERSIZED]`"
		}
		sb.WriteString(fmt.Sprintf("| `%s` | **%s** | `%s` | %s | %s |\n",
			item.SourcePath, item.SourceType, item.TargetPath, tokStr, item.Action))
	}

	sb.WriteString("\n#### 📊 Context Window & Token Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Context Size**: ~%d tokens (%d characters)\n", report.TotalTokens, report.TotalCharacters))
	sb.WriteString(fmt.Sprintf("- **Always-On Tokens**: ~%d tokens\n", report.AlwaysOnTokens))
	sb.WriteString(fmt.Sprintf("- **On-Demand Tokens**: ~%d tokens\n", report.OnDemandTokens))

	sb.WriteString("\n> [!NOTE]\n")
	sb.WriteString("> Estimated token counts (~Tokens) are calculated locally without external API calls or cost, indicating baseline prompt context overhead per conversation turn.\n")

	return sb.String()
}

func formatSpecMatrixCLITable(matrix *SpecMatrixReport) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("📋 [token-hop Specification Matrix] %s  ──►  %s\n\n",
		strings.ToUpper(matrix.FromPlatform), strings.ToUpper(matrix.ToPlatform)))

	headers := []string{"Entity", fmt.Sprintf("%s Location", matrix.FromPlatform), fmt.Sprintf("%s Location", matrix.ToPlatform), "Target Behavior"}
	colWidths := []int{11, 30, 30, 36}

	topBorder := buildBorder("┌", "┬", "┐", "─", colWidths)
	midBorder := buildBorder("├", "┼", "┤", "─", colWidths)
	botBorder := buildBorder("└", "┴", "┘", "─", colWidths)

	sb.WriteString(topBorder)
	sb.WriteByte('\n')
	sb.WriteString(buildRow(headers, colWidths))
	sb.WriteByte('\n')
	sb.WriteString(midBorder)
	sb.WriteByte('\n')

	for _, item := range matrix.Items {
		row := []string{
			string(item.EntityType),
			truncateString(item.SourceLocation, colWidths[1]),
			truncateString(item.TargetLocation, colWidths[2]),
			truncateString(item.TargetBehavior, colWidths[3]),
		}
		sb.WriteString(buildRow(row, colWidths))
		sb.WriteByte('\n')
	}
	sb.WriteString(botBorder)
	sb.WriteByte('\n')

	return sb.String()
}

func formatSpecMatrixMarkdown(matrix *SpecMatrixReport) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("### 📋 Platform Specification Matrix: %s ──► %s\n\n",
		matrix.FromPlatform, matrix.ToPlatform))

	sb.WriteString(fmt.Sprintf("| Entity | %s Location & Syntax | %s Location & Syntax | Target Behavior |\n",
		matrix.FromPlatform, matrix.ToPlatform))
	sb.WriteString("| :--- | :--- | :--- | :--- |\n")

	for _, item := range matrix.Items {
		sb.WriteString(fmt.Sprintf("| **%s** | `%s`<br>*(%s)* | `%s`<br>*(%s)* | %s |\n",
			item.EntityType, item.SourceLocation, item.SourceSyntax, item.TargetLocation, item.TargetSyntax, item.TargetBehavior))
	}

	return sb.String()
}

func buildBorder(left, mid, right, fill string, colWidths []int) string {
	var parts []string
	for _, w := range colWidths {
		parts = append(parts, strings.Repeat(fill, w+2))
	}
	return left + strings.Join(parts, mid) + right
}

func buildRow(cols []string, colWidths []int) string {
	var parts []string
	for i, col := range cols {
		w := colWidths[i]
		parts = append(parts, fmt.Sprintf(" %-*s ", w, col))
	}
	return "│" + strings.Join(parts, "│") + "│"
}

func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	if maxLen <= 3 {
		return str[:maxLen]
	}
	return str[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
