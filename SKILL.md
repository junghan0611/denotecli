---
name: denote-org
description: "Comprehensive support for Denote-based PKM systems with org-mode foundation. Handles 3,000+ file knowledge bases with Denote file naming (timestamp--title__tags.org), knowledge graph navigation via [[denote:ID]] links, multi-silo management, and literate programming. When Claude needs to: (1) Find Denote files by ID/tags in large knowledge bases, (2) Navigate knowledge graphs with thousands of [[denote:]] links, (3) Execute org-mode code blocks with :tangle/:results options, (4) Manage multiple Denote silos (~/org/, ~/claude-memory/...), or (5) Work with Denote-structured PKM systems at scale."
---

# Denote-Org Skills

## Overview

This skill enables Claude to work with Denote-based Personal Knowledge Management (PKM) systems at scale. Denote is a simple yet powerful note-taking system for Emacs that combines file naming conventions, org-mode structure, and knowledge graph links. This skill provides procedural knowledge and executable scripts for handling 3,000+ file knowledge bases efficiently.

## When to Use This Skill

Claude should use this skill when:

1. **Working with Denote files** - Files named `YYYYMMDDTHHMMSS--title__tag1_tag2.org`
2. **Large-scale PKM systems** - Knowledge bases with hundreds or thousands of org files
3. **Knowledge graph navigation** - Following `[[denote:TIMESTAMP]]` links across files
4. **Multi-silo management** - Searching across multiple Denote directories (~/org/, ~/claude-memory/, project docs/)
5. **Literate programming** - Executing org-mode code blocks with `:tangle` and `:results` options
6. **Denote metadata extraction** - Parsing file names and org-mode frontmatter

## What is Denote?

Denote is a note-taking system created by Protesilaos Stavrou that uses:

- **File naming**: `YYYYMMDDTHHMMSS--title__tag1_tag2.org`
  - Timestamp = Unique identifier (Denote ID)
  - Title = Hyphen-separated (supports 한글/Korean)
  - Tags = Underscore-separated (English)

- **Org-mode frontmatter**:
  ```org
  #+title:      Title
  #+date:       [2025-10-21 Tue 10:53]
  #+filetags:   :tag1:tag2:tag3:
  #+identifier: 20251021T105353
  ```

- **Knowledge graph**: `[[denote:20251021T105353]]` links between files

- **Silo concept**: Separate directories for different knowledge domains

**Validated with 3,000+ org files in production PKM systems.**

## Core Capabilities

### 1. Denote File Name Parsing

Parse Denote file names to extract metadata:

```python
# Use scripts/denote_parser.py
from scripts.denote_parser import parse_denote_filename

metadata = parse_denote_filename("20251021T105353--project-title__llmlog_denote.org")
# Returns:
# {
#     'timestamp': '20251021T105353',
#     'title': 'project title',
#     'tags': ['llmlog', 'denote'],
#     'extension': 'org'
# }
```

**When to use:**
- Extracting Denote ID from filename
- Getting title and tags without opening file
- Processing multiple files in batch

**Edge cases handled:**
- Korean (한글) titles
- Filenames without tags
- Multiple hyphens in title

### 2. Finding Denote Files

Search for Denote files by ID, tags, or date patterns:

```python
# Use scripts/denote_finder.py
from scripts.denote_finder import find_denote_file, find_by_tags

# Find by ID
filepath = find_denote_file("20251021T105353", silos=["~/org/"])

# Find by tags
files = find_by_tags(["llmlog", "denote"], silo="~/org/")

# Find by date pattern
files = find_by_date_range("202510*", silo="~/org/")
```

**When to use:**
- User asks "Find file 20251021T105353"
- Search by tags: "Show me all llmlog files"
- Date range queries: "Files from October 2025"

**Performance:**
- Uses caching for repeated searches
- Generator pattern for large result sets
- Typical search: <100ms for 3,000+ files

### 3. Denote Link Resolution

Resolve `[[denote:ID]]` links to actual file paths:

```python
# Use scripts/denote_links.py
from scripts.denote_links import resolve_denote_link

filepath = resolve_denote_link("denote:20251021T105353", silos=["~/org/", "~/claude-memory/"])
# Returns: ~/org/llmlog/20251021T105353--full-filename.org
```

**When to use:**
- File contains `[[denote:...]]` links
- Need to open linked files
- Validate link integrity

**Error recovery:**
- Exact match (full timestamp)
- Fallback: Partial match (first 8 chars)
- Returns suggestions if not found

### 4. Knowledge Graph Operations

Build and navigate knowledge graphs from Denote links:

```python
# Use scripts/denote_graph.py
from scripts.denote_graph import build_knowledge_graph, get_connected_nodes

# Build graph (cached)
graph = build_knowledge_graph("~/org/")

# Get connected files
connected = get_connected_nodes("20251021T105353", hops=2)

# Detect circular references
cycles = detect_circular_refs(graph)
```

**When to use:**
- "Show files linked to this"
- "Find related notes"
- Discover knowledge clusters
- Detect broken or circular links

**Performance:**
- First build: ~10 seconds for 3,000 files
- Cached: <1 second
- Incremental updates supported

### 5. Multi-Silo Management

Search across multiple Denote directories:

```python
# Use scripts/denote_silo.py
from scripts.denote_silo import find_in_silos

# Search in multiple locations
filepath = find_in_silos(
    "20251021T105353",
    silos=["~/org/", "~/claude-memory/", "~/repos/*/docs/"],
    priority=True  # Return first match
)
```

**When to use:**
- Files could be in multiple locations
- Cross-silo link resolution
- Silo statistics and analysis

**Common silos:**
- `~/org/` - Main knowledge base (3,000+ files)
- `~/claude-memory/` - AI memory system
- `~/repos/*/docs/` - Project-specific docs

### 6. Org-mode Code Block Execution

Execute org-mode code blocks (literate programming):

```python
# Use scripts/org_execute.py
from scripts.org_execute import execute_org_blocks

# Execute all code blocks
results = execute_org_blocks("file.org")

# Execute specific named block
result = execute_named_block("file.org", "example-block")

# Tangle code to files
tangle_org_file("file.org", output_dir="/tmp/")
```

**When to use:**
- User asks to "execute code blocks"
- Literate programming workflows
- Generate files from `:tangle` options

**Org-babel compatibility:**
- Supports `:tangle`, `:results`, `:name` options
- Languages: bash, python, etc.
- Preserves org structure
- Updates `#+RESULTS` blocks

## Quick Start Examples

### Example 1: Find and Open Denote File

```
User: "Open file 20251021T105353"

Claude: Uses denote_finder.py
Result: ~/org/llmlog/20251021T105353--anthropic-skill-creator-분석__llmlog_skills.org
```

### Example 2: Knowledge Graph Exploration

```
User: "Show me all files connected to this note"

Claude: Uses denote_graph.py
Result:
- Direct links: 5 files
- 2-hop connections: 23 files
- Clusters: AI/Skills, PKM, Tools for Life
```

### Example 3: Tag-based Search

```
User: "Find all llmlog files from October 2025"

Claude: Uses denote_finder.py
Result: 47 files matching criteria
```

### Example 4: Execute Literate Programming

```
User: "Execute the bash blocks in this org file"

Claude: Uses org_execute.py
Result:
- Executed 3 code blocks
- Updated #+RESULTS
- Tangled to /tmp/script.sh
```

## Common Pitfalls

### Pitfall 1: Denote ID Collision
**Problem:** Files created in same second have identical timestamps
**Solution:** Scripts use full filename as tiebreaker

### Pitfall 2: Cross-Silo Links
**Problem:** Link from ~/org/ to ~/claude-memory/ file
**Solution:** Always search multiple silos, use `find_in_silos()`

### Pitfall 3: Circular References
**Problem:** A → B → C → A creates cycle in knowledge graph
**Solution:** `detect_circular_refs()` identifies cycles

### Pitfall 4: Korean Encoding
**Problem:** 한글 filenames may cause encoding issues
**Solution:** All scripts use UTF-8, `os.path.expanduser()`

### Pitfall 5: Large Result Sets
**Problem:** Query returns 1,000+ files, exceeds token limit
**Solution:** Scripts use generators, return summary + pagination

## File Structure

This skill uses a **Capabilities-Based** structure with executable scripts and reference documentation:

```
denote-org/
├── SKILL.md (this file)
├── scripts/
│   ├── denote_parser.py      - Parse filenames and frontmatter
│   ├── denote_finder.py      - Search files (cached, generator)
│   ├── denote_links.py       - Resolve [[denote:ID]] links
│   ├── denote_silo.py        - Multi-silo management
│   ├── denote_graph.py       - Knowledge graph operations
│   ├── org_parser.py         - Parse org structure
│   ├── org_execute.py        - Execute code blocks
│   └── org_export.py         - Export to markdown/PDF/HTML
├── references/
│   ├── denote-spec.md        - Denote specification
│   ├── orgmode-syntax.md     - Org-mode syntax reference
│   └── literate-programming.md - Literate programming guide
└── assets/
    └── (none - this skill doesn't use output templates)
```

## Resources

### scripts/

Executable Python scripts for efficient Denote and org-mode operations. These scripts handle:

- **File operations**: Finding, parsing, and resolving links across 3,000+ files
- **Graph operations**: Building and navigating knowledge graphs
- **Code execution**: Running org-babel blocks safely
- **Export**: Converting org to other formats

**Performance optimized:**
- Caching with `@lru_cache` for expensive operations
- Generator pattern for memory efficiency
- Typical operations: <100ms (finder), <1s (graph, cached)

### references/

Detailed documentation for Denote and org-mode:

- **denote-spec.md**: Complete Denote specification (file naming, frontmatter, links, silos)
- **orgmode-syntax.md**: Org-mode syntax reference (headings, properties, timestamps, code blocks)
- **literate-programming.md**: Literate programming with org-babel (`:tangle`, `:results`, `:name`)

**Load these when:**
- Working with complex org structures
- Need detailed Denote specification
- Implementing literate programming workflows

### assets/

Not used in this skill. Denote-org processes existing files rather than creating from templates.

**Note:** The `assets/` directory can be deleted as it's not needed for this skill.

## Integration with Other Systems

This skill is part of the **"Tools for Life" 9-layer architecture**, specifically Layer 3 (Knowledge Management):

- **Layer 4 (claude-config)**: Uses Denote format for AI memory (~/claude-memory/)
- **Layer 3.5 (embedding-config)**: Vectorizes org files for semantic search
- **Layer 5a (memex-kb)**: Migrates to Denote format
- **Layer 7 (notes)**: Exports Denote org → markdown for publishing

**Philosophy:** Extends Anthropic's Life Sciences paradigm (domain context + AI = expert collaboration) from Biology research to personal knowledge work.

## Performance Considerations

When working with 3,000+ files:

- **Caching**: Graph building is cached (10s first time, <1s after)
- **Generators**: Large searches use streaming (memory efficient)
- **Timeouts**: Long operations have configurable limits
- **Pagination**: Large result sets are paginated automatically

## Troubleshooting

### Links Not Resolving
- Check if file exists in specified silos
- Try partial match (first 8 chars of timestamp)
- Verify silo paths use `os.path.expanduser()`

### Graph Building Slow
- First build takes ~10 seconds for 3,000 files (normal)
- Subsequent builds use cache (<1 second)
- Use incremental update for changed files only

### Code Block Execution Fails
- Check `:tangle` path exists
- Verify code block language support
- Review `:results` option syntax

### Encoding Errors
- All scripts use UTF-8
- Korean (한글) filenames fully supported
- File paths expanded with `os.path.expanduser()`

## See Also

- [Denote Package](https://protesilaos.com/emacs/denote) by Protesilaos Stavrou
- [Org-mode](https://orgmode.org/) documentation
- README.md in this repository for project context and Life Sciences paradigm connection
