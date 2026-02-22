# denotecli

**Denote knowledge base CLI for AI agents — search, read, analyze, and manage 3,000+ org-mode notes**

> Go stdlib only. Single binary. JSON output. Korean-native.

---

## What This Does

CLI that gives AI agents structured access to a [Denote](https://protesilaos.com/emacs/denote)/org-mode knowledge base.
Not just search — semantic navigation, Korean↔English bridging, tag governance, and knowledge graph traversal.

```bash
./run.sh build              # Build + install
./run.sh test               # 82+ tests
./run.sh showcase           # Visual search pattern check
./run.sh cover              # Coverage report
```

---

## Commands (10)

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

All output is JSON. See [SKILL.md](SKILL.md) for detailed examples per command.

---

## Testing

```bash
./run.sh test       # All 82+ tests
./run.sh showcase   # Visual: all search patterns with input→output
./run.sh cover      # Coverage report (logic functions: 85-100%)
```

**Showcase tests** (`go test -v -run TestShowcase`) print every search pattern with actual results — makes edge cases visible at a glance:

```
=== search patterns ===
  [한글 다중단어] query="에릭 호퍼" → 1건
    → 20250314T125213 에릭호퍼-방랑자의-철학 [autobiography,bib,philosophy]
  [영어 + 태그 필터] query="설정" tags="emacs" → 2건
    → 20250101T100000 emacs-설정-가이드 [config,emacs]
    → 20250201T100000 둠이맥스-설정 [config,doomemacs,emacs]
```

---

## Project Structure

```
denotecli/
├── run.sh                    # Build, test, showcase, cover
├── SKILL.md                  # AI skill definition (pi-skills compatible)
├── README.md
├── docs/
│   └── obsidian-cli-comparison.md
└── denotecli/                # Go source (single package)
    ├── main.go               # CLI routing + flag parsing (329 lines)
    ├── parser.go             # Denote filename + frontmatter + link parser
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
    ├── stemmer.go            # Porter stemmer (embedded, no deps)
    ├── *_test.go             # Unit tests per module
    ├── integration_test.go   # 7 agent workflow scenarios
    └── showcase_test.go      # Visual search pattern display
```

---

## As an AI Skill

See [SKILL.md](SKILL.md) for full skill documentation including:
- **Why this exists** (beyond rg/fd)
- **Typical workflow** for agents
- **All command examples** with JSON output
- **Flag reference** with applies-to column

---

## Related Projects

| Project | Description |
|---------|-------------|
| [pi-skills](https://github.com/junghan0611/pi-skills) | AI skill collection — denotecli distributed as a skill |
| [zotero-config](https://github.com/junghan0611/zotero-config) | Headless Zotero + **bibcli** (sister CLI for 8K+ bibliography) |

---

**Author**: [@junghanacs](https://github.com/junghan0611) · **License**: Apache 2.0
