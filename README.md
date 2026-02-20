# orgmode-skills

Denote knowledge base CLI (`denotecli`) for AI agents. Searches, reads, and analyzes 3,000+ org-mode files.

## Quick Start

```bash
cd denotecli && go build -o denotecli .
./denotecli search "에릭 호퍼" --dirs ~/org
./denotecli read 20250314T152111 --dirs ~/org
./denotecli tags --dirs ~/org --top 10
```

## As a Skill

### Claude Code / pi-coding-agent

```bash
# Symlink to skills directory
ln -s /path/to/orgmode-skills ~/.claude/skills/denote-org
```

### OpenClaw / Container

```bash
# Build and copy binary
cd denotecli && go build -o denotecli .
# Use --dirs /data/org for container paths
```

## License

Apache 2.0
