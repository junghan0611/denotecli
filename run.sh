#!/usr/bin/env bash
# run.sh — orgmode-skills 프로젝트 메인 진입점
#
# Usage:
#   ./run.sh build [INSTALL_DIR]   — Build denotecli + install + copy skill
#   ./run.sh test                  — Run all tests
#   ./run.sh bench                 — Run benchmarks
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

case "${1:-}" in
    build)
        echo "Building denotecli..."
        INSTALL_DIR="${2:-$HOME/.local/bin}"
        mkdir -p "$INSTALL_DIR"
        (cd "$SCRIPT_DIR/denotecli" && go build -o "$INSTALL_DIR/denotecli" .)
        echo "Installed: $INSTALL_DIR/denotecli"
        # Install skill to pi-skills
        SKILL_DIR="$HOME/repos/gh/pi-skills/denotecli"
        if [[ -d "$HOME/repos/gh/pi-skills" ]]; then
            mkdir -p "$SKILL_DIR"
            cp "$SCRIPT_DIR/SKILL.md" "$SKILL_DIR/SKILL.md"
            echo "Skill installed: $SKILL_DIR/SKILL.md"
        fi
        ;;
    test)
        echo "Running tests..."
        (cd "$SCRIPT_DIR/denotecli" && go test -v ./...)
        ;;
    bench)
        echo "Running benchmarks..."
        (cd "$SCRIPT_DIR/denotecli" && go test -bench=. -benchtime=3s -v)
        ;;
    -h|--help|help|"")
        cat <<'EOF'
orgmode-skills — Denote knowledge base CLI

Usage:
  ./run.sh build [DIR]   Build denotecli, install to DIR (default: ~/.local/bin)
  ./run.sh test          Run all tests
  ./run.sh bench         Run benchmarks
EOF
        ;;
    *)
        echo "Unknown command: $1" >&2
        "$0" --help
        exit 1
        ;;
esac
