# orgmode-skills

Denote knowledge base CLI (`denotecli`) for AI agents. Searches, reads, and analyzes 3,000+ org-mode files with Korean support.

**Go stdlib only. Single binary. JSON output.**

## Install

```bash
git clone https://github.com/junghan0611/denotecli.git
cd denotecli

# Build + install to ~/.local/bin + copy skill to pi-skills
./run.sh build
```

Requires Go 1.21+. No external dependencies.

## Quick Start

```bash
denotecli search "에릭 호퍼" --dirs ~/org --max 5
denotecli read 20250314T152111 --dirs ~/org
denotecli tags --dirs ~/org --top 10
```

## As a Skill

### Claude Code / pi-coding-agent

```bash
ln -s /path/to/orgmode-skills ~/.claude/skills/denote-org
```

### OpenClaw / Container

```bash
# Use --dirs /data/org for container paths
denotecli search "query" --dirs /data/org
```

## Commands

| Command | Description |
|---------|-------------|
| `search <query>` | Search notes by title, tags, ID |
| `read <id>` | Read note content + metadata + links |
| `tags` | Tag statistics across all files |

See [SKILL.md](SKILL.md) for full documentation.

## License

Apache 2.0
