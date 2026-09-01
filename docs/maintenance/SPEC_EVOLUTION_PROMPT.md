# Agent Specification Evolution & Target Update Prompt
**AI 코딩 에이전트 스펙 변경 및 신규 타겟 추가를 위한 AI 프롬프트 가이드**

> 이 문서는 Cursor, Google Antigravity, GitHub Copilot, Claude Code 등 AI 도구들의 **설정 디렉토리, 파일 확장자, 프론트매터 스키마, 토큰 제한 정책이 변경되거나 신규 AI 도구가 출시되었을 때**, AI 에이전트(Antigravity, Claude, ChatGPT 등)에게 전달하여 `token-hop` 코드를 안전하고 일관되게 최신화하도록 지시하는 **스펙 진화 프롬프트(Spec Evolution Prompt)**입니다.

---

## 📋 복사하여 사용할 수 있는 프롬프트 템플릿

```markdown
# Role: Senior Compiler & Toolchain Engineer for Token-Hop

## Goal
Update the `token-hop` codebase to support the latest specifications or new AI coding assistant environments.

## Target AI Tool / Specification Change
- **Tool Name**: [예: Cursor / Google Antigravity / Windsurf / Claude Code / Roo Code]
- **Target Version / Year**: [예: 2026 / Latest]
- **Changed / New Specifications**:
  - Config Path: [예: .cursor/rules/*.mdc -> .cursor/rules/*.md]
  - File Extension: [예: .mdc / .instructions.md / .prompt.md / SKILL.md]
  - Frontmatter Schema: [예: alwaysApply: bool, globs: [], description: string]
  - Token / Character Limit: [예: 6,000 characters per file, 12,000 total]
  - Activation Mode: [예: AlwaysOn | Glob | ModelDecision | OnDemand]

---

## Systematic Execution Steps (Checklist)

Please perform the update across the codebase following these exact 6 steps:

### 1. Update Universal Agent IR (UA-IR) AST
- Check `pkg/ir/types.go`.
- If new entity types (e.g., Subagent, Hook), activation modes, or binding fields are introduced, update the `EntityType`, `ActivationMode`, or `UAMetadata` structs while maintaining backwards compatibility.

### 2. Update / Create Parser (`pkg/parser/`)
- Modify or create `pkg/parser/<tool_name>.go`.
- Ensure it parses the new directory structure, splits YAML frontmatter safely, and normalizes fields into standard `ir.UADocument` objects.
- Handle fallback logic for previous versions (e.g. legacy `.agent/rules` vs `.agents/rules`).

### 3. Update / Create Target Emitter (`pkg/emitter/`)
- Modify or create `pkg/emitter/<tool_name>.go`.
- Implement `Emit(docs []*ir.UADocument) ([]string, error)` to generate clean Markdown files with the exact target frontmatter and formatting.
- Ensure `SanitizeFileName()` and `RenderMarkdownWithTitle()` are used to avoid duplicate headings.

### 4. Register in Engine & CLI (`pkg/engine/` & `cmd/token-hop/`)
- Update `pkg/engine/engine.go`: Register the tool name in `ParseSource()` and `EmitTarget()`.
- Update `cmd/token-hop/main.go`: Add the tool name to CLI flag descriptions (`--from`, `--to`, `--target`).

### 5. Add / Update Test Fixture & Verification (`test/`)
- Create or update test fixtures in `test/fixtures/<tool_name>/`.
- Add an automated test function in `test/` (e.g. `Test<ToolName>Conversion`).
- Run `go test -v ./...` and ensure all tests pass in <0.01s.

### 6. Synchronize Documentation
- Update `docs/design/token-hop-design-plan.md` (Mapping Matrix & Target Specs).
- Update `README.md` & `README_KO.md` (Supported Agents & Mapping Table).
- Run `make build && make test` to ensure zero compilation errors.

---

## Official Documentation & Reference Sources
When updating or adding support, cross-reference the official specifications below:
- **Google Antigravity**: https://cloud.google.com/gemini/docs/codeassist
- **Cursor AI Rules**: https://docs.cursor.com/context/rules-for-ai
- **GitHub Copilot Instructions**: https://docs.github.com/en/copilot/customizing-copilot/adding-custom-instructions-for-github-copilot
- **Claude Code**: https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview
- **Windsurf Rules**: https://docs.codeium.com/windsurf/memories
- **Roo Code / Cline**: https://docs.roocode.com
```

---

## 🛠️ 주요 도구별 최근 변경 이력 및 추적 포인트 (Reference Tracker)

| 도구 (Tool) | 현재 표준 경로 | 감시해야 할 변경 포인트 |
| :--- | :--- | :--- |
| **Google Antigravity** | `.agents/rules/*.md`<br>`.agents/workflows/*.md`<br>`.agents/skills/*/SKILL.md` | • JIT Skill 메타데이터 YAML 필드 확장<br>• Workflows 슬래시 커맨드 파라미터 규격 |
| **Cursor AI** | `.cursor/rules/*.mdc` | • `.mdc` 확장자의 순수 `.md` 변경 가능성<br>• `alwaysApply`, `globs`, `description` 외 신규 필드<br>• 단일 파일 6,000자 / 합계 12,000자 용량 정책 변동 |
| **GitHub Copilot** | `.github/copilot-instructions.md`<br>`.github/instructions/*.md`<br>`.github/prompts/*.md` | • `applyTo` 외 다중 Glob 배열 지원 여부<br>• 커스텀 에이전트/모드 스키마 (`.github/agents/`) 표준화 |
| **Claude Code** | `./CLAUDE.md`<br>`.claude/skills/*/SKILL.md` | • 디렉터리 캐스케이딩 탐색 우선순위<br>• 서브에이전트 정의 및 MCP 바인딩 구조 |
| **Windsurf (Codeium)**| `.windsurfrules`<br>`.codeium/rules/` | • 윈드서프 글로벌/워크스페이스 룰 규격<br>• Cascade 에이전트 지침 스키마 |
| **Roo Code / Cline** | `.roomodes`<br>`.clinerules` | • JSON/YAML 기반 Custom Modes 스키마 변경<br>• 도구 권한(Allowed Tools) 바인딩 형식 |

---

## 🌐 공식 레퍼런스 및 스펙 문서 링크 (Official Documentation Links)

스펙 업데이트 및 파서/생성기 변경 시 항상 아래 공식 문서의 최신 버전을 대조하여 작업하십시오.

1. **Google Antigravity / Gemini Code Assist**:
   - [Google Cloud Gemini Code Assist Documentation](https://cloud.google.com/gemini/docs/codeassist)
   - Antigravity IDE Custom Agents, Rules & Workflows Specification

2. **Cursor AI**:
   - [Cursor Rules Official Documentation](https://docs.cursor.com/context/rules-for-ai)
   - [Cursor System Overview & Context](https://docs.cursor.com/)

3. **GitHub Copilot**:
   - [GitHub Copilot Custom Instructions Guide](https://docs.github.com/en/copilot/customizing-copilot/adding-custom-instructions-for-github-copilot)
   - [GitHub Copilot Prompt Files (.github/prompts)](https://docs.github.com/en/copilot)

4. **Claude Code (Anthropic)**:
   - [Claude Code Overview & Tooling](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview)
   - [Claude Code Memory & CLAUDE.md Guide](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/memory)

5. **Windsurf (Codeium)**:
   - [Windsurf Rules & AI Context Guide](https://docs.codeium.com/windsurf/memories)

6. **Roo Code / Cline**:
   - [Roo Code Official Documentation](https://docs.roocode.com)
   - [Roo Code Custom Modes (.roomodes) Guide](https://github.com/RooVetGit/Roo-Code)

