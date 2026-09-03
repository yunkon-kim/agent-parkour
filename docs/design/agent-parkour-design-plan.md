# Agent-Parkour Engine: Cross-Agent Prompt Compiler Architecture & Plan
**AI 개발 환경 간 설정 이질성 해소를 위한 크로스 에이전트 프롬프트 컴파일러 설계 및 오픈소스 기획**

- **프로젝트명**: `agent-parkour` (구 `token-hop`)
- **CLI 명칭**: `parkour` (약칭: `pk`)
- **버전**: v0.1.0-draft
- **작성일**: 2026-09-01
- **저자**: Yunkon Kim

---

## 1. 배경 및 문제 정의 (Background & Problem Definition)

### 1.1 멀티 에이전트 IDE 환경의 파편화 (Context Fragmentation)
최근 소프트웨어 개발 생태계에서는 **Claude Code, Cursor, GitHub Copilot, Google Antigravity, Gemini CLI, Roo Code(Cline)** 등 다양한 AI 에이전트 도구가 병행 활용되고 있습니다. 개발자는 각 AI 도구의 토큰 제한 및 사용량 한도(Token & Usage Limits), 작업 특성(UI 디자인, 백엔드 아키텍처, CLI 자동화 등)에 따라 도구를 수시로 전환하며 작업합니다.

그러나 각 환경은 페르소나, 코딩 규칙(Rule), 워크플로우(Workflow), 스킬(Skill), 지시문(Instruction), 프롬프트(Prompt)를 정의하는 방식에서 극심한 파편화를 보입니다.

```
+-------------------------------------------------------------------------------+
|                       Multi-Agent Ecosystem Fragmentation                     |
+-------------------------------------------------------------------------------+
| [Google Antigravity] | .agents/rules/*.md, .agents/workflows/, .agents/skills/ |
| [Cursor AI]          | .cursor/rules/*.mdc (globs, alwaysApply, description)  |
| [Claude Code]        | ./CLAUDE.md, ~/.claude/CLAUDE.md, .claude/skills/      |
| [GitHub Copilot]     | .github/copilot-instructions.md, .github/prompts/      |
| [Gemini CLI]         | ./GEMINI.md, ~/.gemini/GEMINI.md                       |
| [Roo Code / Cline]   | .roomodes, .clinerules, Custom Modes                   |
+-------------------------------------------------------------------------------+
```

### 1.2 핵심 문제점
1. **맥락 유실 및 드리프트 (Context Rot & Drift)**:
   특정 도구(예: Cursor의 `.mdc`)에서 최적화된 고급 지침이 타 도구(예: Antigravity, Claude Code)로 전환될 때 반영되지 않거나 수동 복사 과정에서 누락·구문 오류가 발생합니다.
2. **컨텍스트 윈도우 오버플로우 (Context Window Overflow)**:
   규칙 파일 길이가 늘어나면(500줄 초과 또는 파일당 6,000~12,000자 초과 시) LLM의 지침 준수율이 급격히 저하되고 판단 혼란이 발생합니다.
3. **양방향 동기화 부재**:
   한 IDE에서 규칙을 수정했을 때 다른 IDE 규칙 및 원본 SSOT로 역반영되지 않아 프로젝트 규칙이 분기(Divergence)됩니다.

### 1.3 핵심 용어 및 약어 해설 (Glossary & Key Concepts)
본 설계서에서 사용되는 주요 핵심 용어와 기술 약어의 정의는 다음과 같습니다.

- **SSOT (Single Source of Truth, 단일 진실 원천)**:
  여러 AI 도구(Cursor, Claude, Copilot 등)마다 각기 다른 설정 파일들이 파편화되어 존재할 때, 이들의 근본이 되는 **단 하나의 마스터 기준 문서**(`AGENTS.md` 또는 프로젝트 메타데이터)를 의미합니다. `token-hop`은 SSOT를 수정하면 나머지 플랫폼 파일들이 자동으로 파생 생성되도록 보장합니다.
- **IR (Intermediate Representation, 중간 표현식)**:
  소스 포맷(예: Cursor의 `.mdc`)을 타겟 포맷(예: Antigravity의 `.md`)으로 직접 1:1 변환하지 않고, 컴파일러 내부에서 플랫폼 중립적으로 다루기 위해 표준화한 **추상 데이터 중간 규격**(Universal Agent IR, UA-IR)입니다. (컴파일러의 바이트코드나 LLVM IR과 유사한 역할)
- **AST (Abstract Syntax Tree, 추상 구문 트리)**:
  YAML Frontmatter 메타데이터, 마크다운 본문, 코드 블록, 조건부 활성화(Glob) 규칙 등의 텍스트 구문을 컴퓨터가 분석·변환하기 쉽도록 **계층적 트리(Tree) 구조의 데이터 객체**로 구조화한 모델입니다.
- **JIT (Just-In-Time, 온디맨드 동적 적재)**:
  모든 코딩 규칙을 에이전트 시작 시점에 무조건 주입(Always On)하지 않고, 특정 도구나 복잡한 절차가 필요한 시점에만 **필요한 순간에 동적으로 컨텍스트에 로드하고 작업 후 해제**하여 토큰 소모를 최소화하는 기법입니다.
- **MDC (Markdown with Config)**:
  Cursor AI에서 채택한 규칙 파일 확장자(`.mdc`)로, YAML 형식의 프론트매터(적용 대상 Glob 패턴, 상시 적용 여부 등)와 마크다운 지침 본문이 결합된 형태입니다.
- **MCP (Model Context Protocol)**:
  AI 모델이 외부 데이터베이스, 파일 시스템, API 등의 도구(Tool)를 표준화된 방식으로 안전하게 호출할 수 있도록 정의된 개방형 프로토콜입니다.
- **3-Way Diff (3자 비교 병합)**:
  원본 SSOT, 이전 캐시(Base), 특정 IDE에서 사용자가 직접 수정한 파일(Target)의 3가지 상태를 비교하여, 충돌 없이 변경사항만을 정확히 병합·역동기화하는 알고리즘입니다.

---

## 2. 시스템 아키텍처 개요 (System Architecture)

`token-hop`은 플랫폼 독립적인 **Universal Agent Intermediate Representation (UA-IR, 범용 에이전트 중간 표현식)** 추상 구문 트리(AST)를 기반으로 작동하는 3단계 컴파일러 파이프라인(Parser $\rightarrow$ UA-IR AST $\rightarrow$ Target Emitter) 및 양방향 동기화 엔진입니다.

```
                       [ Single Source of Truth: AGENTS.md / UA-IR ]
                                              │
                                              ▼
                        ┌───────────────────────────────────────────┐
                        │          token-hop Lexer & Parser         │
                        │    (Frontmatter + AST Tree Builder)       │
                        └─────────────────────┬─────────────────────┘
                                              │
                                              ▼
                        ┌───────────────────────────────────────────┐
                        │      Universal Agent IR (UA-IR AST)       │
                        │ (Skill | Rule | Workflow | Agent | Prompt)│
                        └─────────────────────┬─────────────────────┘
                                              │
                   ┌──────────────────────────┴──────────────────────────┐
                   ▼                                                     ▼
    ┌─────────────────────────────┐                       ┌─────────────────────────────┐
    │  Context Budgeting Engine   │                       │   Bi-Directional 3-Way      │
    │ (Token/Char Slicing & JIT)  │                       │   Diff & Drift Sync Engine  │
    └──────────────┬──────────────┘                       └──────────────┬──────────────┘
                   └──────────────────────────┬──────────────────────────┘
                                              │
                                              ▼
                        ┌───────────────────────────────────────────┐
                        │        Target Emitter & Code Gen          │
                        └──────┬──────────┬──────────┬──────────┬───┘
                               │          │          │          │
         ┌─────────────────────┘          │          │          └─────────────────────┐
         ▼                                ▼          ▼                                ▼
┌──────────────────┐             ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│Google Antigravity│             │    Cursor AI     │ │   Claude Code    │ │  GitHub Copilot  │
│.agents/rules/*.md│             │.cursor/rules/*.md│ │./CLAUDE.md       │ │.github/copilot-  │
│.agents/workflows/│             │(globs/always)    │ │.claude/skills/   │ │ instructions.md  │
│.agents/skills/   │             └──────────────────┘ └──────────────────┘ │.github/prompts/  │
└──────────────────┘                                                       └──────────────────┘
```

---

## 3. 이종 시스템 간 핵심 엔티티 매핑 (Dimensional Mapping Model)

AI 개발 지원 환경의 설정 요소를 **6대 핵심 엔티티**로 분류하고 정규화 매핑합니다.

| 엔티티 (Entity) | 개념적 역할 및 정의 | Google Antigravity | Cursor AI | Claude Code | GitHub Copilot |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Rule** | 프로젝트 전반 또는 파일별 코딩 제약 규칙 | `.agents/rules/*.md`<br>(Glob / Always On) | `.cursor/rules/*.mdc`<br>(`globs`, `alwaysApply`) | `./CLAUDE.md`<br>디렉터리 캐스케이딩 | `.github/instructions/*.md`<br>(`applyTo: glob`) |
| **Skill** | 특정 작업 시 온디맨드(JIT) 동적 적재 패키지 | `.agents/skills/<name>/`<br>`SKILL.md` (JIT Load) | `.cursor/rules/<name>.mdc`<br>(Apply Intelligently) | `.claude/skills/<name>/`<br>`SKILL.md` | `.github/skills/<name>/`<br>`SKILL.md` |
| **Workflow** | 다단계 절차 및 검증 루프를 갖춘 프롬프트 체인 | `.agents/workflows/*.md`<br>(`/workflow-name`) | `@rule` 참조 기반 프롬프트 | CLI Lifecycle Hooks &<br>순차 프롬프트 문서 | `.github/prompts/*.prompt.md`<br>(`/scaffold` 등) |
| **Instruction** | 페르소나 및 시스템 차원 기본 글로벌 지침 | `~/.gemini/GEMINI.md`<br>+ 프로젝트 `AGENTS.md` | `.cursor/rules/base.mdc`<br>(`alwaysApply: true`) | `~/.claude/CLAUDE.md`<br>+ 루트 `CLAUDE.md` | `.github/copilot-instructions.md` |
| **Prompt** | 재사용 가능한 개별 작업 템플릿 | 워크플로우 내 파라미터형<br>Prompt Artifact | 채팅 `@` 멘션 및 템플릿 | 단일 텍스트 프롬프트 파일 | `.github/prompts/*.prompt.md` |
| **Agent** | 역할, 도구 권한, 특화 영역을 갖춘 독립 에이전트 | Subagent 정의 및<br>Hook 바인딩 | MCP 바인딩 +<br>특화 시스템 프롬프트 | `.claude/subagents/` &<br>권한 설정 | `.github/agents/*.agent.md` |

### 3.1 실사례 기반 마이그레이션 패턴: GitHub Copilot `.github/` $\rightarrow$ Google Antigravity (Cloud-Barista `cm-beetle` 사례)

실제 오픈소스 멀티클라우드 프로젝트(예: `cloud-barista/cm-beetle`)에서 GitHub Copilot의 잦은 토큰 한도 초과 및 사용량 제한(Token & Usage Limits)으로 인해 Google Antigravity로 환경을 전환할 때 발생하는 변환 매핑 및 해결책은 다음과 같습니다.

```
[GitHub Copilot 기존 설정 (.github/)]                     [Google Antigravity 대상 구조]
┌──────────────────────────────────────────────┐        ┌──────────────────────────────────────────────┐
│ .github/copilot-instructions.md              │ ───►   │ AGENTS.md (Root SSOT)                        │
│ (레포지토리 전역 최상위 지침)                │        │ .agents/rules/00-global-guidelines.md        │
├──────────────────────────────────────────────┤        ├──────────────────────────────────────────────┤
│ .github/instructions/go-conventions.md       │ ───►   │ .agents/rules/go-conventions.md              │
│ (applyTo: "**/*.go" 경로 기반 규칙)          │        │ (globs: ["**/*.go"] Frontmatter 정규화)      │
├──────────────────────────────────────────────┤        ├──────────────────────────────────────────────┤
│ .github/prompts/scaffold-api.prompt.md       │ ───►   │ .agents/workflows/scaffold-api.md            │
│ (슬래시 커맨드형 API 스캐폴딩 템플릿)        │        │ (/scaffold-api 워크플로우 체이닝 변환)       │
├──────────────────────────────────────────────┤        ├──────────────────────────────────────────────┤
│ .github/agents/cloud-architect.agent.md      │ ───►   │ .agents/skills/cloud-architect/SKILL.md      │
│ (도구 권한 및 특화 모드 페르소나 지침)       │        │ (온디맨드 JIT Skill 적재로 토큰 80% 절감)    │
└──────────────────────────────────────────────┘        └──────────────────────────────────────────────┘
```

#### 주요 변환 메커니즘 및 기술적 해소 방안:
1. **`copilot-instructions.md` $\rightarrow$ `AGENTS.md` & Always On Rule**:
   프로젝트 아키텍처, 빌드/테스트 명령어, 전역 규칙을 추출하여 프로젝트 루트의 `AGENTS.md`로 일원화하고, Antigravity가 자동으로 인식하도록 연결합니다.
2. **`applyTo` 프론트매터 $\rightarrow$ Glob 기반 `.agents/rules/*.md`**:
   Copilot의 `applyTo: "**/*.go"` 단일 문자열 구문을 파싱하여 Antigravity 표준 Glob 패턴 배열로 정규화 트랜스파일링합니다.
3. **`.prompt.md` $\rightarrow$ 워크플로우 체이닝 (`.agents/workflows/*.md`)**:
   단순 텍스트 주입 방식의 Copilot 프롬프트를 Antigravity의 단계별(Planning $\rightarrow$ Execution $\rightarrow$ Verification Loop) 슬래시 명령어 워크플로우로 자동 변환합니다.
4. **`agents/*.agent.md` / 커스텀 모드 $\rightarrow$ 온디맨드 JIT `SKILL.md`**:
   상시 컨텍스트를 차지하던 에이전트 지침을 필요할 때만 호출되는 JIT Skill로 전환하여, Antigravity 턴(Turn) 시작 시 토큰 낭비를 원천 차단합니다.

---

## 4. Universal Agent IR (UA-IR v1.0.0) 사양

### 4.1 스키마 정의
UA-IR은 YAML Frontmatter와 정규화된 Markdown Payload로 구성됩니다.

```yaml
ua_version: "1.0.0"
metadata:
  id: "rule-react-component-standards"
  type: "Rule" # [Rule | Skill | Workflow | Instruction | Prompt | Agent]
  name: "React Component Conventions"
  description: "Enforces strict TypeScript functional patterns and named exports for UI components."
  version: "1.0.0"
  author: "Token-Hop Project"
  tags: ["frontend", "react", "typescript"]

activation:
  mode: "Glob" # [AlwaysOn | Glob | ModelDecision | OnDemand | ManualMention]
  globs:
    - "src/components/**/*.tsx"
    - "src/ui/**/*.tsx"
  exclude_globs:
    - "**/*.test.tsx"
    - "**/*.stories.tsx"
  triggers:
    - "create component"
    - "refactor UI"
  slash_command: "scaffold-component"

context_budget:
  priority: "High" # [Critical | High | Medium | Low]
  max_tokens: 400
  max_characters: 1600
  decompose_strategy: "JIT_Skill" # [JIT_Skill | Truncate | Split_Glob]

bindings:
  allowed_tools: ["read_file", "write_to_file", "run_command"]
  mcp_servers: []

payload:
  markdown_body: |
    # React Component Conventions
    - Use functional components with explicit type definitions for Props.
    - Always use named exports (`export const ComponentName = ...`) instead of default exports.
    - Colocate styles and unit tests in the same directory.
```

---

## 5. 핵심 서브시스템 상세 설계

### 5.1 파서 및 어휘 분석기 (Parser & Lexer)
- **Multi-Format Ingestion**:
  - `AGENTS.md` (SSOT 마스터 파일)
  - `.cursor/rules/*.mdc` (Cursor MDC 포맷)
  - `.agents/rules/*.md`, `.agents/workflows/*.md` (Antigravity 포맷)
  - `CLAUDE.md`, `.claude/skills/` (Claude Code 포맷)
  - `.github/copilot-instructions.md`, `.github/instructions/*.md` (Copilot 포맷)
- 소스 파일의 Frontmatter(YAML/JSON)와 Markdown Body를 파싱하여 `UA-IR Document AST`로 변환.

### 5.2 컨텍스트 예산 최적화기 (Context Budgeting & JIT Decomposer)
- **타겟별 용량 제약 검사**:
  - Cursor: 파일당 6,000자 / 전체 12,000자 제한
  - Claude Code: 200줄 미만 권장
  - Antigravity: 파일당 12,000자 제한
- **JIT 분할 및 스킬화**:
  - AlwaysOn 지침이 예산을 초과하는 경우, 독립적인 작업 절차를 추출하여 Antigravity의 `.agents/skills/<name>/SKILL.md` 또는 Cursor의 온디맨드 `.mdc`로 분할 컴파일.

### 5.3 타겟 코드 생성기 (Target Emitter Engine)
- **Antigravity Emitter**:
  - AlwaysOn $\rightarrow$ `.agents/rules/` (Frontmatter 없는 순수 Markdown 또는 Always On 메타데이터)
  - Glob $\rightarrow$ `.agents/rules/*.md` (YAML Frontmatter 내 `globs` 매핑)
  - Workflow $\rightarrow$ `.agents/workflows/*.md` (`/workflow-name` 슬래시 명령어 바인딩)
  - Skill $\rightarrow$ `.agents/skills/<name>/SKILL.md` (JIT On-Demand 구조)
- **Cursor Emitter**:
  - `.cursor/rules/<name>.mdc` (`description`, `globs`, `alwaysApply` YAML Frontmatter)
- **Claude Code Emitter**:
  - `./CLAUDE.md`, `.claude/skills/<name>/SKILL.md`
- **GitHub Copilot Emitter**:
  - `.github/copilot-instructions.md`, `.github/instructions/*.instructions.md` (`applyTo` 프론트매터)

### 5.4 양방향 동기화 및 3-Way Diff 엔진 (Bi-directional Sync & Drift Control)
- `.token-hop/cache.json`에 각 파일의 생성 시점 해시(SHA-256) 저장.
- 특정 IDE 환경에서 규칙 수정 시 파일 변경 이벤트(FS Watcher) 감지.
- **3-Way Diff**를 통해 원본 SSOT(`AGENTS.md`)와 타겟 변경사항을 충돌 없이 병합하고 타겟 환경 전체로 재컴파일 배포.

---

## 6. 품질 검증 파이프라인 (Promptfoo Integration)

변환된 규칙이 실제 LLM에서 행동 지침을 충실히 따르는지 자동으로 검증합니다.

```yaml
# promptfooconfig.yaml
description: "Token-Hop Cross-Agent Transpilation Precision Evaluation"

prompts:
  - "file://.agents/rules/react-components.md"
  - "file://.cursor/rules/react-components.mdc"
  - "file://.github/instructions/react-components.instructions.md"

providers:
  - id: "anthropic:messages:claude-3-7-sonnet-latest"
  - id: "google:gemini-2.5-pro"
  - id: "openai:chat:gpt-4o"

tests:
  - description: "Verify named export compliance"
    vars:
      input: "Create a React button component in src/components/Button.tsx"
    assert:
      - type: "icontains"
        value: "export const Button"
      - type: "not-icontains"
        value: "export default"
      - type: "llm-rubric"
        value: "The response must not contain default exports and must strictly type props."
```

### 평가 지표 (Metrics)
1. **구문 정확성 (Syntax Accuracy)**: Frontmatter 스키마 파싱 오류 0건 (`regex` / `is-json`).
2. **제약 조건 이행률 (Constraint Compliance)**: 금지 패턴 위반 여부 (`llm-rubric` $\ge 0.95$).
3. **토큰 예산 효율성 (Token Efficiency)**: 단일 규칙 파일 400 토큰 미만 유지.
4. **시맨틱 동등성 (Semantic Parity)**: 원본 SSOT와 컴파일된 결과물 간 코사인 유사도 $\ge 0.90$.

---

## 7. 하이브리드 엔진 설계 및 선택적 생성형 AI 증강 (Hybrid Engine & Selective AI Augmentation)

`token-hop`은 **결정론적 컴파일 코어(Deterministic Core)**를 기반으로 빠르고 비용 없이 동작하며, 사용자의 필요에 따라 **Claude, Gemini, OpenAI, Ollama(로컬 LLM)** 등 생성형 AI API를 선택적으로 결합하는 **하이브리드 아키텍처**를 채택합니다.

```
                  ┌─────────────────────────────────────────────────────────┐
                  │                 token-hop Hybrid Engine                 │
                  └────────────────────────────┬────────────────────────────┘
                                               │
                      ┌────────────────────────┴────────────────────────┐
                      ▼                                                 ▼
        ┌───────────────────────────┐                     ┌───────────────────────────┐
        │  Deterministic Core       │                     │  Optional AI Augmentation │
        │  (Rule-based / 100% Free) │                     │  (--ai / .env Enabled)    │
        ├───────────────────────────┤                     ├───────────────────────────┤
        │ • AST Lexing & Parsing    │                     │ • Semantic Decomposer     │
        │ • Glob/Frontmatter Normal │                     │ • Smart Trigger Gen       │
        │ • Heading-based Slicing   │                     │ • 3-Way Semantic Merge    │
        │ • Static File Generation  │                     │ • Promptfoo LLM Rubric    │
        │ • Execution Time: < 10ms  │                     │ • Multi-Provider Adapter  │
        └───────────────────────────┘                     └─────────────┬─────────────┘
                                                                        │
                                       ┌─────────────────┬──────────────┴──┬─────────────────┐
                                       ▼                 ▼                 ▼                 ▼
                                ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐
                                │Google Gemini│   │  Anthropic  │   │   OpenAI    │   │ Local LLM   │
                                │(2.5 Pro)    │   │(Claude 3.7) │   │  (GPT-4o)   │   │  (Ollama)   │
                                └─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘
```

### 7.1 기본 결정론적 모드 (Deterministic Core)
- **API 호출 0회, 비용 0원, 네트워크 불필요**:
  별도의 LLM API 키 없이도 AST 파싱, Frontmatter 변환, Glob 패턴 정규화, 마크다운 렌더링, 정적 파일 생성을 **완전 무과금·무지연(10ms 이내)**으로 수행합니다.
- **CI/CD 친화성**: API 사용량 한도나 네트워크 불안정 없이 일관된 컴파일 결과를 재현합니다.

### 7.2 선택적 AI 증강 태스크 (AI-Augmented Capabilities)
사용자가 `.env` 또는 `--ai` 플래그로 AI 연동을 활성화할 경우 다음과 같은 고도화 기능을 수행합니다.

1. **지능형 의미론적 규칙 분할 (Semantic Decomposer)**:
   단순 줄/단어 수 기준 분할을 넘어, LLM이 텍스트의 맥락을 분석하여 독립적인 서브 규칙 및 온디맨드 JIT `SKILL.md`로 자동 모듈화합니다.
2. **시맨틱 트리거 및 설명 자동 생성 (Smart Trigger & Description Generator)**:
   Antigravity의 `Model Decision` 및 Cursor의 `Apply Intelligently`에 필요한 고품질 `description`과 활성화 키워드를 규칙 본문에서 자동 추출·생성합니다.
3. **지능형 3-Way 충돌 해결 (AI-Powered Conflict Resolution)**:
   동일한 규칙이 여러 IDE에서 상이하게 수정되어 3-Way 병합 충돌이 발생할 경우, AI가 양쪽의 의도를 종합하여 최적의 병합안을 제안합니다.
4. **Promptfoo 품질 평가 및 LLM Rubric 채점 (Automated Quality Evaluation)**:
   컴파일된 규칙이 타겟 LLM에서 지침을 100% 준수하는지 실시간 호출 및 채점을 수행합니다.

### 7.3 환경 변수 및 멀티 프로바이더 설정 (`.env.example`)
루트 디렉터리의 `.env` 파일을 통해 공급자(Provider)를 유연하게 교체할 수 있습니다.

```bash
# Token-Hop AI 설정 (.env)
TOKEN_HOP_AI_ENABLED=true
TOKEN_HOP_AI_PROVIDER=gemini # [gemini | anthropic | openai | ollama]
TOKEN_HOP_AI_MODEL=gemini-2.5-pro

# Provider별 API Key
GEMINI_API_KEY=your_gemini_api_key
ANTHROPIC_API_KEY=your_anthropic_api_key
OPENAI_API_KEY=your_openai_api_key

# 로컬 오프라인 LLM (API 비용 0원 & 사내 데이터 보안)
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL=deepseek-r1:14b
```

### 7.4 무중단 자동 폴백 (Graceful Fallback)
API 키가 설정되지 않았거나, 할당량 초과(Quota Exceeded) 또는 오프라인 환경인 경우 시스템은 오류로 중단되지 않고 **자동으로 결정론적(Deterministic) 정적 컴파일 모드로 전환**됩니다.

### 7.5 다중 생성형 AI 통합 표준 인터페이스 (`pkg/ai/`)
다양한 LLM 프로바이더를 추상화한 `AIProvider` 표준 인터페이스를 제공하여, 엔진 코드는 변경 없이 프로바이더 구현체만 교체됩니다.

```go
type AIProvider interface {
    Name() string
    Model() string
    DecomposeRule(ctx context.Context, title, content string, maxTokens int) (*DecomposeResult, error)
    GenerateDescription(ctx context.Context, content string) (string, error)
}
```

- **구현 프로바이더**:
  - `GeminiProvider`: Google Gemini REST API v1beta (Gemini 2.5 Pro / Flash)
  - `ClaudeProvider`: Anthropic Claude Messages API (Claude 3.5 / 3.7 Sonnet)
  - `OpenAIProvider`: OpenAI Chat Completions API (GPT-4o / o3-mini)
  - `OllamaProvider`: Local Ollama REST API (DeepSeek-R1 / Qwen 2.5 / Llama 3)
  - `MockProvider`: 오프라인 단위 테스트 및 CI/CD 결정론적 목업 엔진

---

## 8. CLI 도구 (`token-hop`) 설계

### 8.1 주요 명령어 명세
```bash
# 1. 프로젝트 초기화 (AGENTS.md SSOT 스캐폴딩 생성)
token-hop init [--template react | fullstack | python | go]

# 2. 기존 프로젝트 설정 가져오기 (GitHub Copilot, Cursor, Claude 설정 -> SSOT 자동 생성)
token-hop import --from copilot .github/
token-hop import --from cursor .cursor/
token-hop import --from claude CLAUDE.md

# 3. 크로스 에이전트 타겟 컴파일 (기본: 결정론적 코어, --ai: 지능형 증강)
token-hop compile --target antigravity,cursor,claude,copilot [--ai] [--watch]

# 4. 다이렉트 변환 (단일 명령으로 Copilot -> Antigravity 즉시 변환)
token-hop convert --from copilot --to antigravity [--dry-run] [--generate-prompt]

# 5. 크로스 에이전트 설정 매핑 및 규격 사전 시각화 (테이블 출력 및 파일 저장)
token-hop describe --from copilot --to antigravity [--out plan.md | --spec | --format markdown|json]

# 6. 타겟 맞춤형 2차 AI 정제 프롬프트 생성 (Antigravity // turbo, Cursor @rule 등)
token-hop prompt --to antigravity [--file <path>] [--dir <dir>] [--out <output_path>]

# 7. 실시간 변경 감지 및 양방향 동기화 데몬
token-hop watch --ssot AGENTS.md

# 8. Promptfoo 기반 변환 규칙 품질 검증
token-hop eval --config promptfooconfig.yaml --threshold 0.95

# 9. 토큰 예산 및 컨텍스트 윈도우 정밀 감사
token-hop audit --max-tokens 500
```

### 8.2 프로젝트 설정 파일 (`token-hop.yaml`)
```yaml
version: "1.0"
ssot: "AGENTS.md"
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
context_budget:
  max_tokens_per_rule: 400
  max_characters_per_rule: 1800
  auto_decompose: true
ai:
  enabled: false # CLI 실행 시 --ai 플래그 또는 .env 설정으로 활성화
  provider: "gemini" # [gemini | anthropic | openai | ollama]
  model: "gemini-2.5-pro"
```

---

## 9. 구현 로드맵 및 패키지 구조

### 9.1 단계별 마일스톤 (Milestones)

- **Phase 1: Core AST & Dual Target Transpiler (v0.1.0)**
  - UA-IR 핵심 데이터 구조 및 YAML Frontmatter 파서
  - SSOT (`AGENTS.md`) 렉서 및 AST 변환기
  - Google Antigravity (`.agents/rules/`, `.agents/workflows/`) Emitter
  - Cursor AI (`.cursor/rules/*.mdc`) Emitter
  - CLI `token-hop init`, `token-hop compile` 구현

- **Phase 2: Full Multi-Agent Support & Context Budgeting (v0.2.0)**
  - Claude Code (`CLAUDE.md`, `.claude/skills/`) Emitter
  - GitHub Copilot (`.github/instructions/`, `.github/prompts/`) Emitter
  - Gemini CLI (`GEMINI.md`) & Roo Code (`.roomodes`) Emitter
  - Context Budgeting Engine (JIT Skill 분할 및 토큰 최적화)

- **Phase 3: Bi-directional Sync & Quality Engine (v0.3.0)**
  - File Watcher & 3-Way Diff Engine (`.token-hop/cache.json`)
  - Promptfoo 통합 테스트 러너 (`token-hop eval`)
  - CI/CD GitHub Actions Workflow 템플릿 배포

### 9.2 디렉터리 및 패키지 구조
```
token-hop/
├── cmd/
│   └── token-hop/           # CLI 진입점 (main.go / cli.py)
├── pkg/
│   ├── ir/                  # UA-IR 데이터 모델 및 스키마
│   ├── parser/              # Markdown, MDC, YAML Frontmatter 파서
│   │   ├── agents_md.go
│   │   ├── mdc.go
│   │   └── frontmatter.go
│   ├── emitter/             # 타겟 플랫폼별 코드 생성기
│   │   ├── antigravity.go
│   │   ├── cursor.go
│   │   ├── claude.go
│   │   └── copilot.go
│   ├── audit/               # 토큰 산출, 컨텍스트 감사 및 JIT 분할기
│   ├── ai/                  # 다중 LLM 프로바이더 통합 인터페이스
│   ├── sync/                # 3-Way Diff 및 파일 감시 데몬
│   └── eval/                # Promptfoo 연동 및 품질 평가 모듈
├── docs/
│   ├── README.md            # 문서 목차 (TOC)
│   ├── design/              # 아키텍처 및 스펙 문서
│   │   └── token-hop-design-plan.md
│   └── maintenance/         # 스펙 최신화 및 프롬프트 가이드
│       └── SPEC_EVOLUTION_PROMPT.md
├── token-hop.yaml           # 프로젝트 기본 설정
└── README.md
```

---

## 10. 공식 레퍼런스 및 스펙 문서 (Official References)

1. **Google Antigravity / Gemini Code Assist**:
   - [Google Cloud Gemini Code Assist Docs](https://cloud.google.com/gemini/docs/codeassist)
   - Antigravity Agents, Rules (`.agents/rules/`), Workflows (`.agents/workflows/`), Skills (`.agents/skills/`)
2. **Cursor AI**:
   - [Cursor Rules Documentation](https://docs.cursor.com/context/rules-for-ai)
   - [Cursor Context & MDC Spec](https://docs.cursor.com/)
3. **GitHub Copilot**:
   - [GitHub Copilot Custom Instructions](https://docs.github.com/en/copilot/customizing-copilot/adding-custom-instructions-for-github-copilot)
   - [GitHub Copilot Prompts & Scopes](https://docs.github.com/en/copilot)
4. **Claude Code (Anthropic)**:
   - [Claude Code Overview](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview)
   - [Claude Code Memory (`CLAUDE.md`) & Tools](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/memory)
5. **Windsurf (Codeium)**:
   - [Windsurf Rules & Memories](https://docs.codeium.com/windsurf/memories)
6. **Roo Code / Cline**:
   - [Roo Code Custom Modes (.roomodes)](https://github.com/RooVetGit/Roo-Code)
   - [Roo Code Official Docs](https://docs.roocode.com)

