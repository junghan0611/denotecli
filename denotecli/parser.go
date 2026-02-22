// parser.go
package main

import (
	"regexp"
	"strings"
)

// DenoteFile represents parsed Denote file metadata.
type DenoteFile struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Date  string   `json:"date"`
	Path  string   `json:"path"`
}

// DenoteContent extends DenoteFile with file content.
type DenoteContent struct {
	DenoteFile
	Content string   `json:"content"`
	Links   []string `json:"links"`
}

// OutlineEntry represents a single org heading.
type OutlineEntry struct {
	Level   int    `json:"level"`
	Title   string `json:"title"`
	Line    int    `json:"line"`
	Tags    string `json:"tags,omitempty"`
}

// DenoteOutline extends DenoteFile with outline (headings only).
type DenoteOutline struct {
	DenoteFile
	Outline []OutlineEntry `json:"outline"`
	Links   []string       `json:"links"`
}

// Frontmatter holds parsed org-mode frontmatter fields.
type Frontmatter struct {
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Filetags    []string `json:"filetags"`
	Identifier  string   `json:"identifier"`
	Description string   `json:"description"`
}

var denoteRe = regexp.MustCompile(`^(\d{8}T\d{6})--(.+?)(?:__(.+))?\.org$`)
var linkRe = regexp.MustCompile(`\[\[denote:(\d{8}T\d{6})\]`)

// ParseFilename extracts Denote metadata from a filename.
func ParseFilename(filename string) (DenoteFile, bool) {
	m := denoteRe.FindStringSubmatch(filename)
	if m == nil {
		return DenoteFile{}, false
	}

	df := DenoteFile{
		ID:    m[1],
		Title: m[2],
	}

	if m[3] != "" {
		df.Tags = strings.Split(m[3], "_")
	}

	// Derive date from ID: 20251107T082610 -> 2025-11-07
	if len(df.ID) >= 8 {
		df.Date = df.ID[0:4] + "-" + df.ID[4:6] + "-" + df.ID[6:8]
	}

	return df, true
}

// ParseFrontmatter extracts org-mode frontmatter from file content.
func ParseFrontmatter(content string) Frontmatter {
	var fm Frontmatter
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || (len(line) > 0 && line[0] == '*') {
			if fm.Title != "" || fm.Identifier != "" {
				break // stop at first heading or blank line after we've seen frontmatter
			}
			continue
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
		}
	}
	return fm
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
