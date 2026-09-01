# Token-Hop Test Suite & Fixtures

이 디렉터리는 `token-hop` (`thop`)의 단위 테스트 및 실 프로젝트 대상 변환 검증 픽스처(Fixture)를 포함합니다.

## 디렉터리 구성

- **`fixtures/cm-beetle/.github/`**:
  - 실제 오픈소스 멀티클라우드 프로젝트 [`cloud-barista/cm-beetle`](https://github.com/cloud-barista/cm-beetle)의 실제 GitHub Copilot 설정 파일들입니다.
  - `copilot-instructions.md` (전역 지침)
  - `instructions/` (6개 파일: `go`, `analyzer`, `markdown`, `tb-sync`, `transx`, `ui`)
  - `prompts/` (5개 파일: `api-guide`, `git-commit`, `release-staging`, `sync-tb`, `sync-tb-model`)
- **`cm_beetle_conversion_test.go`**:
  - `cm-beetle` Copilot 설정 파싱 $\rightarrow$ Google Antigravity(`.agents/rules/`, `.agents/workflows/`, `AGENTS.md`) 변환 $\rightarrow$ Cursor(`.cursor/rules/*.mdc`) 변환 $\rightarrow$ 토큰 수 및 컨텍스트 한도 감사(Token & Context Audit) 자동 검증 테스트.

## 테스트 실행 방법

```bash
# 전체 테스트 실행
go test -v ./...

# 특정 테스트만 실행
go test -v -run TestCmBeetleCopilotToAntigravityConversion ./test

# CLI 바이너리로 직접 변환 테스트
./bin/thop convert --from copilot --to antigravity --input test/fixtures/cm-beetle/.github --output /tmp/test-antigravity
./bin/thop audit --input test/fixtures/cm-beetle/.github
```
