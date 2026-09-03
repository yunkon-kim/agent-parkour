# Agent-Parkour Documentation

이 디렉터리는 **Agent-Parkour Engine** (크로스 에이전트 프롬프트 컴파일러)의 설계, 아키텍처 및 기술 문서를 포함합니다.

## 목차 (Table of Contents)

### 1. 설계 및 아키텍처 (Design & Architecture)
- [docs/design/agent-parkour-design-plan.md](design/agent-parkour-design-plan.md): 
  - 핵심 용어 및 약어 해설 (SSOT, IR, AST, JIT, MDC, MCP, 3-Way Diff)
  - 멀티 에이전트(Antigravity, Cursor, Claude Code, Copilot, Gemini CLI, Roo Code) 파편화 분석 및 SSOT 아키텍처
  - 6대 핵심 엔티티(Rule, Skill, Workflow, Instruction, Prompt, Agent) 차원 변환 모델
  - Universal Agent IR (UA-IR v1.0.0) 스키마 명세
  - 컴파일러 파이프라인 (Parser $\rightarrow$ Context Budgeting & JIT Decomposer $\rightarrow$ Target Emitter)
  - 3-Way Diff 기반 양방향 동기화 (Bi-directional Sync)
  - Promptfoo 기반 품질 평가 체계 (CI/CD)
  - CLI 명세 (`parkour init/convert/describe/prompt/audit`) 및 구현 로드맵

### 2. 유지보수 및 스펙 최신화 (Maintenance & Spec Evolution)
- [docs/maintenance/SPEC_EVOLUTION_PROMPT.md](maintenance/SPEC_EVOLUTION_PROMPT.md):
  - IDE/LLM 도구들의 경로, 확장자, 프론트매터 스키마 변경 시 `agent-parkour` 엔진을 체계적으로 업데이트하기 위한 **스펙 진화 프롬프트 템플릿 및 6단계 체크리스트**
  - Google Antigravity 워크플로우: `.agents/workflows/update-agent-specs.md` (`/update-agent-specs`)
