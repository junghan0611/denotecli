# denotecli

**Denote knowledge base CLI for AI agents — search, read, analyze, and manage 3,000+ org-mode notes**

> Go stdlib only. Single binary. JSON output. Korean-native.

> **AI Agent Skill**: 에이전트용 스킬 문서는 [agent-config](https://github.com/junghan0611/agent-config) 리포의 `skills/denotecli/SKILL.md`에서 관리합니다.

---

## What This Does

CLI that gives AI agents structured access to a [Denote](https://protesilaos.com/emacs/denote)/org-mode knowledge base.
Not just search — semantic navigation, Korean↔English bridging, tag governance, and knowledge graph traversal.

```bash
./run.sh build              # Build + install
./run.sh test               # 105 tests
./run.sh showcase           # Visual search pattern check
./run.sh cover              # Coverage report
```

---

## Denote File Format

```
YYYYMMDDTHHMMSS[==SIGNATURE]--title-with-hyphens[__tag1_tag2].org
```

- **ID** = unique timestamp identifier (the key for everything)
- **Signature** = optional alphanumeric code (`==5a2`, `==0za`), used for [Denote signatures](https://protesilaos.com/emacs/denote#h:4e9c7512-84dc-4dfb-9fa9-e15d51178e5d) (e.g. syntopicon/propaedia ordering)
- **Frontmatter**: `#+title:`, `#+date:`, `#+filetags:`, `#+identifier:`
- **Links**: `[[denote:YYYYMMDDTHHMMSS]]`

Examples:
```
20251107T082610--제목-하이픈-구분__tag1_tag2_tag3.org          # standard
20250904T075937==5a2--힣-ai-에이전트__agents_ai.org           # with signature
20250421T125513==0--†-syntopicon-신토피콘__metameta.org       # single-char signature
20251021T105353--simple-title.org                              # no tags
```

---

## Commands (12)

### Search

```bash
denotecli search "에릭 호퍼" --tags bib --max 5       # title/tag/ID search
denotecli search-headings "양자역학" --level 1 --max 10 # heading search across all files
denotecli search-content "LSP 설정" --tags emacs       # full-text grep with tag filter
```

### Read

```bash
denotecli read 20250314T125213 --outline --level 2     # heading structure (TOC)
denotecli read 20250314T125213 --offset 41 --limit 20  # specific section by line range
```

### Navigate

```bash
denotecli graph 20250314T125213                         # outgoing + incoming links
denotecli keyword-map "이맥스"                          # Korean↔English tag mapping
denotecli keyword-map "emacs"                           # bidirectional
```

### Day / Timeline

```bash
denotecli day 2023-02-22                               # 특정 날짜 저널/노트/datetree 통합 (생성+수정 노트)
denotecli day --years-ago 3                             # N년 전 오늘
denotecli timeline-journal --month 2023-02             # 월간 저널 활동 개요
```

`denotecli day`는 다섯 축을 함께 반환합니다.

| Field | Source | Meaning |
|---|---|---|
| `journal` | `journal/<date>__journal.org` 또는 weekly | 그날 저널 시간 엔트리 |
| `datetree` | `*--diary.org` reverse datetree | 그날 datetree 엔트리 + CLOCK |
| `notes_created` | 파일명 ID prefix == 그날 | 그날 만든 Denote 노트 |
| `notes_modified` | 본문의 `#+hugo_lastmod:` == 그날 | 그날 **수정된** 노트. 생성일은 다른 날 |
| `years_ago` | system time | 같은 월일 N년 전인 경우 N |

`notes_created`와 `notes_modified`는 **상보(complementary)** 관계로 정의됩니다 — 같은 날 만들고 그날 `hugo_lastmod`까지 박은 파일은 `notes_created`에만 들어가고 `notes_modified`에서는 제외됩니다 (중복 방지). `notes_modified`는 명시적인 `#+hugo_lastmod:` 만 신뢰합니다 (mtime/`#+date:` fallback 없음). 다양한 org 타임스탬프 포맷 — `[2025-06-10]`, `[2025-03-29 Sat 02:06]`, `<2024-01-03 Wed 16:54>`, `2023-06-19`, `Time-stamp: <...>` — 모두 첫 `YYYY-MM-DD` 패턴으로 정규화됩니다.

`notes_modified` 패스는 모든 노트의 frontmatter head(첫 50줄)만 읽기 때문에 ~3,300 노트 코퍼스에서 100ms 수준을 유지합니다.

### Manage

```bash
denotecli create --title "대화 기록" --tags llmlog,emacs --dir ~/org/llmlog
denotecli rename-tag --from llms --to llm --dry-run     # batch tag rename (preview)
denotecli rename-tag --from llms --to llm               # actual rename (filename + frontmatter)
denotecli tags --top 20                                  # tag statistics
denotecli tags --suggest                                 # stem-based duplicate detection
```

---

## Performance

| Command | Scope | Time |
|---------|-------|------|
| `search` | 3K files, filenames | ~16ms |
| `search-headings` | 3K files, 60K headings | ~30ms |
| `search-content` | 3K files, 14MB text | ~270ms |
| `keyword-map` | meta notes | ~24ms |
| `graph` | 3K files, backlink scan | ~85ms |
| `tags --suggest` | 2K+ tags, Porter stemmer | ~23ms |

---

## Install

```bash
git clone https://github.com/junghan0611/denotecli.git
cd denotecli
./run.sh build    # → ~/.local/bin/denotecli
```

Requires Go 1.21+. No external dependencies (stdlib only).

---

## Output

All output is JSON. Signature field included when present (`"signature": "5a2"`), omitted when absent.

```json
{
  "id": "20250904T075937",
  "signature": "5a2",
  "title": "힣-ai-에이전트-편재성-기억-연결",
  "tags": ["agents", "ai"],
  "date": "2025-09-04",
  "path": "/home/.../meta/20250904T075937==5a2--힣-ai-에이전트-편재성-기억-연결__agents_ai.org"
}
```

---

## Testing

```bash
./run.sh test       # All 105 tests
./run.sh showcase   # Visual: all search patterns with input→output
./run.sh cover      # Coverage report (logic functions: 85-100%)
```

**Showcase tests** (`go test -v -run TestShowcase`) print every search pattern with actual results — makes edge cases visible at a glance.

---

## Project Structure

```
denotecli/
├── run.sh                    # Build, test, showcase, cover
├── README.md
├── docs/
│   └── obsidian-cli-comparison.md
└── denotecli/                # Go source (single package, 15 modules + 16 test files)
    ├── main.go               # CLI routing + flag parsing
    ├── parser.go             # Denote filename (with signature) + frontmatter + link parser
    ├── search.go             # Directory scanner + title/tag search
    ├── search_headings.go    # Heading search across all files
    ├── search_content.go     # Full-text content search
    ├── read.go               # Read by ID + outline extraction
    ├── graph.go              # Outgoing/incoming link traversal
    ├── keyword_map.go        # Korean↔English keyword mapping
    ├── create.go             # Note creation with Denote naming
    ├── rename_tag.go         # Batch tag rename (filename + frontmatter)
    ├── tags.go               # Tag statistics
    ├── tag_suggest.go        # Stem-based duplicate detection
    ├── day.go                # Day query (journal + datetree + notes)
    ├── timeline_journal.go   # Monthly journal timeline
    ├── stemmer.go            # Porter stemmer (embedded, no deps)
    └── *_test.go             # 105 tests (unit + integration + showcase)
```

---

## Related Projects

| Project | Description |
|---------|-------------|
| [agent-config](https://github.com/junghan0611/agent-config) | AI agent configuration — denotecli SKILL.md managed here |
| [zotero-config](https://github.com/junghan0611/zotero-config) | Headless Zotero + **bibcli** (sister CLI for 8K+ bibliography) |

---

**Author**: [@junghanacs](https://github.com/junghan0611) · **License**: Apache 2.0
