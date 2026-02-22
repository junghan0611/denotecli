# denotecli vs obsidian-cli 커버리지 비교

> 2026-02-22 작성. obsidian-cli SKILL.md 기준 대조.

## 설계 차이

| | obsidian-cli | denotecli |
|---|---|---|
| **데이터** | Markdown + DB (Obsidian 인덱스) | Org-mode 플레인 텍스트 (DB 없음) |
| **볼트** | 다중 볼트, config에서 resolve | 단일 경로 `~/org` (--dirs) |
| **ID 체계** | 파일명/경로 | Denote timestamp ID (`YYYYMMDDTHHMMSS`) |
| **링크** | `[[wikilink]]` | `[[denote:ID]]` |
| **에이전트 출력** | 미상 (텍스트?) | JSON 고정 |

## 기능 매핑

### ✅ 커버됨

| obsidian-cli | denotecli | 비고 |
|---|---|---|
| `search "query"` | `search "query"` | denotecli이 ID+태그까지 검색. 우위 |
| `print-default --path-only` | `--dirs ~/org` | 볼트 개념 불필요 |
| 직접 .md 편집 | 직접 .org 편집 | 동일 |

### ❌ 없음 → 추가 예정

| obsidian-cli | 추가 방향 | 구현 |
|---|---|---|
| `search-content "query"` | `search-content "query"` | rg 래핑. DB 없으니 rg가 최적 |
| `create "path" --content` | `create --title --tags --content` | Denote 네이밍+헤더 자동생성 |

### ⬇️ 낮은 우선순위 (추가 안 함)

| obsidian-cli | 이유 |
|---|---|
| `move "old" "new"` | Emacs denote-rename이 링크 리팩터까지 처리 |
| `delete "path"` | `rm`이면 됨 |
| `set-default` / vault 관리 | 단일 경로 체계, 불필요 |

### 🔮 denotecli 고유 (obsidian-cli에 없음)

| 기능 | 상태 | 에이전트 가치 |
|---|---|---|
| `read ID` (메타+링크 파싱) | ✅ 운용중 | 핵심 — ID로 정확히 접근 |
| `tags --top N` | ✅ 운용중 | 지식베이스 조감도 |
| `read --outline` | 🔧 구현 예정 | **핵심** — 에이전트가 head 100 대신 헤딩 구조로 판단 |
| `search --headings` | 🔧 구현 예정 | org 헤딩 = 문서의 의미 단위. 본문보다 헤딩 검색이 정확 |
| `graph` | 🔧 구현 예정 | outgoing/incoming 링크 탐색 |
| `keyword-map` | 🔧 구현 예정 | 한영 양방향 키워드 매핑 |
| `create` | 🔧 구현 예정 | llmlog 자동생성 |

## 에이전트 워크플로우 (목표)

```
1. tags          → 조감도: 어떤 주제가 많은가
2. keyword-map   → 한글 "양자역학" → 영어 태그 "quantum" 연결
3. search        → 제목/태그로 후보 노트 찾기
4. search-content→ rg 기반 본문 grep (스니펫+라인)
5. read --outline→ 긴 문서의 헤딩 구조 파악
6. read          → 필요한 섹션만 offset/limit로 읽기
7. graph         → 관련 노트 네트워크 탐색
```

## 구현 우선순위

| 순위 | 기능 | 이유 |
|---|---|---|
| 1 | `read --outline` | 에이전트 효율 극대화. org 헤딩이 곧 목차 |
| 2 | `search-content` | obsidian-cli 정렬 + rg 래핑으로 간단 |
| 3 | `create` | llmlog 자동화 |
| 4 | `keyword-map` | 온톨로지 기초 |
| 5 | `graph` | 링크 네트워크 |
