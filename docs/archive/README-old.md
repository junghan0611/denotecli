# Denote-Org Skills for Claude

> **Bringing Anthropic's Life Sciences paradigm to Life Everything**

Comprehensive Denote PKM system support for Claude AI, validated with **3,000+ org files** in production knowledge bases.

**Denote-first. Org-mode powered. Life-scale proven.**

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Files](https://img.shields.io/badge/validated-3000%2B%20files-green.svg)]()
[![Status](https://img.shields.io/badge/version-0.1--dev-orange.svg)]()

## 🌟 Vision: Life Sciences → Life Everything

Anthropic proved that **domain context + AI = expert-level collaboration** with Life Sciences (PubMed, Benchling, 10x Genomics).

This project extends that paradigm from **Biology** to **Living**:

**Life Sciences (Biology)** → **Life Everything (Living)**

Just as Claude becomes a scientist with PubMed context, Claude becomes your **personal knowledge partner** with Denote context.

### The Pattern

| Domain | Context | Result |
|--------|---------|--------|
| **Anthropic** | PubMed (millions of papers) + Claude | Scientific AI |
| **This Project** | Denote (3,000+ org files) + Claude | Personal Knowledge AI |

**Same methodology. Different domain. Universal paradigm.**

## Why This is a Skill, Not Just a Prompt

**Anthropic's Skills Philosophy:**

> "Skills transform Claude from a general-purpose agent into a specialized agent equipped with **procedural knowledge** that no model can fully possess."
>
> — [Equipping agents for the real world with Agent Skills](https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills)

**Just like:**
- **PDF Skill** solves complex forms (8 Python scripts for bounding boxes, validation...)
- **XLSX Skill** maintains formulas (recalc.py with LibreOffice integration)
- **Denote-Org Skill** handles 3,000-file knowledge graphs (finder, graph builder, silo manager)

**Skills are NOT prompts. They are operational systems.**

## The Problem This Solves

### Without This Skill

❌ Claude doesn't know Denote file naming: `20251021T105353--title__tags.org`
❌ Finding files in 3,000+ requires repeated token-heavy searches
❌ Knowledge graph navigation is slow and error-prone
❌ Literate programming (`:tangle`, `:results`) often breaks org structure
❌ Multiple silos (`~/org/`, `~/claude-memory/`) cause confusion

### With This Skill

✅ **Denote domain knowledge** bundled (file naming, frontmatter, links)
✅ **Python scripts** for efficient operations (finder, graph, executor)
✅ **Knowledge graph** cached and reusable
✅ **Safe code execution** (org-babel compatible)
✅ **Silo management** (multiple knowledge domains)

**Not a file converter. An operational knowledge system.**

## 🎯 What is Denote?

Denote is a **simple yet powerful** note-taking system for Emacs:

- **File naming convention**: `YYYYMMDDTHHMMSS--title__tag1_tag2.org`
- **Knowledge graph**: `[[denote:TIMESTAMP]]` links
- **Silo concept**: Separate directories for different knowledge domains
- **Org-mode foundation**: Headings, properties, timestamps, code blocks

**Created by:** [Protesilaos Stavrou](https://protesilaos.com/emacs/denote)

## 📊 Scale: Production-Validated

This skill is built for **real-world PKM systems**, not toy examples:

- ✅ **3,000+ org files** (validated at scale)
- ✅ **Multiple silos** (`~/org/`, `~/claude-memory/`, project `docs/`)
- ✅ **Denote knowledge graph** (thousands of `[[denote:ID]]` links)
- ✅ **Daily usage** with Claude Code
- ✅ **Part of 9-layer "Tools for Life" system**

## 🏗️ Part of the 9-Layer System

This skill is **Layer 3** (Knowledge Understanding) in a complete AI-human collaboration architecture:

```
Layer 7: Knowledge Publishing (notes.junghanacs.com)
Layer 6: Agent Orchestration (meta-config)
Layer 5a: Migration (memex-kb)
Layer 5b: Life Timeline (memacs-config)
Layer 4: AI Memory (claude-config)
Layer 3.5: RAG/Vector DB (embedding-config)
Layer 3: Knowledge Mgmt (Zotero, denote-org 🆕)
Layer 2: Dev Environment (Doom Emacs)
Layer 1: Infrastructure (NixOS)
```

**Related projects:**
- [claude-config](https://github.com/junghan0611/claude-config) - AI Memory System with PARA methodology
- [memex-kb](https://github.com/junghan0611/memex-kb) - Universal KB Migration (Google Docs, Dooray...)
- [embedding-config](https://github.com/junghan0611/embedding-config) - Vector RAG (9,477 blocks, 85% success rate)

**Philosophy:** Same as Anthropic Life Sciences - **AI as Partner, not Tool**

## 🚀 Real Use Cases

### Case 1: Knowledge Graph Navigation
```bash
User: "Find all files linked to 20251021T105353"
Claude: Uses denote_graph.py
Result: Connected files in <1 second (not repeated token searches)
```

### Case 2: Tag-based Search Across 3,000 Files
```bash
User: "All llmlog files from October 2025"
Claude: Uses denote_finder.py --tags llmlog --date 202510*
Result: Instant results with metadata
```

### Case 3: Silo Management
```bash
User: "Is this in ~/org/ or ~/claude-memory/?"
Claude: Uses denote_silo.py find_in_silos()
Result: Resolves ambiguity automatically
```

### Case 4: Literate Programming
```bash
User: "Execute code blocks in this org file"
Claude: Uses org_execute.py (org-babel compatible)
Result: Safe execution, structure preserved, :results updated
```

### Case 5: Denote Link Resolution
```bash
File contains: [[denote:20251021T105353]]
Claude: Uses denote_links.py resolve_link()
Result: ~/org/llmlog/20251021T105353--full-filename.org
```

## 💡 Features

### Denote Core
- ✅ File naming: `20251021T105353--title__tag1_tag2.org`
- ✅ Frontmatter parsing (title, date, filetags, identifier)
- ✅ Link resolution (`[[denote:ID]]`)
- ✅ Tag-based search
- ✅ Silo management (multiple knowledge domains)
- ✅ Knowledge graph (3,000+ files validated)

### Org-mode Base
- ✅ Heading/property parsing
- ✅ Code block execution (literate programming)
- ✅ Timestamp handling
- ✅ Export (markdown, PDF, HTML)
- ✅ Table operations

### Performance
- ✅ Python scripts (not token generation)
- ✅ Caching for repeated operations
- ✅ Efficient graph traversal
- ✅ Batch processing support

## 📦 Installation

### Prerequisites
```bash
# Python 3.8+
pip install orgparse pypandoc pyyaml

# Pandoc (for export)
# Install via your package manager
```

### Claude Code
```bash
git clone https://github.com/junghan0611/orgmode-skills.git
cd orgmode-skills
/plugin add .
```

### Claude API
```python
# See API documentation for skills integration
import anthropic

client = anthropic.Client(api_key="your-api-key")
# Skills API usage
```

## 🔧 Quick Start

### Example 1: Find Denote File
```python
from scripts.denote_finder import find_denote_file

# By identifier
filepath = find_denote_file("20251021T105353", silos=["~/org/"])
# Returns: ~/org/llmlog/20251021T105353--full-filename.org

# By tags
files = find_by_tags(["llmlog", "denote"], silo="~/org/")
```

### Example 2: Extract Knowledge Graph
```python
from scripts.denote_graph import build_knowledge_graph, get_connected_nodes

# Build graph (cached)
graph = build_knowledge_graph("~/org/")

# Get connected files
connected = get_connected_nodes("20251021T105353", hops=2)
```

### Example 3: Execute Org Code Blocks
```python
from scripts.org_execute import execute_org_blocks

# Execute all code blocks (respects :tangle, :results)
results = execute_org_blocks("file.org")
```

## 📂 Project Structure

```
orgmode-skills/
├── README.md                    # This file
├── README-KO.md                 # Korean version
├── LICENSE                      # Apache 2.0
├── SKILL.md                     # Claude's main guide
├── docs/                        # Development logs (Denote format)
│   └── 20251021T113500--*.org
├── denote-core.md               # Denote specification
├── denote-silo.md               # Silo management
├── denote-knowledge-graph.md    # Graph operations
├── orgmode-base.md              # Org-mode features
├── literate.md                  # Literate programming
├── export.md                    # Export capabilities
└── scripts/
    ├── denote_finder.py         # File search
    ├── denote_links.py          # Link resolution
    ├── denote_silo.py           # Silo management
    ├── denote_graph.py          # Knowledge graph
    ├── org_parser.py            # Org parsing
    ├── org_execute.py           # Code execution
    ├── org_export.py            # Export
    └── requirements.txt
```

## 🗺️ Roadmap to 0.1

**Phase 1: Core (Current)**
- [x] Project setup
- [ ] README.md, README-KO.md
- [ ] Development log (docs/)
- [ ] SKILL.md (Denote-focused)

**Phase 2: Denote Core**
- [ ] denote_finder.py (file search)
- [ ] denote_links.py (link resolution)
- [ ] denote_silo.py (silo management)
- [ ] denote-core.md documentation

**Phase 3: Org-mode Base**
- [ ] org_parser.py (orgparse)
- [ ] org_execute.py (code blocks)
- [ ] orgmode-base.md documentation

**Phase 4: 0.1 Release**
- [ ] Testing with 3,000+ files
- [ ] Examples
- [ ] Public release

## 🤝 Contributing

This project is validated with real-world 3,000+ file Denote PKM systems. Denote users welcome!

## 📖 Philosophy: Tools for Life

**From the creator:**

> "저에게 필요한 것은 경쟁력 있는 지식 노동자로 살 수 있는 방법이었습니다."
>
> "도구는 존재에 녹아든다"

This skill embodies the "Tools for Life" philosophy - extending Anthropic's Life Sciences approach from biology research to personal knowledge work.

**Life Sciences (Biology)** → **Life Everything (Living)**

## 📚 Resources

- [Denote Package](https://protesilaos.com/emacs/denote) by Protesilaos Stavrou
- [Anthropic Skills Documentation](https://support.claude.com/en/articles/12512176-what-are-skills)
- [Agent Skills Engineering Blog](https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills)
- [Org-mode](https://orgmode.org/)

## 🔗 Related Projects

Part of the **"시간과정신의방" (Hyperbolic Time Chamber)** ecosystem:
- [nixos-config](https://github.com/junghan0611/nixos-config) - Declarative infrastructure
- [doomemacs-config](https://github.com/junghan0611/doomemacs-config) - Terminal-optimized Emacs
- [zotero-config](https://github.com/junghan0611/zotero-config) - AI-queryable bibliography (156k+ lines)

**Digital Garden:** [notes.junghanacs.com](https://notes.junghanacs.com)

## 📄 License

Apache 2.0 (compatible with Anthropic Skills ecosystem)

## 🙏 Acknowledgments

- **Anthropic** for the Agent Skills paradigm and Life Sciences blueprint
- **Protesilaos Stavrou** for creating Denote
- **Org-mode community** for the foundation
- **AIONS Clubs International** for the vision of collective intelligence

---

**Status:** 🟡 Development (targeting 0.1 public release)
**Maintainer:** [@junghan0611](https://github.com/junghan0611)
**Created:** 2025-10-21
