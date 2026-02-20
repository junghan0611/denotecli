# Phase 2 Design — graph + keyword-map

> denotecli v0.3.0 target

## Scope

에이전트 B(OpenClaw) 요구사항 중 미구현 2개 커맨드.
간단하게 구현 → 에이전트가 바로 사용 → 이후 확장.

## 1. `graph` — 링크 탐색

```
denotecli graph <id> [--depth N] [--direction outgoing|incoming|both] [--dirs DIR,...]
```

### 동작

- **outgoing**: 대상 파일 read → `ExtractLinks()` (기존 코드 재사용)
- **incoming**: `rg "denote:<id>" --include="*.org" -l` 실행 → 결과 ParseFilename
- **both**: 위 둘 다
- **depth>1**: 1단계 결과 ID들로 재귀 (visited set 순환 방지, max depth 제한)
- **rg 없으면**: pure Go fallback (ScanDirs + 전체 파일 읽기)

### 출력

```json
{
  "id": "20250314T152111",
  "title": "시작 노트",
  "direction": "both",
  "depth": 1,
  "outgoing": [{"id": "...", "title": "...", "tags": [...], "date": "...", "path": "..."}],
  "incoming": [{"id": "...", "title": "...", "tags": [...], "date": "...", "path": "..."}]
}
```

### 구현

- 새 파일: `graph.go`, `graph_test.go`
- main.go에 `case "graph": cmdGraph()` 추가

## 2. `keyword-map` — 한↔영 키워드 매핑

```
denotecli keyword-map <query> [--dirs DIR,...] [--map-file PATH]
```

### 동작

1. `~/org/meta/` 에서 `*keyword-map*` 파일 탐색 (또는 `--map-file` 직접 지정)
2. org heading(`* english`) → 하위 리스트(`- 한글1, 한글2`) 파싱
3. query가 한글이면 → 매핑된 영어 키워드 반환 (ko_to_en)
4. query가 영어이면 → 매핑된 한글 키워드 반환 (en_to_ko)

### 매핑 파일 포맷 (~/org/meta/TIMESTAMP--denote-keyword-map__meta_denotecli.org)

```org
#+title: Denote Keyword Map
#+filetags: :meta:denotecli:

* philosophy
- 철학, 사상, 세계관, 인생관, 가치관

* creativity
- 창의성, 창조성, 발상, 아이디어, 독창성

* knowledge
- 지식, 앎, 인식, 지혜, 학문
```

### 출력

```json
{
  "query": "창조",
  "mappings": ["creativity", "innovation", "imagination"],
  "direction": "ko_to_en"
}
```

### 구현

- 새 파일: `keywordmap.go`, `keywordmap_test.go`
- main.go에 `case "keyword-map": cmdKeywordMap()` 추가

## Design Decisions

| 결정 | 이유 |
|---|---|
| rg 외부 호출 + fallback | incoming 링크 3,000+ 파일 grep에서 10x+ 성능 차이 |
| 매번 전체 스캔 | Denote 철학 "파일시스템=DB", 캐시 복잡도 회피 |
| 매핑 데이터 ~/org/ 내 | Denote 생태계 안에서 관리, Emacs에서 편집 가능 |
| depth 기본값 1 | 대부분 사용 시나리오가 직접 연결 |

## SKILL.md 업데이트

graph, keyword-map 커맨드 사용법 추가.

## br Issues (prefix: dc)

- dc-279: read --outline
- dc-2jo: graph
- dc-28y: keyword-map
- dc-2u6: 통합 테스트
