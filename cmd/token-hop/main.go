package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yunkon-kim/token-hop/pkg/ai"
	"github.com/yunkon-kim/token-hop/pkg/audit"
	"github.com/yunkon-kim/token-hop/pkg/engine"
)

var (
	version = "v0.2.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "token-hop",
		Short: "token-hop (thop): Cross-Agent Prompt Compiler & Context Synchronizer",
		Long: `token-hop (thop) is a universal prompt compiler and context synchronizer
that eliminates configuration fragmentation across Google Antigravity, Cursor,
Claude Code, GitHub Copilot, Gemini CLI, and Roo Code.`,
	}

	// 1. Version Command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("token-hop (thop) %s - Cross-Agent Prompt Compiler\n", version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// 2. Convert Command
	var fromFormat, toFormat, inputPath, outputPath string
	var useAI, autoDecompose bool
	var aiProvider, aiModel string

	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert instructions across AI agent formats with automatic path mapping",
		Example: `  # Convert from GitHub Copilot to Google Antigravity (auto-locates .github/ files)
  thop convert --from copilot --to antigravity

  # Convert from GitHub Copilot to all supported targets (Antigravity, Cursor)
  thop convert --from copilot --to all

  # Convert to a specific target (auto-detects source in repo)
  thop convert --to cursor

  # Optional AI-powered semantic decomposition for oversized rules
  thop convert --from copilot --to antigravity --decompose --ai --provider gemini`,
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)

			// Require --to flag to avoid unintended batch overwrite
			if !cmd.Flags().Changed("to") || toFormat == "" || toFormat == "auto" {
				fmt.Println("⚠️  Please specify a target format with --to <format> or --to all")
				fmt.Println()
				fmt.Println("👉 Common Examples:")
				fmt.Println("   thop convert --from copilot --to antigravity")
				fmt.Println("   thop convert --from copilot --to all")
				fmt.Println("   thop convert --to cursor")
				fmt.Println()
				fmt.Println("💡 Run 'thop convert --help' for all available options.")
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
					fmt.Printf("🤖 [token-hop AI] Augmented with %s (%s)\n", strings.ToUpper(provider.Name()), provider.Model())
				}
			}

			fmt.Printf("🚀 [token-hop] Converting from '%s' to '%s'...\n", strings.ToUpper(fromFormat), strings.Join(targets, ", "))
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
	rootCmd.AddCommand(convertCmd)

	// 3. Audit Command
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
			fmt.Printf("📊 [token-hop Token Audit Report for: %s (%s)]\n", input, strings.ToUpper(from))
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

	// 4. Init Command
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize token-hop SSOT scaffolding and token-hop.yaml in current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🎉 Initializing token-hop SSOT in current project...")
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
				os.WriteFile("AGENTS.md", []byte(content), 0644)
				fmt.Println("   • Created AGENTS.md (Root SSOT)")
			}
			fmt.Println("✅ token-hop initialized successfully!")
			return nil
		},
	}
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
