# denotecli

**Denote knowledge base CLI for AI agents: search, read, and analyze 3,000+ org-mode notes**

> Go stdlib only. Single binary. JSON output. Korean-native.

[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)

---

## What This Does

CLI tool that gives AI agents (Claude Code, OpenClaw, etc.) structured access to a [Denote](https://protesilaos.com/emacs/denote)/org-mode knowledge base. Parses filenames, frontmatter, and inter-note links. Returns JSON.

```
./run.sh build              # Build + install to ~/.local/bin
denotecli search "에릭 호퍼"  # Search 3,000+ notes (~16ms)
denotecli read 20250314T152111  # Read full content + metadata + links
denotecli tags --top 20     # Tag statistics across all files
```

---

## Install

```bash
git clone https://github.com/junghan0611/denotecli.git
cd denotecli
./run.sh build    # Builds binary → ~/.local/bin/denotecli
```

Requires Go 1.21+. No external dependencies (stdlib only).

---

## Commands

### search

```bash
denotecli search "에릭 호퍼" --dirs ~/org --max 5
denotecli search "emacs" --dirs ~/org --tags emacs
denotecli search "창조" --dirs ~/org --title-only
```

- Multiple words = AND (all must match)
- Searches: Denote ID, title (from filename), tags
- Case-insensitive (Korean included)

### read

```bash
denotecli read 20250314T152111 --dirs ~/org
denotecli read 20241206T085900 --dirs ~/org --limit 50
denotecli read 20250314T152111 --dirs ~/org --outline
```

Returns full content + parsed frontmatter + outgoing `[[denote:ID]]` links.
`--outline` returns only the heading structure (level, title, line number) — use this first to understand document structure before reading specific sections.

### tags

```bash
denotecli tags --dirs ~/org --top 20
denotecli tags --dirs ~/org --pattern "emacs|vim"
```

---

## Output

All output is JSON:

```json
// search
[{"id": "20251107T082610", "title": "제목", "tags": ["tag1"], "date": "2025-11-07", "path": "..."}]

// read
{"id": "...", "title": "...", "content": "...", "links": ["20240601T204208"]}

// tags
{"total_files": 2839, "total_tags": 2162, "tags": [{"name": "bib", "count": 966}]}
```

---

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--dirs DIR,...` | Search directories (comma-separated) | `~/org` |
| `--tags TAG` | Filter by tag (comma-separated, OR) | all |
| `--title-only` | Search title field only | false |
| `--max N` | Max search results | 20 |
| `--outline` | Show heading structure only | false |
| `--offset N` | Start line for read | 0 |
| `--limit N` | Lines to read (0=all) | 0 |
| `--pattern PAT` | Tag name regex filter | all |
| `--top N` | Top N tags | 50 |

---

## As an AI Skill

### Claude Code

```bash
# Option 1: Symlink as skill
ln -s /path/to/denotecli ~/.claude/skills/denote-org

# Option 2: Use via pi-skills (pre-installed)
# See https://github.com/junghan0611/pi-skills
```

### OpenClaw / Container

```bash
denotecli search "query" --dirs /data/org
```

See [SKILL.md](SKILL.md) for full skill documentation.

---

## Denote File Format

```
YYYYMMDDTHHMMSS--title-with-hyphens__tag1_tag2.org
```

- **ID** = timestamp (`20250314T152111`) — the unique key for everything
- **Frontmatter**: `#+title:`, `#+date:`, `#+filetags:`, `#+identifier:`
- **Links**: `[[denote:YYYYMMDDTHHMMSS]]`

---

## Project Structure

```
denotecli/
├── run.sh                 # Build + install entry point
├── SKILL.md               # AI skill definition (pi-skills compatible)
├── denotecli/
│   ├── main.go            # CLI routing (search/read/tags)
│   ├── parser.go          # Denote filename + frontmatter + link parser
│   ├── search.go          # Directory scanner + search engine
│   ├── read.go            # Read by ID with frontmatter enrichment
│   ├── tags.go            # Tag aggregation + statistics
│   ├── parser_test.go     # Parser tests (9 cases incl. U+00A0)
│   ├── search_test.go     # Search tests (7 cases)
│   └── tags_test.go       # Tags + read tests (4 cases)
└── docs/                  # Design docs + archive
```

---

## Related Projects

| Project | Description |
|---------|-------------|
| [pi-skills](https://github.com/junghan0611/pi-skills) | AI skill collection for Claude Code — denotecli is distributed here as a skill |
| [zotero-config](https://github.com/junghan0611/zotero-config) | Headless Zotero-to-BibTeX workflow + **bibcli** (sister CLI for 8,000+ bibliography entries) |

---

## Links

- **Digital Garden**: [notes.junghanacs.com](https://notes.junghanacs.com)
- **Denote**: [protesilaos.com/emacs/denote](https://protesilaos.com/emacs/denote)

---

**Author**: [@junghanacs](https://github.com/junghan0611)

## License

Apache 2.0
