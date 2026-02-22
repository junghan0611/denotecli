---
name: denotecli
description: "Search, read, and analyze 3,000+ Denote/org-mode notes. Supports title/tag search, heading search across all files, outline extraction, and full content reading. Use when working with ~/org/, Denote files, org-mode knowledge bases, or when user asks about notes, journal entries, or bibliography."
---

# denotecli — Denote Knowledge Base CLI

Search, read, and analyze 3,000+ Denote/org-mode notes (notes, bib, journal, llmlog).

Binary is bundled in the skill directory. Invoke via `{baseDir}/denotecli`.

All output is JSON.

## Typical Workflow

```
1. search "에릭 호퍼"              → find notes by title/tag (fast, filename-only)
2. search-headings "창조"          → find topics inside notes (scans all headings)
3. read <ID> --outline --level 2   → see document structure before reading
4. read <ID> --offset 41 --limit 20 → read specific section by line range
```

## Commands

### search — find notes by title, tag, ID

```bash
{baseDir}/denotecli search "에릭 호퍼" --dirs ~/org --max 5
{baseDir}/denotecli search "emacs" --dirs ~/org --tags emacs
{baseDir}/denotecli search "창조" --dirs ~/org --title-only
```

- Multiple words = AND (all must match)
- Searches: Denote ID, title (from filename), tags
- Case-insensitive (Korean included)
- `--tags TAG`: filter by tag (comma-separated, OR)
- `--title-only`: search title field only

```json
[{"id": "20251107T082610", "title": "제목", "tags": ["tag1", "tag2"], "date": "2025-11-07", "path": "/home/..."}]
```

### search-headings — find topics inside notes

```bash
{baseDir}/denotecli search-headings "양자역학" --dirs ~/org --max 10
{baseDir}/denotecli search-headings "창조" --dirs ~/org --level 1 --max 5
```

- Searches org headings (`* heading`) across ALL files (~3K files, ~60K headings, ~30ms)
- Returns file metadata + matched heading with line number
- `--level N`: only search headings up to level N (0=all)

```json
[{"id": "...", "title": "...", "tags": [...], "path": "...", "heading": {"level": 1, "title": "양자역학의 해석", "line": 23}}]
```

### read — read note content

```bash
{baseDir}/denotecli read 20250314T152111 --dirs ~/org
{baseDir}/denotecli read 20241206T085900 --dirs ~/org --offset 40 --limit 30
```

- Returns full content + parsed frontmatter + outgoing `[[denote:ID]]` links
- Use `--offset`/`--limit` to read specific line ranges (from outline)

```json
{"id": "...", "title": "...", "tags": [...], "date": "...", "path": "...", "content": "...", "links": ["20240601T204208"]}
```

### read --outline — see document structure

```bash
{baseDir}/denotecli read 20250314T152111 --dirs ~/org --outline
{baseDir}/denotecli read 20250314T152111 --dirs ~/org --outline --level 2
```

- Returns org heading structure: level, title, line number, org tags
- Use before full read — line numbers let you target `--offset`/`--limit` precisely
- `--level N`: filter headings up to level N (0=all)

```json
{"id": "...", "title": "...", "tags": [...], "outline": [{"level": 1, "title": "1장 서론", "line": 5}, {"level": 2, "title": "1.1 배경", "line": 7}], "links": [...]}
```

### tags — knowledge base overview

```bash
{baseDir}/denotecli tags --dirs ~/org --top 20
{baseDir}/denotecli tags --dirs ~/org --pattern "emacs|vim"
```

```json
{"total_files": 3156, "total_tags": 2162, "tags": [{"name": "bib", "count": 966}, ...]}
```

## Flags

| Flag | Applies to | Description | Default |
|------|-----------|-------------|---------|
| `--dirs DIR,...` | all | Search directories (comma-separated) | `~/org` |
| `--max N` | search, search-headings | Max results | 20 |
| `--tags TAG` | search | Filter by tag (comma-separated, OR) | all |
| `--title-only` | search | Search title field only | false |
| `--level N` | search-headings, read --outline | Max heading level (0=all) | 0 |
| `--outline` | read | Show heading structure instead of content | false |
| `--offset N` | read | Start line (1-indexed from outline) | 0 |
| `--limit N` | read | Lines to read (0=all) | 0 |
| `--pattern PAT` | tags | Tag name regex filter | all |
| `--top N` | tags | Top N tags | 50 |

## Denote File Format

- **Filename**: `YYYYMMDDTHHMMSS--title-with-hyphens__tag1_tag2.org`
- **ID** = unique timestamp identifier (the key for everything)
- **Frontmatter**: `#+title:`, `#+date:`, `#+filetags:`, `#+identifier:`
- **Links**: `[[denote:YYYYMMDDTHHMMSS]]`

## Knowledge Base Structure

| Directory | Purpose | Scale |
|-----------|---------|-------|
| `notes/` | Main notes | 800+ |
| `bib/` | Bibliography | 900+ |
| `journal/` | Weekly journals | 700+ |
| `llmlog/` | LLM conversation logs | 300+ |
| `meta/` | Meta topics | - |
| `archives/` | Archived notes | - |
| root `.org` files | diary, tasks, etc. | ~10 |

## Environment Paths

| Environment | Root Path |
|-------------|-----------|
| **Local** (Claude Code) | `~/org` |
| **Container** (OpenClaw) | `~/org` |

Multiple directories: `--dirs ~/org/notes,~/org/bib,~/org/journal,~/org/llmlog`
