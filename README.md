<div align="center">

# Ksecret

[![CI](https://github.com/jeonjw85/Ksecret/actions/workflows/ci.yml/badge.svg)](https://github.com/jeonjw85/Ksecret/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev/)

고유식별정보와 국내/글로벌 서비스 키를 찾는 한국 맞춤 시크릿 스캐너

</div>

```text
$ ksecret scan testdata

ksecret: 이슈 2건

kr-rrn               심각     hits.env:2:8
  주민등록번호  9001...3459

kr-toss              심각     hits.env:11:6
  토스페이먼츠 키  test...i9j0
```

## 설치

### 일반 설치 (macOS / Linux)

```bash
curl -sSfL https://raw.githubusercontent.com/jeonjw85/Ksecret/main/scripts/install.sh | sh
```

> 설치 위치 기본값 : `/usr/local/bin`

### Go가 설치되어 있을 시 설치(권장)

```bash
go install github.com/jeonjw85/Ksecret@latest
```

> Windows는 위 `go install` 을 이용하세요 (install.sh 미지원)

## 사용법

> 주민번호,사업자번호,카드는 오탐을 줄이기 위해 정규식만 보지 않고 체크섬이 맞을 때만 보고됩니다

```bash
ksecret scan .                                       # 현재 디렉토리 스캔 (콘솔)
ksecret scan ./src --format json                     # 스캔 후 JSON 출력
ksecret scan . --format sarif --output results.sarif # SARIF 파일로 저장
ksecret install --pre-commit                         # pre-commit 훅 설치
```

| 종료 코드 | 의미            |
| --------- | --------------- |
| 0         | 이슈 없음       |
| 1         | 시크릿/PII 발견 |
| 2         | 실행 오류       |

### pre-commit

`ksecret install --pre-commit` 은 `.git/hooks/pre-commit` 에 staged 파일만 스캔하는 훅을 넣어 커밋할 때만 돌아갑니다

> 훅은 `#!/bin/sh` 스크립트입니다. Windows에서는 Git for Windows(Git Bash)가 있어야 동작합니다 — Git 기본 설치에 포함되어 있습니다.

직접 쓰려면:

```bash
ksecret scan --staged
```

### GitHub Action

```yaml
- uses: jeonjw85/Ksecret@v1
  with:
      path: .
      format: sarif
      fail: true
```

`format: sarif` 이면 워크스페이스에 `results.sarif` 를 남깁니다.  
GitHub code scanning에 올리려면 그 파일을 `github/codeql-action/upload-sarif` 로 업로드하면 됩니다.

이슈가 있어도 job을 통과시키려면 `fail: false` 로 변경하세요

## 탐지영역

- 고유식별정보: 주민등록번호, 외국인등록번호, 운전면허, 여권
- 사업자등록번호, 카드번호
- 국내 서비스 키: 카카오, 네이버, 토스, 쿨에스엠에스/솔라피, 포트원, 네이버 클라우드
- 글로벌 서비스 키: AWS, GCP, Slack, JWT

> 주민/외국인/사업자/카드 번호는 오탐 방지를 위해 정해진 체크섬 실패 시 보고되지 않습니다

## 설정

프로젝트 루트의 `.ksecret.yml` 을 읽습니다 (예시: `.ksecret.example.yml`)

```yaml
allowlist:
    paths: ["testdata/**", "**/*_test.go"]
    rules: ["jwt"]
    patterns: ["EXAMPLE"]
exclude:
    paths: ["vendor/**"]
max_file_bytes: 1048576
```

- `allowlist.paths` / `exclude.paths`: glob (`**` 지원)
- `allowlist.rules`: 해당 규칙 전부 무시
- `allowlist.patterns`: 그 문자열이 들어 있는 라인은 무시
- `max_file_bytes`: 기본 1MB. 바이너리(NUL)와 `.git`, `node_modules` 등은 제외

`--config` 로 다른 파일 지정 가능

## 예시 출력

```text
ksecret: 이슈 2건

kr-rrn               심각     hits.env:2:8
  주민등록번호  9001...3459

kr-toss              심각     hits.env:11:6
  토스페이먼츠 키  test...i9j0
```

JSON:

```json
{
    "findings": [
        {
            "rule_id": "kr-rrn",
            "severity": "critical",
            "file": "hits.env",
            "line": 2,
            "column": 8,
            "secret": "9001...3459",
            "message": "주민등록번호"
        }
    ],
    "counts": {
        "kr-rrn": 1
    }
}
```

시크릿은 앞 4글자 + `...` + 뒤 4글자로만 보여 줍니다.  
콘솔, JSON, SARIF 모두 같고, 터미널이나 CI 로그에 원문이 그대로 안 남게 하려는 겁니다. 8글자 이하는 맨 앞 한 글자만 남기고 나머지는 `...` 처리합니다.

## 개발

```bash
go test -race -shuffle=on -count=1 ./...
go build -trimpath -ldflags="-s -w" -o bin/ksecret .
```

## License

[MIT](LICENSE)
