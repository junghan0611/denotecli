# denotecli Design

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Go CLI that parses Denote filenames and org frontmatter from ~/org/ (3,000+ files), providing search/read/tags for AI agents (JSON-only output).

**Architecture:** Single Go module at `denotecli/`, stdlib only. Regex-based Denote filename parser, filepath.WalkDir for file discovery. Subcommands route via os.Args.

**Tech Stack:** Go 1.21+, stdlib only (encoding/json, regexp, os, path/filepath, strings, fmt, bufio, sort, strconv)

**Targets:**
- Local: Claude Code, OpenClaw agents
- Paths: `~/org/` (local), `/data/org/` (container) — configurable via `--dirs`

---

## Repository Structure

```
orgmode-skills/
├── SKILL.md                     # pi-skills compatible skill definition
├── README.md                    # concise project description
├── .gitignore
├── denotecli/
│   ├── go.mod
│   ├── main.go                  # CLI routing (search/read/tags)
│   ├── parser.go                # Denote filename + frontmatter parsing
│   ├── parser_test.go
│   ├── search.go                # Search logic (filepath.WalkDir)
│   ├── search_test.go
│   ├── read.go                  # File read (ID → content + links)
│   ├── tags.go                  # Tag extraction/statistics
│   └── tags_test.go
├── docs/
│   ├── plans/                   # Design docs
│   └── archive/                 # Legacy Python scripts/docs
```

## Commands

### `denotecli search <query> [flags]`

Search notes by filename metadata and frontmatter.

| Flag | Default | Description |
|------|---------|-------------|
| `<query>` | required | Search term (Korean/English, multi-word AND) |
| `--tags` | "" | Tag filter (comma-separated) |
| `--dirs` | `~/org` | Search directories (comma-separated) |
| `--title-only` | false | Search title only |
| `--max` | 20 | Max results |

**Search targets:** filename (ID, title, tags) + frontmatter (title, filetags, description).
Content search: prefix query with `content:` to grep file body.

**Output:**
```json
[
  {
    "id": "20251107T082610",
    "title": "제목-하이픈-구분",
    "tags": ["tag1", "tag2"],
    "date": "2025-11-07",
    "path": "/home/user/org/notes/20251107T082610--제목__tag1_tag2.org"
  }
]
```

### `denotecli read <id> [flags]`

Read a note by Denote ID.

| Flag | Default | Description |
|------|---------|-------------|
| `<id>` | required | Denote ID (e.g., 20251107T082610) |
| `--offset` | 0 | Start line |
| `--limit` | 0 (all) | Lines to read |
| `--dirs` | `~/org` | Search directories |

**Output:**
```json
{
  "id": "20251107T082610",
  "title": "제목",
  "tags": ["tag1", "tag2"],
  "date": "[2025-11-07 Fri 08:26]",
  "path": "/path/to/file.org",
  "content": "full file content...",
  "links": ["20240601T204208", "20250314T152111"]
}
```

### `denotecli tags [flags]`

Tag statistics across all Denote files.

| Flag | Default | Description |
|------|---------|-------------|
| `--pattern` | "" | Tag filter (regex) |
| `--top` | 50 | Top N tags |
| `--dirs` | `~/org` | Search directories |

**Output:**
```json
{
  "total_files": 3100,
  "total_tags": 450,
  "tags": [
    {"name": "emacs", "count": 280},
    {"name": "programming", "count": 215}
  ]
}
```

## Core Data Structures

```go
// DenoteFile represents parsed Denote file metadata.
type DenoteFile struct {
    ID    string   `json:"id"`
    Title string   `json:"title"`
    Tags  []string `json:"tags"`
    Date  string   `json:"date"`
    Path  string   `json:"path"`
}

// DenoteContent extends DenoteFile with file content.
type DenoteContent struct {
    DenoteFile
    Content string   `json:"content"`
    Links   []string `json:"links"`
}
```

## Filename Parsing

```go
var denoteRe = regexp.MustCompile(`^(\d{8}T\d{6})--(.+?)(?:__(.+))?\.org$`)
// group 1: ID, group 2: title (hyphen-separated), group 3: tags (underscore-separated)
```

## Implementation Tasks (TDD)

### Task 1: Go module + Denote filename parser
- Create: `denotecli/go.mod`, `parser.go`, `parser_test.go`
- Test: Korean titles, no-tag files, multiple hyphens

### Task 2: Search logic
- Create: `denotecli/search.go`, `search_test.go`
- Test: multi-word AND, tag filter, title-only, max limit

### Task 3: Read + Tags
- Create: `denotecli/read.go`, `tags.go`, `tags_test.go`
- Test: ID lookup, frontmatter parsing, link extraction, tag aggregation

### Task 4: CLI main + JSON output
- Create: `denotecli/main.go`
- Manual test with real ~/org/ data

### Task 5: Repo cleanup + SKILL.md + benchmark
- Move legacy files to docs/archive/
- Write SKILL.md (pi-skills format)
- Write README.md (concise)
- Benchmark with 3,000+ files (target: <100ms for search)
