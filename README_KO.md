# 🦘 token-hop (`thop`)

[English](README.md) | **[한국어](README_KO.md)**

> **"토큰 제한에 걸리셨나요? 가볍게 Hop!"** 💨  
> *"Hit a token limit? Just Hop across your AI coding agents!"*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](file:///home/ubuntu/dev/yunkon-kim/token-hop/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](file:///home/ubuntu/dev/yunkon-kim/token-hop/go.mod)
[![Architecture: UA-IR](https://img.shields.io/badge/Architecture-UA--IR%20v1.0-brightgreen)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/design/token-hop-design-plan.md)

> ⚡ **완전 무료 & 로컬 동작**: `token-hop`은 기본적으로 **API 키 없이 100% 로컬에서 무료**로 동작합니다. 생성형 AI(Claude/Gemini/OpenAI/Ollama)는 지능형 규칙 분할이 필요할 때만 선택적으로 연동됩니다.

---

`token-hop`(`thop`)은 **Claude Code, Cursor, GitHub Copilot, Google Antigravity, Gemini CLI, Roo Code** 등 다양한 AI 코딩 도구 간의 프롬프트, 규칙(Rules), 워크플로우(Workflows), 스킬(Skills) 설정 이질성을 해결하는 **초고속 크로스 에이전트 프롬프트 컴파일러이자 컨텍스트 동기화 도구**입니다.

AI 도구별 토큰 제한(Token Limit)이나 사용량 한도(Usage Cap)에 도달하여 다른 도구로 전환할 때마다 지침을 일일이 복사하고 형식을 맞추느라 고생할 필요가 없습니다. 단 하나의 마스터 기준 문서(`AGENTS.md`)만 수정하면, 모든 AI 도구용 설정으로 **0.01초 만에 자동 컴파일**됩니다.

---

## ⚡ 왜 token-hop인가요? (Why token-hop?)

1. 🎯 **단일 진실 원천 (Single Source of Truth, SSOT)**  
   프로젝트 루트의 `AGENTS.md` 파일 하나로 Antigravity(`.agents/`), Cursor(`.cursor/rules/`), Copilot(`.github/`), Claude(`CLAUDE.md`) 설정을 중앙 집중 관리합니다.

2. 🚀 **Zero-Cost & 초고속 결정론적(Deterministic) 코어**  
   복잡한 AST 파싱 및 코드 생성이 외부 LLM API 호출 없이 **로컬에서 10ms 이내(무과금)**로 즉시 실행됩니다.

3. 🧩 **지능형 토큰 예산 & JIT Skill 자동 분할 (Context Budgeting)**  
   Cursor 6,000자, Antigravity 12,000자 등 플랫폼별 용량 한계를 자동 감지하고, 대용량 지침은 온디맨드 JIT `SKILL.md`로 분해하여 **상시 토큰 소모를 80% 이상 절감**합니다.

4. 🔄 **1초 만에 프로젝트 마이그레이션 (`thop convert`)**  
   GitHub Copilot(`.github/`)을 사용하던 기존 프로젝트를 단 한 줄의 명령어로 Google Antigravity(`.agents/`) 구조로 완벽 변환합니다.

5. 🤖 **선택적 생성형 AI 증강 지원 (Hybrid AI Augmentation)**  
   필요에 따라 Claude, Gemini, OpenAI 및 **로컬 무과금 LLM (Ollama)**을 선택적으로 연동하여 의미론적 규칙 분할, 시맨틱 트리거 자동 생성, 3-Way 병합 충돌 해결을 수행할 수 있습니다.

---

## ⌨️ 빠른 시작 (Quick Start)

### 1. 설치 (Installation)

단 몇 초 만에 설치하여 바로 사용할 수 있습니다:

```bash
# Go를 통한 설치 (권장)
go install github.com/yunkon-kim/token-hop/cmd/token-hop@latest

# 또는 설치 스크립트로 바이너리 직접 다운로드 (Linux / macOS)
curl -fsSL https://raw.githubusercontent.com/yunkon-kim/token-hop/main/install.sh | bash

# 초고속 입력을 위한 단축어 설정 (선택 사항)
alias thop=token-hop
```

### 2. 바로 사용하기 (Instant Usage)

어떤 프로젝트 디렉터리에서든 즉시 실행:

```bash
# 1) 기존 GitHub Copilot(.github/) 설정을 Google Antigravity(.agents/)로 1초 만에 변환
thop convert --from copilot --to antigravity

# 2) 프로젝트 내 모든 규칙의 토큰 예산 및 문자수 정밀 감사
thop audit

# 3) SSOT(AGENTS.md)로부터 모든 AI 타겟(Antigravity, Cursor, Copilot)으로 동시 컴파일
thop compile

# 4) 신규 프로젝트에 SSOT 템플릿 초기화
thop init
```

---

### 🛠️ 소스코드 직접 빌드 (기여자용)

```bash
git clone https://github.com/yunkon-kim/token-hop.git
cd token-hop
make build && make test
```

---

## 📊 변환 매핑 매트릭스 (Mapping Matrix)

| 엔티티 (Entity) | UA-IR 정규화 역할 | Google Antigravity | Cursor AI | GitHub Copilot | Claude Code |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Rule** | 파일 경로/정적 제약 | `.agents/rules/*.md` | `.cursor/rules/*.mdc` | `.github/instructions/*.md` | `./CLAUDE.md` |
| **Skill** | 온디맨드 JIT 동적 절차 | `.agent/skills/<name>/` | `.cursor/rules/*.mdc` | `.github/skills/<name>/` | `.claude/skills/` |
| **Workflow** | 다단계 절차 및 검증 루프 | `.agents/workflows/*.md` | `@rule` 체이닝 문서 | `.github/prompts/*.md` | CLI 프롬프트 체인 |
| **Instruction** | 전역 시스템 기본 지침 | `AGENTS.md` | `.cursor/rules/base.mdc`| `.github/copilot-instructions.md` | `CLAUDE.md` |

---

## 📚 상세 문서 (Documentation)

- [전체 문서 목차 (docs/README.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/README.md)
- [시스템 아키텍처 및 구현 계획 명세서 (docs/design/token-hop-design-plan.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/docs/design/token-hop-design-plan.md)
- [테스트 스위트 및 픽스처 안내 (test/README.md)](file:///home/ubuntu/dev/yunkon-kim/token-hop/test/README.md)

---

## 📄 라이선스 (License)

[MIT License](file:///home/ubuntu/dev/yunkon-kim/token-hop/LICENSE) © 2026 Yunkon Kim
