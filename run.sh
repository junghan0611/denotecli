#!/usr/bin/env bash
# run.sh — denotecli 프로젝트 진입점
#
# Usage:
#   ./run.sh build          — Build + install + copy skill
#   ./run.sh test           — Run all tests
#   ./run.sh showcase       — Run showcase tests (visual pattern check)
#   ./run.sh cover          — Run tests with coverage report
#   ./run.sh bench          — Run benchmarks
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_DIR="$SCRIPT_DIR/denotecli"

case "${1:-}" in
    build)
        echo "Building denotecli..."
        INSTALL_DIR="${2:-$HOME/.local/bin}"
        mkdir -p "$INSTALL_DIR"
        (cd "$GO_DIR" && go build -o "$INSTALL_DIR/denotecli" .)
        (cd "$GO_DIR" && GOOS=linux GOARCH=arm64 go build -o "$INSTALL_DIR/denotecli-linux-arm64" .)
        echo "Installed: $INSTALL_DIR/denotecli ($(uname -m))"
        echo "Installed: $INSTALL_DIR/denotecli-linux-arm64 (aarch64)"
        # Install skill to pi-skills (SKILL.md + binary)
        SKILL_DIR="$HOME/repos/gh/pi-skills/denotecli"
        if [[ -d "$HOME/repos/gh/pi-skills" ]]; then
            mkdir -p "$SKILL_DIR"
            cp "$SCRIPT_DIR/SKILL.md" "$SKILL_DIR/SKILL.md"
            cp "$INSTALL_DIR/denotecli" "$SKILL_DIR/denotecli"
            cp "$INSTALL_DIR/denotecli-linux-arm64" "$SKILL_DIR/denotecli-linux-arm64"
            echo "Skill installed: $SKILL_DIR/ (SKILL.md + x86_64 + arm64)"
        fi
        ;;
    test)
        echo "Running all tests..."
        (cd "$GO_DIR" && go test -v -count=1 ./...)
        ;;
    showcase)
        echo "Running showcase (visual pattern check)..."
        (cd "$GO_DIR" && go test -v -run "TestShowcase" -count=1 ./...)
        ;;
    cover)
        echo "Running coverage..."
        (cd "$GO_DIR" && go test -cover -coverprofile=cover.out -count=1 ./...)
        echo ""
        echo "=== Coverage by function ==="
        (cd "$GO_DIR" && go tool cover -func=cover.out | grep -v "0.0%" | tail -20)
        echo ""
        (cd "$GO_DIR" && go tool cover -func=cover.out | tail -1)
        rm -f "$GO_DIR/cover.out"
        ;;
    bench)
        echo "Running benchmarks..."
        (cd "$GO_DIR" && go test -bench=. -benchtime=3s -v)
        ;;
    -h|--help|help|"")
        cat <<'EOF'
denotecli — Denote knowledge base CLI for AI agents

Usage:
  ./run.sh build [DIR]   Build + install to DIR (default: ~/.local/bin)
  ./run.sh test          Run all 82+ tests
  ./run.sh showcase      Run showcase tests (visual search pattern check)
  ./run.sh cover         Run tests with coverage report
  ./run.sh bench         Run benchmarks
EOF
        ;;
    *)
        echo "Unknown command: $1" >&2
        "$0" --help
        exit 1
        ;;
esac
