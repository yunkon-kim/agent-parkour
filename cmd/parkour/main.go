package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yunkon-kim/agent-parkour/pkg/ai"
	"github.com/yunkon-kim/agent-parkour/pkg/audit"
	"github.com/yunkon-kim/agent-parkour/pkg/describer"
	"github.com/yunkon-kim/agent-parkour/pkg/engine"
	"github.com/yunkon-kim/agent-parkour/pkg/refiner"
)

var (
	version = "v0.0.3"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "parkour",
		Aliases: []string{"pk", "agent-parkour", "thop", "token-hop"},
		Short:   "agent-parkour (parkour / pk): Vault across AI coding agents without token walls",
		Long: `agent-parkour (parkour, pk) is a universal prompt compiler and context synchronizer
that eliminates configuration fragmentation across Google Antigravity, Cursor,
Claude Code, GitHub Copilot, Gemini CLI, and Roo Code.

Hit a token wall? Vault across your AI coding agents in under 10ms! 🏃💨`,
	}

	// 1. Version Command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("agent-parkour (parkour / pk) %s - Cross-Agent Prompt Compiler\n", version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// 2. Convert Command
	var fromFormat, toFormat, inputPath, outputPath string
	var useAI, autoDecompose, dryRun, generatePrompt bool
	var aiProvider, aiModel string

	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert instructions across AI agent formats with automatic path mapping",
		Example: `  # Convert from GitHub Copilot to Google Antigravity (auto-locates .github/ files)
  parkour convert --from copilot --to antigravity

  # Preview conversion without writing files (dry-run mode)
  parkour convert --from copilot --to antigravity --dry-run

  # Convert and automatically generate target AI refinement prompt (refine-prompt.md)
  parkour convert --from copilot --to antigravity --generate-prompt

  # Convert from GitHub Copilot to all supported targets (Antigravity, Cursor)
  parkour convert --from copilot --to all

  # Convert to a specific target (auto-detects source in repo)
  parkour convert --to cursor

  # Optional AI-powered semantic decomposition for oversized rules
  parkour convert --from copilot --to antigravity --decompose --ai --provider gemini`,
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)

			// Require --to flag to avoid unintended batch overwrite
			if !cmd.Flags().Changed("to") || toFormat == "" || toFormat == "auto" {
				fmt.Println("⚠️  Please specify a target format with --to <format> or --to all")
				fmt.Println()
				fmt.Println("👉 Common Examples:")
				fmt.Println("   parkour convert --from copilot --to antigravity")
				fmt.Println("   parkour convert --from copilot --to all")
				fmt.Println("   parkour convert --to cursor")
				fmt.Println()
				fmt.Println("💡 Run 'parkour convert --help' for all available options.")
				return nil
			}

			// 1. Auto-detect or resolve source path
			checkPath := "."
			if cmd.Flags().Changed("input") {
				checkPath = inputPath
			}
			detectedFrom, detectedInput := eng.AutoDetectSource(checkPath)

			if !cmd.Flags().Changed("from") || fromFormat == "auto" {
				fromFormat = detectedFrom
			}
			if !cmd.Flags().Changed("input") {
				// If user explicitly specified --from, resolve standard folder if present
				if fromFormat == "copilot" && checkPath == "." {
					if _, err := os.Stat(".github"); err == nil {
						inputPath = ".github"
					} else {
						inputPath = detectedInput
					}
				} else {
					inputPath = detectedInput
				}
			}

			// 2. Resolve target list
			targets := eng.AutoDetectTargets(fromFormat, toFormat)
			if len(targets) == 0 {
				return fmt.Errorf("no valid targets specified for %q", toFormat)
			}

			// 2.1 If Dry-Run is enabled, preview with Describe instead of writing files
			if dryRun {
				for _, target := range targets {
					report, err := eng.Describe(fromFormat, target, inputPath, outputPath, 400)
					if err != nil {
						return fmt.Errorf("dry-run preview failed for %s: %w", target, err)
					}
					fmt.Print(describer.FormatMappingReport(report, describer.FormatTable))
					if len(targets) > 1 {
						fmt.Println()
					}
				}
				return nil
			}

			// 3. Setup AI Provider if requested or configured
			aiCfg := ai.LoadConfig()
			if cmd.Flags().Changed("ai") {
				aiCfg.Enabled = useAI
			}
			if cmd.Flags().Changed("provider") {
				aiCfg.Provider = aiProvider
			}
			if cmd.Flags().Changed("model") {
				aiCfg.Model = aiModel
			}

			if aiCfg.Enabled || useAI || autoDecompose {
				provider, err := ai.NewProvider(aiCfg)
				if err != nil {
					fmt.Printf("⚠️  AI Provider note: %v (falling back to deterministic core)\n", err)
				} else {
					eng.SetAIProvider(provider)
					fmt.Printf("🤖 [agent-parkour AI] Augmented with %s (%s)\n", strings.ToUpper(provider.Name()), provider.Model())
				}
			}

			fmt.Printf("🚀 [agent-parkour] Vaulting from '%s' to '%s'...\n", strings.ToUpper(fromFormat), strings.Join(targets, ", "))
			fmt.Printf("   Source Path : %s\n", inputPath)
			fmt.Printf("   Output Dir  : %s\n", outputPath)
			if autoDecompose {
				fmt.Printf("   AI Mode     : Semantic JIT Decomposition (>400 tokens -> JIT Skills)\n")
			}
			fmt.Println()

			// 4. Perform conversion for each target format
			var totalFiles []string
			var lastAudit *audit.Report

			for _, target := range targets {
				fmt.Printf("🔨 Generating %s configuration...\n", strings.ToUpper(target))
				writtenFiles, auditReport, err := eng.ConvertWithAI(context.Background(), fromFormat, target, inputPath, outputPath, autoDecompose, 400)
				if err != nil {
					fmt.Printf("   ❌ Error generating %s: %v\n", target, err)
					continue
				}
				lastAudit = auditReport
				totalFiles = append(totalFiles, writtenFiles...)
				fmt.Printf("   ✅ Successfully generated %d files for %s\n", len(writtenFiles), target)

				// 4.1 Generate AI Refinement Prompt if requested
				if generatePrompt {
					refinePrompt, err := eng.GenerateRefinePrompt(fromFormat, target, inputPath, "", 400)
					if err == nil && refinePrompt != "" {
						promptFileName := filepath.Join(outputPath, "refine-prompt-"+target+".md")
						switch target {
						case "antigravity":
							promptFileName = filepath.Join(outputPath, ".agents", "refine-prompt.md")
						case "cursor":
							promptFileName = filepath.Join(outputPath, ".cursor", "refine-prompt.md")
						case "copilot":
							promptFileName = filepath.Join(outputPath, ".github", "refine-prompt.md")
						}
						_ = os.MkdirAll(filepath.Dir(promptFileName), 0755)
						_ = os.WriteFile(promptFileName, []byte(refinePrompt), 0644)
						fmt.Printf("   📝 Generated target AI refinement prompt -> %s\n", promptFileName)
					}
				}
			}

			// 5. Display Summary & Audit
			exceededCount := 0
			if lastAudit != nil {
				fmt.Printf("\n📊 Context Window & Token Summary:\n")
				fmt.Printf("   • Total Documents : %d\n", lastAudit.TotalDocuments)
				fmt.Printf("   • Total Tokens    : ~%d tokens\n", lastAudit.TotalTokens)
				fmt.Printf("   • Total Characters: %d chars\n", lastAudit.TotalCharacters)

				for _, item := range lastAudit.Items {
					if item.ExceedsLimit {
						exceededCount++
						fmt.Printf("   ⚠️  [%s] %s (~%d tokens) is larger than recommended size (>400 tokens)\n", item.Type, item.ID, item.Tokens)
						if item.Recommendation != "" {
							fmt.Printf("      └─ Tip: %s\n", item.Recommendation)
						}
					}
				}
				if exceededCount == 0 {
					fmt.Println("   ✨ All instructions are within safe context window limits!")
				}
			}

			fmt.Printf("\n🎉 Conversion completed! Total %d files generated.\n\n", len(totalFiles))

			// 6. Action Items & Next Steps for the User
			fmt.Println("👉 Next Steps & Action Items:")
			fmt.Println("   1. Review changes   : Run 'git status' or 'git diff' to review the generated files.")

			for _, target := range targets {
				switch target {
				case "antigravity":
					fmt.Println("   2. Google Antigravity: Open in Antigravity. Rules (.agents/rules/) and slash workflows (.agents/workflows/) are ready.")
				case "cursor":
					fmt.Println("   2. Cursor AI        : Open in Cursor. Rules in .cursor/rules/*.mdc are active.")
				case "copilot":
					fmt.Println("   2. GitHub Copilot   : Open in VSCode. Instructions in .github/ are active.")
				}
			}

			if exceededCount > 0 && !autoDecompose {
				fmt.Printf("   3. Optimize Context : %d oversized rule(s) exceed recommended limit (>400 tokens). Consider '--decompose --ai' for JIT skills.\n", exceededCount)
				fmt.Println("   4. Commit Changes   : git add . && git commit -m \"docs: sync AI guidelines to " + strings.Join(targets, ", ") + "\"")
			} else {
				fmt.Println("   3. Commit Changes   : git add . && git commit -m \"docs: sync AI guidelines to " + strings.Join(targets, ", ") + "\"")
			}
			return nil
		},
	}
	convertCmd.Flags().StringVarP(&fromFormat, "from", "f", "auto", "Source agent format (copilot, cursor, antigravity)")
	convertCmd.Flags().StringVarP(&toFormat, "to", "t", "", "Target agent format (antigravity, cursor, copilot, or 'all')")
	convertCmd.Flags().StringVarP(&inputPath, "input", "i", ".", "Path to source directory or file (auto-detected if omitted)")
	convertCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "Output directory for target files")
	convertCmd.Flags().BoolVar(&useAI, "ai", false, "Enable Generative AI augmentation")
	convertCmd.Flags().BoolVarP(&autoDecompose, "decompose", "d", false, "Semantically decompose oversized rules into JIT skills")
	convertCmd.Flags().StringVar(&aiProvider, "provider", "gemini", "AI provider (gemini, claude, openai, ollama, mock)")
	convertCmd.Flags().StringVar(&aiModel, "model", "", "AI model override")
	convertCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview mapping plan in table format without writing files")
	convertCmd.Flags().BoolVar(&generatePrompt, "generate-prompt", false, "Generate target-specific AI refinement prompt (refine-prompt.md) after conversion")
	convertCmd.Flags().BoolVar(&generatePrompt, "prompt", false, "Alias for --generate-prompt")
	rootCmd.AddCommand(convertCmd)

	// 3. Describe Command
	var descFrom, descTo, descInput, descOutput, descFormat, descOutFile string
	var descSpec bool
	var descMaxTokens int

	describeCmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe cross-agent mapping plan and specifications in tabular format before conversion",
		Example: `  # Preview file mapping plan from Copilot to Antigravity
  parkour describe --from copilot --to antigravity

  # Auto-detect source in current project and preview mapping to Cursor
  parkour describe --to cursor

  # Save mapping report to a file (auto-detects Markdown for .md or JSON for .json)
  parkour describe --from copilot --to antigravity --out mapping.md
  parkour describe --from copilot --to antigravity --out plan.json

  # Show conceptual specification matrix between platforms
  parkour describe --from copilot --to antigravity --spec

  # Output as Markdown or JSON to stdout
  parkour describe --from copilot --to antigravity --format markdown
  parkour describe --from copilot --to antigravity --format json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)

			// Helper to determine effective format
			resolveFormat := func(explicitFormat, outFile string) describer.OutputFormat {
				if cmd.Flags().Changed("format") {
					return describer.OutputFormat(explicitFormat)
				}
				if outFile != "" {
					lower := strings.ToLower(outFile)
					if strings.HasSuffix(lower, ".md") {
						return describer.FormatMarkdown
					}
					if strings.HasSuffix(lower, ".json") {
						return describer.FormatJSON
					}
				}
				return describer.FormatTable
			}

			// 1. Spec Matrix Mode (--spec)
			if descSpec {
				fromP := descFrom
				if fromP == "" || fromP == "auto" {
					fromP, _ = eng.AutoDetectSource(".")
				}
				toP := descTo
				if toP == "" || toP == "auto" {
					targets := eng.AutoDetectTargets(fromP, "")
					if len(targets) > 0 {
						toP = targets[0]
					} else {
						toP = "antigravity"
					}
				}
				matrix := describer.BuildSpecMatrix(fromP, toP)
				effectiveFormat := resolveFormat(descFormat, descOutFile)
				content := describer.FormatSpecMatrix(matrix, effectiveFormat)

				if descOutFile != "" {
					if err := os.WriteFile(descOutFile, []byte(content), 0644); err != nil {
						return fmt.Errorf("failed to write spec matrix to %s: %w", descOutFile, err)
					}
					fmt.Printf("📄 Specification matrix successfully exported to: %s (%s format)\n", descOutFile, effectiveFormat)
					return nil
				}

				fmt.Print(content)
				return nil
			}

			// 2. Resolve source path & format
			checkPath := "."
			if cmd.Flags().Changed("input") {
				checkPath = descInput
			}
			detectedFrom, detectedInput := eng.AutoDetectSource(checkPath)

			if !cmd.Flags().Changed("from") || descFrom == "auto" {
				descFrom = detectedFrom
			}
			if !cmd.Flags().Changed("input") {
				if descFrom == "copilot" && checkPath == "." {
					if _, err := os.Stat(".github"); err == nil {
						descInput = ".github"
					} else {
						descInput = detectedInput
					}
				} else {
					descInput = detectedInput
				}
			}

			// 3. Resolve target list
			targets := eng.AutoDetectTargets(descFrom, descTo)
			if len(targets) == 0 {
				return fmt.Errorf("no valid target platform specified (use --to <target>)")
			}

			effectiveFormat := resolveFormat(descFormat, descOutFile)
			var fullOutput strings.Builder

			for _, target := range targets {
				report, err := eng.Describe(descFrom, target, descInput, descOutput, descMaxTokens)
				if err != nil {
					return fmt.Errorf("describe failed for target %s: %w", target, err)
				}
				fullOutput.WriteString(describer.FormatMappingReport(report, effectiveFormat))
				if len(targets) > 1 {
					fullOutput.WriteString("\n")
				}
			}

			if descOutFile != "" {
				if err := os.WriteFile(descOutFile, []byte(fullOutput.String()), 0644); err != nil {
					return fmt.Errorf("failed to write description report to %s: %w", descOutFile, err)
				}
				fmt.Printf("📄 Mapping description plan successfully exported to: %s (%s format)\n", descOutFile, effectiveFormat)
				return nil
			}

			fmt.Print(fullOutput.String())
			return nil
		},
	}
	describeCmd.Flags().StringVarP(&descFrom, "from", "f", "auto", "Source agent format (copilot, cursor, antigravity, claude)")
	describeCmd.Flags().StringVarP(&descTo, "to", "t", "auto", "Target agent format (antigravity, cursor, copilot, claude, or 'all')")
	describeCmd.Flags().StringVarP(&descInput, "input", "i", ".", "Path to source directory or file (auto-detected if omitted)")
	describeCmd.Flags().StringVarP(&descOutput, "output", "o", ".", "Target output directory")
	describeCmd.Flags().StringVar(&descFormat, "format", "table", "Output format: table, markdown, json")
	describeCmd.Flags().StringVar(&descOutFile, "out", "", "Export mapping plan directly to a file (.md, .json, .txt)")
	describeCmd.Flags().BoolVar(&descSpec, "spec", false, "Display standard specification matrix between platforms")
	describeCmd.Flags().IntVarP(&descMaxTokens, "max-tokens", "m", 400, "Maximum recommended token limit per rule")
	rootCmd.AddCommand(describeCmd)

	// 4. Prompt Command (Target-Specific AI Refinement Prompt Generator)
	var promptTo, promptFrom, promptFile, promptDir, promptOut, promptGuidance string
	var promptMaxTokens int

	promptCmd := &cobra.Command{
		Use:   "prompt",
		Short: "Generate target-specific AI refinement prompts for semantic 2nd-stage optimization",
		Example: `  # Generate refinement prompt for a converted Antigravity workflow
  parkour prompt --to antigravity --file .agents/workflows/git-commit.md

  # Generate refinement prompt for all files in directory and export to file
  parkour prompt --to antigravity --dir .agents/ --out .agents/refine-prompt.md
  parkour prompt --to cursor --dir .cursor/rules/ --out .cursor/refine-prompt.md

  # Output target AI optimization guideline template to console
  parkour prompt --to claude
  parkour prompt --to antigravity`,
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)

			// 1. Resolve target platform
			toPlatform := promptTo
			if toPlatform == "" || toPlatform == "auto" {
				srcFormat, _ := eng.AutoDetectSource(".")
				targets := eng.AutoDetectTargets(srcFormat, "")
				if len(targets) > 0 {
					toPlatform = targets[0]
				} else {
					toPlatform = "antigravity"
				}
			}

			gen := refiner.NewGenerator(promptMaxTokens)
			var promptContent string
			var err error

			if promptFile != "" {
				promptContent, err = gen.GenerateFromFile(toPlatform, promptFile, promptGuidance)
				if err != nil {
					return fmt.Errorf("failed to generate prompt from file %s: %w", promptFile, err)
				}
			} else if promptDir != "" {
				promptContent, err = gen.GenerateFromDirectory(toPlatform, promptDir, promptGuidance)
				if err != nil {
					return fmt.Errorf("failed to generate prompt from directory %s: %w", promptDir, err)
				}
			} else {
				srcFormat := promptFrom
				checkPath := "."
				if srcFormat == "" || srcFormat == "auto" {
					srcFormat, checkPath = eng.AutoDetectSource(checkPath)
				}
				docs, parseErr := eng.ParseSource(srcFormat, checkPath)
				if parseErr == nil && len(docs) > 0 {
					promptContent, err = gen.GenerateFromDocs(toPlatform, docs, promptGuidance)
				} else {
					promptContent, err = gen.GenerateFromDocs(toPlatform, nil, promptGuidance)
				}
				if err != nil {
					return fmt.Errorf("failed to generate prompt: %w", err)
				}
			}

			if promptOut != "" {
				_ = os.MkdirAll(filepath.Dir(promptOut), 0755)
				if err := os.WriteFile(promptOut, []byte(promptContent), 0644); err != nil {
					return fmt.Errorf("failed to write prompt to %s: %w", promptOut, err)
				}
				fmt.Printf("📝 Refinement prompt successfully saved to: %s\n", promptOut)
				return nil
			}

			fmt.Print(promptContent)
			return nil
		},
	}
	promptCmd.Flags().StringVarP(&promptTo, "to", "t", "antigravity", "Target AI platform (antigravity, cursor, claude, copilot)")
	promptCmd.Flags().StringVarP(&promptFrom, "from", "f", "auto", "Source format if auto-parsing repo")
	promptCmd.Flags().StringVar(&promptFile, "file", "", "Target file to generate refinement prompt for")
	promptCmd.Flags().StringVar(&promptDir, "dir", "", "Directory containing files to include in prompt")
	promptCmd.Flags().StringVarP(&promptOut, "out", "o", "", "File path to save the generated prompt")
	promptCmd.Flags().StringVarP(&promptGuidance, "guidance", "g", "", "Custom project-specific instructions to append to prompt")
	promptCmd.Flags().IntVarP(&promptMaxTokens, "max-tokens", "m", 400, "Recommended token limit threshold")
	rootCmd.AddCommand(promptCmd)

	// 5. Audit Command
	var auditDir string
	var maxTokens int
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit token counts and character limits of instructions",
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)
			from, input := eng.AutoDetectSource(auditDir)
			docs, err := eng.ParseSource(from, input)
			if err != nil {
				return fmt.Errorf("failed to parse instructions in %s: %w", auditDir, err)
			}

			report := audit.AuditDocuments(docs, maxTokens)
			fmt.Printf("📊 [agent-parkour Token Audit Report for: %s (%s)]\n", input, strings.ToUpper(from))
			fmt.Printf("   Total Documents: %d | Total Tokens: ~%d | Total Characters: %d\n\n",
				report.TotalDocuments, report.TotalTokens, report.TotalCharacters)

			fmt.Printf("%-30s | %-12s | %-10s | %-10s | %s\n", "ID", "Type", "Chars", "Tokens", "Status")
			fmt.Printf("%s\n", strings.Repeat("-", 80))
			for _, item := range report.Items {
				status := "OK"
				if item.ExceedsLimit {
					status = "OVERSIZED"
				}
				fmt.Printf("%-30s | %-12s | %-10d | ~%-9d | %s\n",
					item.ID, item.Type, item.Characters, item.Tokens, status)
			}
			return nil
		},
	}
	auditCmd.Flags().StringVarP(&auditDir, "input", "i", ".", "Directory to audit (auto-detects if '.')")
	auditCmd.Flags().IntVarP(&maxTokens, "max-tokens", "m", 400, "Maximum allowed tokens per instruction")
	rootCmd.AddCommand(auditCmd)

	// 6. Init Command
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize agent-parkour SSOT scaffolding and parkour.yaml in current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🏃 Initializing agent-parkour SSOT in current project...")
			if _, err := os.Stat("AGENTS.md"); os.IsNotExist(err) {
				content := `# Project Agent Guidelines (Single Source of Truth)

## Project Overview
- Purpose: (Define project objectives and core technologies)
- Languages & Frameworks: (e.g. Go 1.21+, React, Python)
- Build & Test: ` + "`go build ./... && go test ./...`" + `

## Core Conventions
- Follow standard formatting and linting rules.
- Write concise, explicit, and self-documenting code.
`
				_ = os.WriteFile("AGENTS.md", []byte(content), 0644)
				fmt.Println("   • Created AGENTS.md (Root SSOT)")
			}

			if _, err := os.Stat("parkour.yaml"); os.IsNotExist(err) {
				cfgContent := `version: "1.0"
ssot: "AGENTS.md"

# Target AI configurations
targets:
  - name: "antigravity"
    output_dir: ".agents"
    enable_skills: true
  - name: "cursor"
    output_dir: ".cursor/rules"
  - name: "claude"
    output_file: "CLAUDE.md"
    skills_dir: ".claude/skills"
  - name: "copilot"
    instructions_dir: ".github/instructions"
    prompts_dir: ".github/prompts"
  - name: "gemini"
    output_file: "GEMINI.md"

# Context limits & Optimization
context_limits:
  max_tokens_per_rule: 400
  max_characters_per_rule: 1800
  auto_decompose: true
`
				_ = os.WriteFile("parkour.yaml", []byte(cfgContent), 0644)
				fmt.Println("   • Created parkour.yaml (Project configuration)")
			}

			fmt.Println("✅ agent-parkour initialized successfully! Start vaulting across agents.")
			return nil
		},
	}
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
