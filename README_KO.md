# 🏃 agent-parkour (`parkour` / `pk`)

[English](README.md) | **[한국어](README_KO.md)**

> **"토큰 장벽에 부딪히셨나요? 가볍게 파쿠르 하듯 에이전트를 뛰어넘으세요!"** 🏃💨  
> *"Hit a token wall? Just Parkour across your AI coding agents!"*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](file:///home/ubuntu/dev/yunkon-kim/token-hop/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](file:///home/ubuntu/dev/yunkon-kim/token-hop/go.mod)
[![Architecture: UA-IR](https://img.shields.io/badge/Architecture-UA--IR%20v1.0-brightgreen)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/design/agent-parkour-design-plan.md)

> ⚡ **완전 무료 & 로컬 동작**: `agent-parkour`는 기본적으로 **API 키 없이 100% 로컬에서 무료**로 동작합니다. 생성형 AI 연동은 고급 규칙 분할을 위한 순수 **선택적·실험적(Optional & Experimental)** 기능입니다.

---

`agent-parkour`(`parkour` / `pk`)는 **Claude Code, Cursor, GitHub Copilot, Google Antigravity, Gemini CLI, Roo Code** 등 다양한 AI 코딩 도구 간의 프롬프트, 규칙(Rules), 워크플로우(Workflows), 스킬(Skills) 설정 이질성을 해결하는 **초고속 크로스 에이전트 프롬프트 컴파일러이자 컨텍스트 동기화 도구**입니다.

AI 도구별 토큰 제한(Token Limit)이나 사용량 한도(Usage Cap)라는 벽에 부딪혀 다른 도구로 전환할 때마다 지침을 일일이 복사하고 형식을 맞추느라 고생할 필요가 없습니다. 단 하나의 마스터 기준 문서(`AGENTS.md`)만 수정하면, 모든 AI 도구용 설정으로 **0.01초 만에 자동 컴파일**됩니다.

---

## ⚡ 왜 agent-parkour인가요? (Why agent-parkour?)

1. 🎯 **단일 진실 원천 (Single Source of Truth, SSOT)**  
   프로젝트 루트의 `AGENTS.md` 파일 하나로 Antigravity(`.agents/`), Cursor(`.cursor/rules/`), Copilot(`.github/`), Claude(`CLAUDE.md`) 설정을 중앙 집중 동기화합니다.

2. 🚀 **Zero-Cost & 초고속 결정론적(Deterministic) 코어**  
   복잡한 AST 파싱 및 코드 생성이 외부 LLM API 호출 없이 **로컬에서 10ms 이내(무과금)**로 즉시 실행됩니다.

3. 🛡️ **자동 타임스탬프 백업 (`*.bak_YYYYMMDD_HHMMSS`)**  
   기존 파일이 덮어쓰여질 때마다 이전 원본을 타임스탬프와 함께 자동 백업하여 작업 손실을 원천 차단합니다.

4. 🧩 **지능형 토큰 최적화 & JIT Skill 자동 분할 (Context Optimization)**  
   Cursor 6,000자, Antigravity 12,000자 등 플랫폼별 용량 한계를 자동 감지하고, 대용량 지침은 온디맨드 JIT `SKILL.md`로 분해하여 **상시 토큰 소모를 80% 이상 절감**합니다.

5. 🔄 **1초 만에 크로스 에이전트 마이그레이션 (`parkour convert`)**  
   GitHub Copilot(`.github/`)을 사용하던 기존 프로젝트를 단 한 줄의 명령어로 Google Antigravity(`.agents/`) 구조로 완벽 변환합니다.

---

## ⌨️ 빠른 시작 (Quick Start)

### 1. 설치 (Installation)

단 몇 초 만에 설치하여 바로 사용할 수 있습니다:

```bash
# Go를 통한 설치 (권장)
go install github.com/yunkon-kim/agent-parkour/cmd/parkour@latest

# 또는 설치 스크립트로 바이너리 직접 다운로드 (Linux / macOS)
curl -fsSL https://raw.githubusercontent.com/yunkon-kim/agent-parkour/main/install.sh | bash

# 초고속 입력을 위한 단축어 설정 (설치 스크립트가 자동 구성)
alias pk=parkour
```

### 2. 바로 사용하기 (결정론적 코어 모드)

어떤 프로젝트 디렉터리에서든 추가 설정 없이 즉시 실행 (`parkour` 또는 단축어 `pk` 사용 가능):

```bash
# 1) 변환 전 소스 -> 타겟 파일 매핑 계획을 테이블로 사전 확인
parkour describe --from copilot --to antigravity

# 2) 매핑 계획을 파일로 직접 내보내기 (확장자에 따라 Markdown / JSON 자동 서식 적용)
parkour describe --from copilot --to antigravity --out plan.md
parkour describe --from copilot --to antigravity --out plan.json

# 3) 플랫폼 간 6대 핵심 엔티티(규칙, 스킬, 워크플로우 등) 표준 규격 매트릭스 조회
parkour describe --from copilot --to antigravity --spec

# 4) GitHub Copilot 설정을 Google Antigravity로 자동 대응 변환 (알아서 모든 파일 매핑)
parkour convert --from copilot --to antigravity

# 5) 변환 완료와 동시에 타겟 맞춤형 2차 AI 정제 프롬프트(refine-prompt.md) 자동 생성
parkour convert --from copilot --to antigravity --generate-prompt

# 6) 파일 쓰기 없이 변환 결과를 시뮬레이션 (Dry-Run 모드)
parkour convert --from copilot --to antigravity --dry-run

# 7) GitHub Copilot 설정을 지원되는 모든 타겟(Antigravity, Cursor)으로 일괄 변환
parkour convert --from copilot --to all

# 8) 특정 AI 타겟으로 변환 (소스는 현재 레포에서 자동 감지)
parkour convert --to cursor

# 9) 타겟 맞춤형 2차 AI 정제 프롬프트 생성 (Antigravity // turbo, Cursor @rule 등 반영)
parkour prompt --to antigravity --file .agents/workflows/git-commit.md
parkour prompt --to antigravity --dir .agents/ --out .agents/refine-prompt.md
parkour prompt --to cursor --out .cursor/refine-prompt.md

# 10) 프로젝트 내 모든 규칙의 토큰 수 및 컨텍스트 한도 정밀 감사
parkour audit

# 11) 신규 프로젝트에 SSOT 템플릿(AGENTS.md & parkour.yaml) 초기화
parkour init
```

---

## 🧪 선택적 & 실험적 기능: 생성형 AI 증강 (Optional & Experimental)

> ⚠️ **안내**: 생성형 AI 기능은 **완전한 선택 사항(Optional)이자 실험적(Experimental) 기능**입니다. `agent-parkour`의 핵심 컴파일러는 어떠한 API 키도 요구하지 않고 100% 로컬에서 동작합니다.

대용량 복합 지침을 문맥 단위로 분석하여 모듈형 JIT Skill로 자동 분해하고자 할 때만 선택적으로 생성형 AI를 연동할 수 있습니다:

### 1. `.env` 설정 (선택 사항)
```bash
cp .env.example .env
# .env 파일에서 원하는 프로바이더 및 API 키 설정:
# PARKOUR_AI_ENABLED=true
# PARKOUR_AI_PROVIDER=gemini  # [gemini | claude | openai | ollama]
# GEMINI_API_KEY=your_key_here
```

### 2. 지능형 규칙 분할 실행
```bash
# Google Gemini를 연동하여 대용량 지침을 JIT Skill로 지능형 분할
parkour convert --from copilot --to antigravity --decompose --ai --provider gemini

# 또는 로컬 무료 LLM(Ollama)을 연동하여 API 비용 0원으로 분할
parkour convert --from copilot --to antigravity --decompose --ai --provider ollama
```

---

### 🛠️ 소스코드 직접 빌드 (기여자용)

```bash
git clone https://github.com/yunkon-kim/agent-parkour.git
cd agent-parkour
make build && make test
```

---

## 📊 변환 매핑 매트릭스 (Mapping Matrix)

| 엔티티 (Entity) | UA-IR 정규화 역할 | Google Antigravity | Cursor AI | GitHub Copilot | Claude Code |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Rule** | 파일 경로/정적 제약 | `.agents/rules/*.md` | `.cursor/rules/*.mdc` | `.github/instructions/*.md` | `./CLAUDE.md` |
| **Skill** | 온디맨드 JIT 동적 절차 | `.agents/skills/<name>/` | `.cursor/rules/*.mdc` | `.github/skills/<name>/` | `.claude/skills/` |
| **Workflow** | 다단계 절차 및 검증 루프 | `.agents/workflows/*.md` | `@rule` 체이닝 문서 | `.github/prompts/*.md` | CLI 프롬프트 체인 |
| **Instruction** | 전역 시스템 기본 지침 | `AGENTS.md` | `.cursor/rules/base.mdc`| `.github/copilot-instructions.md` | `CLAUDE.md` |

---

## 📚 상세 문서 (Documentation)

- [전체 문서 목차 (docs/README.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/README.md)
- [시스템 아키텍처 및 구현 계획 명세서 (docs/design/agent-parkour-design-plan.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/design/agent-parkour-design-plan.md)
- [스펙 최신화 및 프롬프트 가이드 (docs/maintenance/SPEC_EVOLUTION_PROMPT.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/maintenance/SPEC_EVOLUTION_PROMPT.md)
- [테스트 스위트 및 픽스처 안내 (test/README.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/test/README.md)

---

## 📄 라이선스 (License)

[MIT License](file:///home/ubuntu/dev/yunkon-kim/token-hop/LICENSE) © 2026 Yunkon Kim
