#!/usr/bin/env python3
"""denote_parser.py

Parse Denote filenames and org-mode frontmatter.

Denote filename format:
    YYYYMMDDTHHMMSS--title-words__tag1_tag2.org

Org frontmatter:
    #+title:      Title
    #+date:       [2025-10-21 Tue 10:53]
    #+filetags:   :tag1:tag2:
    #+identifier: 20251021T105353

Usage:
    # Parse filename only
    python denote_parser.py --filename "20251021T105353--제목__tag.org"

    # Parse file (filename + frontmatter)
    python denote_parser.py path/to/file.org

    # Output as JSON
    python denote_parser.py --json path/to/file.org
"""

import re
import sys
import json
from pathlib import Path
from typing import Optional


# Denote filename pattern
# YYYYMMDDTHHMMSS--title-words__tag1_tag2.ext
DENOTE_FILENAME_RE = re.compile(
    r"^(?P<timestamp>\d{8}T\d{6})"  # YYYYMMDDTHHMMSS
    r"--"
    r"(?P<title>[^_]+?)"           # title (before __)
    r"(?:__(?P<tags>[^.]+))?"      # optional __tags
    r"\.(?P<ext>\w+)$"             # .extension
)

# Frontmatter patterns
TITLE_RE = re.compile(r"^#\+title:\s*(.+)$", re.IGNORECASE)
DATE_RE = re.compile(r"^#\+date:\s*(.+)$", re.IGNORECASE)
FILETAGS_RE = re.compile(r"^#\+filetags:\s*:?([^:]+(?::[^:]+)*):?$", re.IGNORECASE)
IDENTIFIER_RE = re.compile(r"^#\+identifier:\s*(\d{8}T\d{6})$", re.IGNORECASE)


def parse_denote_filename(filename: str) -> Optional[dict]:
    """Parse a Denote filename into components.

    Args:
        filename: Just the filename (not full path)

    Returns:
        dict with timestamp, title, tags, ext or None if not a Denote file
    """
    # Get just the filename if path given
    name = Path(filename).name

    match = DENOTE_FILENAME_RE.match(name)
    if not match:
        return None

    title_raw = match.group("title")
    # Convert hyphens to spaces, handle Korean
    title = title_raw.replace("-", " ").strip()

    tags_raw = match.group("tags")
    tags = tags_raw.split("_") if tags_raw else []

    return {
        "timestamp": match.group("timestamp"),
        "title": title,
        "tags": tags,
        "ext": match.group("ext"),
    }


def parse_frontmatter(content: str) -> dict:
    """Parse org-mode frontmatter from file content.

    Args:
        content: Full file content or first ~20 lines

    Returns:
        dict with title, date, filetags, identifier (all optional)
    """
    result = {}

    # Only check first 20 lines for frontmatter
    lines = content.split("\n")[:20]

    for line in lines:
        line = line.strip()

        if m := TITLE_RE.match(line):
            result["title"] = m.group(1).strip()
        elif m := DATE_RE.match(line):
            result["date"] = m.group(1).strip()
        elif m := FILETAGS_RE.match(line):
            tags_str = m.group(1)
            result["filetags"] = [t for t in tags_str.split(":") if t]
        elif m := IDENTIFIER_RE.match(line):
            result["identifier"] = m.group(1)

    return result


def parse_denote_file(filepath: str) -> dict:
    """Parse both filename and frontmatter of a Denote file.

    Args:
        filepath: Path to the org file

    Returns:
        dict with filename_meta and frontmatter
    """
    path = Path(filepath).expanduser()

    result = {
        "path": str(path),
        "filename": path.name,
        "filename_meta": parse_denote_filename(path.name),
        "frontmatter": {},
    }

    if path.exists():
        try:
            content = path.read_text(encoding="utf-8")
            result["frontmatter"] = parse_frontmatter(content)
        except Exception as e:
            result["error"] = str(e)

    return result


def main(argv=None):
    if argv is None:
        argv = sys.argv[1:]

    if not argv:
        print("Usage: denote_parser.py [--json] [--filename] <file_or_name>",
              file=sys.stderr)
        return 1

    output_json = "--json" in argv
    filename_only = "--filename" in argv
    argv = [a for a in argv if not a.startswith("--")]

    if not argv:
        print("ERROR: No file or filename provided", file=sys.stderr)
        return 1

    target = argv[0]

    if filename_only:
        result = parse_denote_filename(target)
    else:
        result = parse_denote_file(target)

    if output_json:
        print(json.dumps(result, ensure_ascii=False, indent=2))
    else:
        if result is None:
            print("Not a valid Denote filename")
            return 1

        if filename_only:
            print(f"Timestamp: {result['timestamp']}")
            print(f"Title: {result['title']}")
            print(f"Tags: {', '.join(result['tags']) if result['tags'] else '(none)'}")
            print(f"Extension: {result['ext']}")
        else:
            print(f"Path: {result['path']}")
            if result['filename_meta']:
                meta = result['filename_meta']
                print(f"ID: {meta['timestamp']}")
                print(f"Title (filename): {meta['title']}")
                print(f"Tags (filename): {', '.join(meta['tags']) if meta['tags'] else '(none)'}")
            if result['frontmatter']:
                fm = result['frontmatter']
                if 'title' in fm:
                    print(f"Title (frontmatter): {fm['title']}")
                if 'identifier' in fm:
                    print(f"Identifier: {fm['identifier']}")
                if 'filetags' in fm:
                    print(f"Filetags: {', '.join(fm['filetags'])}")
                if 'date' in fm:
                    print(f"Date: {fm['date']}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
