// parser.go
package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// DenoteFile represents parsed Denote file metadata.
type DenoteFile struct {
	ID          string   `json:"id"`
	Signature   string   `json:"signature,omitempty"`
	Title       string   `json:"title"`                  // filename slug (backward-compat)
	HeaderTitle string   `json:"header_title,omitempty"` // #+title: header value (Bug 4 fix)
	Tags        []string `json:"tags"`                   // union: filename slots ∪ #+filetags: (Bug 1 fix)
	Date        string   `json:"date"`
	Lastmod     string   `json:"lastmod,omitempty"`      // #+hugo_lastmod: normalized to YYYY-MM-DD
	HugoLastmod string   `json:"hugo_lastmod,omitempty"` // #+hugo_lastmod: raw org timestamp (keeps HH:MM)
	Description string   `json:"description,omitempty"`  // #+description:
	Path        string   `json:"path"`
}

// Abstract is a leading org callout block — the "[!abstract] 이 노트에 대하여"
// convention used across the garden. Kind is the callout word ("abstract",
// "note", ...), Title the text on the callout line, Body the paragraph(s) below.
type Abstract struct {
	Kind  string `json:"kind"`
	Title string `json:"title,omitempty"`
	Body  string `json:"body"`
}

// DenoteContent extends DenoteFile with file content.
type DenoteContent struct {
	DenoteFile
	Abstract *Abstract `json:"abstract,omitempty"`
	Content  string    `json:"content"`
	Links    []string  `json:"links"`
}

// OutlineEntry represents a single org heading.
type OutlineEntry struct {
	Level int    `json:"level"`
	Title string `json:"title"`
	Line  int    `json:"line"`
	Tags  string `json:"tags,omitempty"`
}

// DenoteOutline extends DenoteFile with outline (headings only).
type DenoteOutline struct {
	DenoteFile
	Abstract *Abstract      `json:"abstract,omitempty"`
	Outline  []OutlineEntry `json:"outline"`
	Links    []string       `json:"links"`
}

// Frontmatter holds parsed org-mode frontmatter fields.
type Frontmatter struct {
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Filetags    []string `json:"filetags"`
	Identifier  string   `json:"identifier"`
	Description string   `json:"description"`
	HugoLastmod string   `json:"hugo_lastmod,omitempty"`
}

var denoteRe = regexp.MustCompile(`^(\d{8}T\d{6})(?:==([^-]+))?--(.+?)(?:__(.+))?\.org$`)
var linkRe = regexp.MustCompile(`\[\[denote:(\d{8}T\d{6})\]`)
var dateRe = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

// ExtractDate finds the first YYYY-MM-DD substring.
// Used to normalize org timestamp variants — [2025-06-10], [2025-03-29 Sat 02:06],
// <2024-01-03 Wed 16:54>, 2023-06-19, "Time-stamp: <2025-01-31 12:54:06 junghan>".
// Returns empty string if no date pattern is found.
func ExtractDate(s string) string {
	return dateRe.FindString(s)
}

// ParseFilename extracts Denote metadata from a filename.
func ParseFilename(filename string) (DenoteFile, bool) {
	m := denoteRe.FindStringSubmatch(filename)
	if m == nil {
		return DenoteFile{}, false
	}

	df := DenoteFile{
		ID:        m[1],
		Signature: m[2],
		Title:     m[3],
	}

	if m[4] != "" {
		df.Tags = strings.Split(m[4], "_")
	}

	// Derive date from ID: 20251107T082610 -> 2025-11-07
	if len(df.ID) >= 8 {
		df.Date = df.ID[0:4] + "-" + df.ID[4:6] + "-" + df.ID[6:8]
	}

	return df, true
}

// ParseFrontmatter extracts org-mode frontmatter from file content.
//
// The frontmatter region ends at the first heading, the first block opener
// (`#+begin_...`), or the first line of real body text — NOT at the first blank
// line. Blank lines are routinely used to group keyword lines, and treating one
// as the terminator silently dropped `#+description:` on notes that place it
// below such a gap (measured on the corpus 2026-09-04: 20230904T144600,
// 20241216T152928 among others). The body-leak guard the blank-line rule was
// there for is kept by isFrontmatterLine: a `#+key:` inside a src block is
// unreachable because `#+begin_src` stops the scan first.
func ParseFrontmatter(content string) Frontmatter {
	var fm Frontmatter
	inDrawer := false
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue // a gap between keyword groups is not the end of frontmatter
		}
		// A leading drawer (:PROPERTIES: / :DOC-CONFIG: … :END:) sits above the
		// keywords on gptel- and config-style notes; its rows are arbitrary text.
		if inDrawer {
			if strings.EqualFold(line, ":END:") {
				inDrawer = false
			}
			continue
		}
		if !strings.EqualFold(line, ":END:") && drawerOpenRe.MatchString(line) {
			inDrawer = true
			continue
		}
		if !isFrontmatterLine(line) {
			break
		}
		switch {
		case strings.HasPrefix(line, "#+title:"):
			fm.Title = strings.TrimSpace(line[8:])
		case strings.HasPrefix(line, "#+date:"):
			fm.Date = strings.TrimSpace(line[7:])
		case strings.HasPrefix(line, "#+filetags:"):
			raw := strings.TrimSpace(line[11:])
			raw = strings.Trim(raw, ":")
			if raw != "" {
				fm.Filetags = strings.Split(raw, ":")
			}
		case strings.HasPrefix(line, "#+identifier:"):
			fm.Identifier = strings.TrimSpace(line[13:])
		case strings.HasPrefix(line, "#+description:"):
			fm.Description = strings.TrimSpace(line[14:])
		case strings.HasPrefix(line, "#+hugo_lastmod:"):
			fm.HugoLastmod = strings.TrimSpace(line[15:])
		}
	}
	return fm
}

// ReadFrontmatterFile opens a file and parses only its frontmatter region.
// Reads at most maxLines lines or until the first heading/blank-line terminator
// in ParseFrontmatter, whichever comes first. Use this for bulk passes (e.g. day
// notes_modified) to avoid loading entire org files into memory.
// Returns empty Frontmatter on any error.
func ReadFrontmatterFile(path string, maxLines int) Frontmatter {
	f, err := os.Open(path)
	if err != nil {
		return Frontmatter{}
	}
	defer f.Close()

	if maxLines <= 0 {
		maxLines = 50
	}
	var sb strings.Builder
	scanner := bufio.NewScanner(f)
	// Frontmatter lines should be short; default token size is fine, but bump for safety.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	count := 0
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteByte('\n')
		count++
		if count >= maxLines {
			break
		}
	}
	return ParseFrontmatter(sb.String())
}

// ExtractLinks finds all [[denote:ID]] links in content.
func ExtractLinks(content string) []string {
	matches := linkRe.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var links []string
	for _, m := range matches {
		id := m[1]
		if !seen[id] {
			seen[id] = true
			links = append(links, id)
		}
	}
	return links
}

// drawerOpenRe matches a drawer opener (`:PROPERTIES:`, `:DOC-CONFIG:`) but not
// its `:END:` terminator.
var drawerOpenRe = regexp.MustCompile(`(?i)^:(?:[A-Za-z][A-Za-z0-9_-]*):$`)

// fileLocalRe matches an Emacs file-local variables line, which may open a file
// with no leading `#`.
var fileLocalRe = regexp.MustCompile(`^-\*-.*-\*-$`)

var calloutRe = regexp.MustCompile(`^\[!([A-Za-z]+)\]\s*(.*)$`)

// ExtractAbstract returns the first org quote block that opens with a callout
// marker (`[!abstract] 이 노트에 대하여`), searching only the region before the
// first heading. Returns nil when the note carries no such block.
func ExtractAbstract(content string) *Abstract {
	lines := strings.Split(content, "\n")
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "*") && len(trimmed) > 1 && trimmed[1] == ' ' {
			return nil // reached the first heading — no leading callout
		}
		if !strings.EqualFold(trimmed, "#+begin_quote") {
			continue
		}
		// Collect until #+end_quote
		var block []string
		closed := false
		for j := i + 1; j < len(lines); j++ {
			bt := strings.TrimSpace(lines[j])
			if strings.EqualFold(bt, "#+end_quote") {
				closed = true
				i = j
				break
			}
			block = append(block, lines[j])
		}
		if !closed || len(block) == 0 {
			return nil
		}
		// First non-empty line must be the callout marker.
		k := 0
		for k < len(block) && strings.TrimSpace(block[k]) == "" {
			k++
		}
		if k >= len(block) {
			continue
		}
		m := calloutRe.FindStringSubmatch(strings.TrimSpace(block[k]))
		if m == nil {
			continue // a plain quote — keep looking
		}
		body := strings.TrimSpace(strings.Join(block[k+1:], "\n"))
		return &Abstract{
			Kind:  strings.ToLower(m[1]),
			Title: strings.TrimSpace(m[2]),
			Body:  body,
		}
	}
	return nil
}

// applyFrontmatter overlays parsed frontmatter onto a DenoteFile. Filename-derived
// values stay as fallbacks; header values win when present. Lastmod is normalized
// to YYYY-MM-DD while HugoLastmod keeps the raw timestamp (callers that compare
// against commit times need the HH:MM the normalization drops).
func applyFrontmatter(f *DenoteFile, fm Frontmatter) {
	if fm.Title != "" {
		f.Title = fm.Title
	}
	if fm.Date != "" {
		f.Date = fm.Date
	}
	if len(fm.Filetags) > 0 {
		f.Tags = fm.Filetags
	}
	if fm.Description != "" {
		f.Description = fm.Description
	}
	if fm.HugoLastmod != "" {
		f.HugoLastmod = fm.HugoLastmod
		f.Lastmod = ExtractDate(fm.HugoLastmod)
	}
}

// isFrontmatterLine reports whether a non-empty, trimmed line may still belong
// to the frontmatter region. Keyword lines and org comments may; a heading, a
// block opener, or plain prose ends the region. Property drawers are handled by
// the caller, which skips them wholesale.
func isFrontmatterLine(line string) bool {
	switch {
	case line[0] == '*':
		return false // heading
	case strings.HasPrefix(strings.ToLower(line), "#+begin_"):
		return false // block opener — everything inside is body
	case strings.HasPrefix(line, "#+"):
		return true // keyword
	case strings.HasPrefix(line, "# "), line == "#":
		return true // org comment
	case fileLocalRe.MatchString(line):
		return true // `-*- mode: Org; … -*-` file-local variables line
	}
	return false // body text
}
