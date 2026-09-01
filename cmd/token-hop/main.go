package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yunkon-kim/token-hop/pkg/ai"
	"github.com/yunkon-kim/token-hop/pkg/budget"
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

	// 2. Convert Command (Smart Auto-Detection & Multi-Target)
	var fromFormat, toFormat, inputPath, outputPath string
	var useAI, autoDecompose bool
	var aiProvider, aiModel string

	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Automatically detect and convert instructions across AI agent formats",
		Example: `  # Zero-config: Automatically detects source in current repo and converts to targets
  thop convert

  # Convert to a specific target
  thop convert --to antigravity
  thop convert --to cursor

  # Explicit source and target
  thop convert --from copilot --to antigravity --input .github --output .

  # Optional AI-powered semantic decomposition
  thop convert --decompose --ai --provider gemini`,
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)

			// 1. Auto-detect source if not explicitly provided
			checkPath := "."
			if cmd.Flags().Changed("input") {
				checkPath = inputPath
			}
			detectedFrom, detectedInput := eng.AutoDetectSource(checkPath)
			if !cmd.Flags().Changed("from") || fromFormat == "auto" {
				fromFormat = detectedFrom
			}
			if !cmd.Flags().Changed("input") {
				inputPath = detectedInput
			}

			// 2. Auto-detect targets if not explicitly provided
			targets := eng.AutoDetectTargets(fromFormat, toFormat)

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

			fmt.Printf("🔍 Source detected: [%s] at '%s'\n", strings.ToUpper(fromFormat), inputPath)
			fmt.Printf("🎯 Target format(s): %s\n", strings.Join(targets, ", "))
			if autoDecompose {
				fmt.Printf("⚡ Mode: Semantic JIT Decomposition (>400 tokens -> JIT Skills)\n")
			}
			fmt.Println()

			// 4. Perform conversion for each target format
			var totalFiles []string
			var lastAudit *budget.AuditReport

			for _, target := range targets {
				fmt.Printf("🔨 Converting to %s...\n", target)
				writtenFiles, auditReport, err := eng.ConvertWithAI(context.Background(), fromFormat, target, inputPath, outputPath, autoDecompose, 400)
				if err != nil {
					fmt.Printf("   ❌ Error converting to %s: %v\n", target, err)
					continue
				}
				lastAudit = auditReport
				totalFiles = append(totalFiles, writtenFiles...)
				fmt.Printf("   ✅ Generated %d files for %s\n", len(writtenFiles), target)
			}

			// 5. Display Summary & Audit
			if lastAudit != nil {
				fmt.Printf("\n📊 Context Budget Summary:\n")
				fmt.Printf("   • Total Documents : %d\n", lastAudit.TotalDocuments)
				fmt.Printf("   • Total Tokens    : ~%d tokens\n", lastAudit.TotalTokens)
				fmt.Printf("   • Total Characters: %d chars\n", lastAudit.TotalCharacters)

				exceededCount := 0
				for _, item := range lastAudit.Items {
					if item.ExceedsBudget {
						exceededCount++
						fmt.Printf("   ⚠️  [%s] %s (~%d tokens) exceeds recommended budget!\n", item.Type, item.ID, item.Tokens)
						if item.Recommendation != "" {
							fmt.Printf("      └─ Recommendation: %s\n", item.Recommendation)
						}
					}
				}
				if exceededCount == 0 {
					fmt.Println("   ✨ All instructions are within safe context window limits!")
				}
			}

			fmt.Printf("\n🎉 All operations completed! Total %d files generated.\n", len(totalFiles))
			return nil
		},
	}
	convertCmd.Flags().StringVarP(&fromFormat, "from", "f", "auto", "Source agent format (auto, copilot, cursor, antigravity)")
	convertCmd.Flags().StringVarP(&toFormat, "to", "t", "auto", "Target agent format(s) (auto, all, antigravity, cursor, copilot)")
	convertCmd.Flags().StringVarP(&inputPath, "input", "i", ".", "Path to source directory or file")
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
		Short: "Audit token budget and character counts of instructions",
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)
			from, input := eng.AutoDetectSource(auditDir)
			docs, err := eng.ParseSource(from, input)
			if err != nil {
				return fmt.Errorf("failed to parse instructions in %s: %w", auditDir, err)
			}

			report := budget.AuditDocuments(docs, maxTokens)
			fmt.Printf("📊 [token-hop Audit Report for: %s (%s)]\n", input, strings.ToUpper(from))
			fmt.Printf("   Total Documents: %d | Total Tokens: ~%d | Total Characters: %d\n\n",
				report.TotalDocuments, report.TotalTokens, report.TotalCharacters)

			fmt.Printf("%-30s | %-12s | %-10s | %-10s | %s\n", "ID", "Type", "Chars", "Tokens", "Status")
			fmt.Printf("%s\n", strings.Repeat("-", 80))
			for _, item := range report.Items {
				status := "OK"
				if item.ExceedsBudget {
					status = "EXCEEDED"
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
