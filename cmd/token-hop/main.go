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
that resolves configuration fragmentation across Google Antigravity, Cursor,
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
		Short: "Directly convert instructions from one agent format to another",
		Example: `  token-hop convert --from copilot --to antigravity --input .github --output .
  thop convert --from copilot --to cursor --input .github --output .
  thop convert --from copilot --to antigravity --decompose --ai --provider gemini`,
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)

			// Setup AI Provider if requested or configured
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
					fmt.Printf("⚠️  AI Provider initialization note: %v (falling back to deterministic core)\n", err)
				} else {
					eng.SetAIProvider(provider)
					fmt.Printf("🤖 [token-hop AI] Augmented with %s (%s)\n", strings.ToUpper(provider.Name()), provider.Model())
				}
			}

			fmt.Printf("🚀 [token-hop] Converting from '%s' to '%s'...\n", fromFormat, toFormat)
			fmt.Printf("   Source: %s\n", inputPath)
			fmt.Printf("   Target: %s\n", outputPath)
			if autoDecompose {
				fmt.Printf("   Mode  : Semantic JIT Decomposition (Tokens > 400 -> JIT Skills)\n")
			}
			fmt.Println()

			writtenFiles, auditReport, err := eng.ConvertWithAI(context.Background(), fromFormat, toFormat, inputPath, outputPath, autoDecompose, 400)
			if err != nil {
				return fmt.Errorf("conversion failed: %w", err)
			}

			fmt.Printf("✅ Successfully compiled %d files to %s:\n", len(writtenFiles), toFormat)
			for _, file := range writtenFiles {
				fmt.Printf("   • %s\n", file)
			}

			fmt.Printf("\n📊 Context Budget Summary:\n")
			fmt.Printf("   • Total Documents : %d\n", auditReport.TotalDocuments)
			fmt.Printf("   • Total Tokens    : ~%d tokens\n", auditReport.TotalTokens)
			fmt.Printf("   • Total Characters: %d chars\n\n", auditReport.TotalCharacters)

			exceededCount := 0
			for _, item := range auditReport.Items {
				if item.ExceedsBudget {
					exceededCount++
					fmt.Printf("   ⚠️  [%s] %s (~%d tokens) exceeds budget limit!\n", item.Type, item.ID, item.Tokens)
					if item.Recommendation != "" {
						fmt.Printf("      └─ Recommendation: %s\n", item.Recommendation)
					}
				}
			}

			if exceededCount == 0 {
				fmt.Println("   ✨ All instructions are within safe context window limits!")
			}

			return nil
		},
	}
	convertCmd.Flags().StringVarP(&fromFormat, "from", "f", "copilot", "Source agent format (copilot, cursor, antigravity)")
	convertCmd.Flags().StringVarP(&toFormat, "to", "t", "antigravity", "Target agent format (antigravity, cursor, copilot)")
	convertCmd.Flags().StringVarP(&inputPath, "input", "i", ".github", "Path to source directory or file")
	convertCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "Output directory for target files")
	convertCmd.Flags().BoolVar(&useAI, "ai", false, "Enable Generative AI augmentation")
	convertCmd.Flags().BoolVarP(&autoDecompose, "decompose", "d", false, "Semantically decompose oversized rules into JIT skills")
	convertCmd.Flags().StringVar(&aiProvider, "provider", "gemini", "AI provider (gemini, claude, openai, ollama, mock)")
	convertCmd.Flags().StringVar(&aiModel, "model", "", "AI model override")
	rootCmd.AddCommand(convertCmd)

	// 3. Compile Command
	var targetList string
	var configPath string
	compileCmd := &cobra.Command{
		Use:   "compile",
		Short: "Compile instructions from SSOT (AGENTS.md) to multiple targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := engine.LoadConfig(configPath)
			if err != nil {
				return err
			}

			eng := engine.NewEngine(cfg)
			docs, err := eng.ParseSource("antigravity", ".")
			if err != nil || len(docs) == 0 {
				docs, err = eng.ParseSource("copilot", ".github")
			}
			if err != nil || len(docs) == 0 {
				return fmt.Errorf("no source instructions found to compile")
			}

			targets := strings.Split(targetList, ",")
			for _, t := range targets {
				t = strings.TrimSpace(t)
				if t == "" {
					continue
				}
				fmt.Printf("🔨 Compiling target: %s...\n", t)
				written, err := eng.EmitTarget(t, docs, ".")
				if err != nil {
					fmt.Printf("   ❌ Error emitting %s: %v\n", t, err)
				} else {
					fmt.Printf("   ✅ Emitted %d files for %s\n", len(written), t)
				}
			}
			return nil
		},
	}
	compileCmd.Flags().StringVarP(&targetList, "target", "t", "antigravity,cursor", "Comma-separated target list")
	compileCmd.Flags().StringVarP(&configPath, "config", "c", "token-hop.yaml", "Path to token-hop.yaml")
	rootCmd.AddCommand(compileCmd)

	// 4. Audit Command
	var auditDir string
	var maxTokens int
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit token budget and character counts of instructions",
		RunE: func(cmd *cobra.Command, args []string) error {
			eng := engine.NewEngine(nil)
			docs, err := eng.ParseSource("copilot", auditDir)
			if err != nil {
				docs, err = eng.ParseSource("antigravity", auditDir)
			}
			if err != nil {
				return fmt.Errorf("failed to parse instructions in %s: %w", auditDir, err)
			}

			report := budget.AuditDocuments(docs, maxTokens)
			fmt.Printf("📊 [token-hop Audit Report for: %s]\n", auditDir)
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
	auditCmd.Flags().StringVarP(&auditDir, "input", "i", ".github", "Directory to audit")
	auditCmd.Flags().IntVarP(&maxTokens, "max-tokens", "m", 400, "Maximum allowed tokens per instruction")
	rootCmd.AddCommand(auditCmd)

	// 5. Init Command
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
