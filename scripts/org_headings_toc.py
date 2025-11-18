#!/usr/bin/env python3
"""org_headings_toc.py

Extract a lightweight table of contents (TOC) from an org-mode file.

This script is intentionally simple and robust: it does not try to fully
parse org-mode semantics, only heading levels and titles. It is designed
for agents to call when they need a quick structural overview of a large
org file before deciding which sections to read in detail.

Behavior:
- Reads the given org file line by line.
- Detects headings by leading '*' characters (Org-style headings).
- Determines heading level by the number of consecutive '*' at the
  beginning of the line.
- Strips TODO keywords and leading tags like "* TODO", "* DONE" when
  extracting the title.
- Prints a simple TOC to stdout in the form:
    LEVEL<TAB>TITLE
  where LEVEL is an integer (1 for top-level headings, 2 for subheadings, ...).

Usage:
    python org_headings_toc.py path/to/file.org

This script is meant to be called by Claude/agents via a shell command
when using the denote-org skill.
"""

import sys
from pathlib import Path


TODO_KEYWORDS = {
    "TODO",
    "DONE",
    "WAITING",
    "HOLD",
    "CANCELLED",
    "NEXT",
}


def extract_headings(path: Path):
    """Yield (level, title) tuples for each heading in the org file."""
    try:
        with path.open("r", encoding="utf-8") as f:
            for line in f:
                if not line.lstrip().startswith("*"):
                    continue

                stripped = line.rstrip("\n")
                # Count leading '*' characters
                i = 0
                while i < len(stripped) and stripped[i] == "*":
                    i += 1
                level = i

                # Require at least one space after the stars
                if level == 0 or level >= len(stripped) or stripped[level] != " ":
                    continue

                # Extract raw title part
                raw_title = stripped[level + 1 :].strip()

                # Remove TODO keywords at the beginning of the title
                parts = raw_title.split()
                if parts and parts[0] in TODO_KEYWORDS:
                    parts = parts[1:]
                title = " ".join(parts).strip()

                if title:
                    yield level, title
    except FileNotFoundError:
        print(f"ERROR: File not found: {path}", file=sys.stderr)
    except UnicodeDecodeError:
        print(f"ERROR: Cannot decode file as UTF-8: {path}", file=sys.stderr)


def main(argv=None):
    if argv is None:
        argv = sys.argv[1:]

    if not argv:
        print("Usage: org_headings_toc.py path/to/file.org", file=sys.stderr)
        return 1

    path = Path(argv[0]).expanduser()

    for level, title in extract_headings(path):
        # LEVEL<TAB>TITLE format for easy parsing by the agent
        print(f"{level}\t{title}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
