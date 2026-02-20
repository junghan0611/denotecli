---
name: denote-org
description: "Search, read, and analyze 3,000+ Denote/org-mode files. Use when working with ~/org/, Denote files (YYYYMMDDTHHMMSS--title__tags.org), org-mode knowledge bases, or when user asks about notes, journal entries, or bibliography."
---

# denotecli - Denote Knowledge Base CLI

Search, read, and analyze notes from the user's Denote/org-mode knowledge base (2,800+ files across notes, bib, journal, llmlog).

## Prerequisites

Binary must be in PATH. Build from source:

```bash
# Clone (if not already)
git clone https://github.com/junghan0611/org-mode-skills.git
cd org-mode-skills

# Build + install to ~/.local/bin
./run.sh build
```

Requires Go 1.21+. No external dependencies (stdlib only).

## Commands

### Search notes

```bash
denotecli search "에릭 호퍼" --dirs ~/org --max 5
denotecli search "emacs" --dirs ~/org --tags emacs
denotecli search "창조" --dirs ~/org --title-only
```

- Multiple words = AND condition (all must match)
- Searches: Denote ID, title (from filename), tags
- Case-insensitive (Korean included)
- `--tags`: filter by tag (comma-separated, OR)
- `--title-only`: only search title field

### Read a note

```bash
denotecli read 20250314T152111 --dirs ~/org
denotecli read 20241206T085900 --dirs ~/org --limit 50
```

Returns full content + parsed frontmatter metadata + outgoing denote links.

### Tag statistics

```bash
denotecli tags --dirs ~/org --top 20
denotecli tags --dirs ~/org --pattern "emacs|vim"
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--dirs DIR,...` | Search directories (comma-separated) | `~/org` |
| `--tags TAG` | Filter by tag (comma-separated) | all |
| `--title-only` | Search title only | false |
| `--max N` | Max results | search: 20 |
| `--offset N` | Start line (read) | 0 |
| `--limit N` | Lines to read (read, 0=all) | 0 |
| `--pattern PAT` | Tag regex filter (tags) | all |
| `--top N` | Top N tags | 50 |

## Output

All output is JSON. Examples:

**search** returns brief entries:
```json
[{"id": "20251107T082610", "title": "제목", "tags": ["tag1", "tag2"], "date": "2025-11-07", "path": "/home/..."}]
```

**read** returns full content:
```json
{"id": "...", "title": "...", "tags": [...], "date": "...", "path": "...", "content": "...", "links": ["20240601T204208"]}
```

**tags** returns statistics:
```json
{"total_files": 2839, "total_tags": 2162, "tags": [{"name": "bib", "count": 966}, ...]}
```

## Denote File Format

- **Filename**: `YYYYMMDDTHHMMSS--title-with-hyphens__tag1_tag2.org`
- **ID** = unique timestamp identifier (the key for everything)
- **Frontmatter**: `#+title:`, `#+date:`, `#+filetags:`, `#+identifier:`
- **Links**: `[[denote:YYYYMMDDTHHMMSS]]`

## Environment Paths

The knowledge base root differs by environment. Use `--dirs` accordingly:

| Environment | Root Path | Example |
|-------------|-----------|---------|
| **Local** (Claude Code) | `~/org` | `denotecli search "query" --dirs ~/org` |
| **Container** (OpenClaw) | `/data/org` | `denotecli search "query" --dirs /data/org` |

Multiple directories (comma-separated):
```bash
denotecli search "query" --dirs /data/org/notes,/data/org/bib,/data/org/journal,/data/org/llmlog
```

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
