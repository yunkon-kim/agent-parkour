# Agent-Parkour Engine: Cross-Agent Prompt Compiler Architecture & Plan
**AI 개발 환경 간 설정 이질성 해소를 위한 크로스 에이전트 프롬프트 컴파일러 설계 및 오픈소스 기획**

- **프로젝트명**: `agent-parkour`
- **CLI 명칭**: `parkour` (약칭: `pk`)
- **버전**: v1.0.0
- **저자**: Yunkon Kim

---

## 1. 배경 및 문제 정의 (Background & Problem Definition)

### 1.1 멀티 에이전트 IDE 환경의 파편화 (Context Fragmentation)
최근 소프트웨어 개발 생태계에서는 **Google Antigravity, GitHub Copilot, Claude Code, Cursor AI** 등 다양한 AI 에이전트 도구가 병행 활용되고 있습니다. 개발자는 각 AI 도구의 토큰 제한 및 사용량 한도(Token & Usage Limits), 작업 특성에 따라 도구를 수시로 전환하며 작업합니다.

그러나 각 환경은 페르소나, 코딩 규칙(Rule), 워크플로우(Workflow), 스킬(Skill), 지시문(Instruction), 프롬프트(Prompt)를 정의하는 방식에서 극심한 파편화를 보입니다.

```
+-------------------------------------------------------------------------------+
|                       Multi-Agent Ecosystem Fragmentation                     |
+-------------------------------------------------------------------------------+
| [Google Antigravity] | .agents/rules/*.md, .agents/workflows/, .agents/skills/ |
| [GitHub Copilot]     | .github/copilot-instructions.md, prompts/, skills/     |
| [Claude Code]        | ./CLAUDE.md, ~/.claude/CLAUDE.md, .claude/skills/      |
| [Cursor AI]          | .cursor/rules/*.mdc (globs, alwaysApply, description)  |
+-------------------------------------------------------------------------------+
```

### 1.2 핵심 용어 및 약어 해설 (Glossary & Key Concepts)
- **SSOT (Single Source of Truth, 단일 진실 원천)**:
  여러 AI 도구마다 파편화된 설정을 단 하나의 마스터 기준 문서(`AGENTS.md`)로 중앙 집중화합니다.
- **IR (Intermediate Representation, 중간 표현식)**:
  플랫폼 중립적인 추상 데이터 중간 규격(Universal Agent IR, UA-IR)입니다.
- **AST (Abstract Syntax Tree, 추상 구문 트리)**:
  Frontmatter 메타데이터, 본문, Glob 규칙을 계층적 트리 구조로 모델링합니다.
- **JIT (Just-In-Time, 온디맨드 동적 적재)**:
  규칙을 항상 켜두지 않고, 필요한 작업 시점에만 동적으로 컨텍스트에 로드합니다.

---

## 2. 시스템 아키텍처 개요 (System Architecture)

`agent-parkour`는 플랫폼 독립적인 **Universal Agent Intermediate Representation (UA-IR)** 추상 구문 트리를 기반으로 작동하는 3단계 컴파일러 파이프라인(Parser $\rightarrow$ UA-IR AST $\rightarrow$ Target Emitter)입니다.

---

## 3. 이종 시스템 간 핵심 엔티티 매핑 (Dimensional Mapping Model)

| 지침 종류 (Entity) | 역할 및 성격 | Google Antigravity | GitHub Copilot | Claude Code | Cursor AI |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Instruction** | 프로젝트 전반에 상시 적용되는 전역 헌장 | `AGENTS.md`<br>*(또는 `.agents/AGENTS.md`)* | `.github/copilot-instructions.md` | `CLAUDE.md`<br>*(루트 전역)* | `.cursorrules`<br>*(또는 `base.mdc`)* |
| **Rule** | 파일·언어 편집 시 적용되는 코딩 규칙 | `.agents/rules/*.md`<br>(`globs`) | `.github/instructions/*.md`<br>(`applyTo` 문자열/배열) | 서브디렉터리 `CLAUDE.md`<br>(디렉터리 스코프) | `.cursor/rules/*.mdc`<br>(`globs`) |
| **Prompt** | `/<name>` 슬래시 커맨드 작업 템플릿 | `.agents/workflows/*.md` | `.github/prompts/*.prompt.md`<br>(`/<cmd>`) | `.claude/workflows/*.md`<br>(`/<cmd>`) | `.cursor/rules/*.mdc` |
| **Workflow** | `/<name>` 슬래시 커맨드 다단계 절차 | `.agents/workflows/*.md`<br>(`/<cmd>`) | `.github/prompts/*.prompt.md` | `.claude/workflows/*.md` | `.cursor/rules/*.mdc` |
| **Skill** | AI가 자율적으로 꺼내보는 전문 참고서/팩 | `.agents/skills/<name>/`<br>`SKILL.md` (JIT) | `.github/skills/<name>/`<br>`SKILL.md` (JIT) | `.claude/skills/<name>/`<br>`SKILL.md` (JIT) | `.cursor/rules/*.mdc`<br>(Description 기반 로딩) |
| **Agent** | 독립 권한과 도구를 갖춘 특화 서브에이전트 | `.agents/skills/` (또는 Agent) | `.github/agents/*.agent.md` | `.claude/agents/` | `.cursor/rules/*.mdc` |

---

## 4. CLI 도구 (`agent-parkour` / `pk`) 명세

```bash
# 1. 크로스 에이전트 설정 매핑 및 사전 검토
parkour describe -i /path/to/repo --from copilot --to antigravity

# 2. 다이렉트 변환
parkour convert --from copilot --to antigravity [--dry-run] [--gen-refine-prompt]

# 3. 토큰 예산 및 컨텍스트 윈도우 감사
parkour audit --input .github/
```
