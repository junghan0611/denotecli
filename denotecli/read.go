// read.go
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var headingRe = regexp.MustCompile(`^(\*+)\s+(.+)$`)

// ReadByID finds a file by Denote ID and returns its full content.
func ReadByID(files []DenoteFile, id string, offset int, limit int) (DenoteContent, error) {
	for _, f := range files {
		if f.ID == id {
			data, err := os.ReadFile(f.Path)
			if err != nil {
				return DenoteContent{}, fmt.Errorf("read %s: %w", f.Path, err)
			}
			content := string(data)

			// Parse frontmatter for richer metadata
			fm := ParseFrontmatter(content)
			if fm.Title != "" {
				f.Title = fm.Title
			}
			if fm.Date != "" {
				f.Date = fm.Date
			}
			if len(fm.Filetags) > 0 {
				f.Tags = fm.Filetags
			}

			// Apply offset/limit
			if offset > 0 || limit > 0 {
				lines := strings.Split(content, "\n")
				if offset >= len(lines) {
					content = ""
				} else {
					end := len(lines)
					if limit > 0 && offset+limit < end {
						end = offset + limit
					}
					content = strings.Join(lines[offset:end], "\n")
				}
			}

			return DenoteContent{
				DenoteFile: f,
				Content:    content,
				Links:      ExtractLinks(string(data)),
			}, nil
		}
	}
	return DenoteContent{}, fmt.Errorf("not found: %s", id)
}

// ExtractOutline parses org headings from content.
func ExtractOutline(content string) []OutlineEntry {
	var outline []OutlineEntry
	for i, line := range strings.Split(content, "\n") {
		m := headingRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		level := len(m[1])
		title := m[2]
		var tags string
		// Org heading tags: * Heading :tag1:tag2:
		if idx := strings.LastIndex(title, ":"); idx > 0 {
			candidate := strings.TrimSpace(title[strings.LastIndex(strings.TrimRight(title, ":"), ":")+1:])
			_ = candidate
			// Look for org-style tags pattern at end: :TAG1:TAG2:
			trimmed := strings.TrimSpace(title)
			parts := strings.Fields(trimmed)
			if len(parts) > 0 {
				last := parts[len(parts)-1]
				if len(last) > 2 && last[0] == ':' && last[len(last)-1] == ':' {
					tags = last
					title = strings.TrimSpace(trimmed[:len(trimmed)-len(last)])
				}
			}
		}
		outline = append(outline, OutlineEntry{
			Level: level,
			Title: title,
			Line:  i + 1, // 1-indexed
			Tags:  tags,
		})
	}
	return outline
}

// ReadOutlineByID returns just the heading structure of a note.
func ReadOutlineByID(files []DenoteFile, id string) (DenoteOutline, error) {
	for _, f := range files {
		if f.ID == id {
			data, err := os.ReadFile(f.Path)
			if err != nil {
				return DenoteOutline{}, fmt.Errorf("read %s: %w", f.Path, err)
			}
			content := string(data)

			fm := ParseFrontmatter(content)
			if fm.Title != "" {
				f.Title = fm.Title
			}
			if fm.Date != "" {
				f.Date = fm.Date
			}
			if len(fm.Filetags) > 0 {
				f.Tags = fm.Filetags
			}

			return DenoteOutline{
				DenoteFile: f,
				Outline:    ExtractOutline(content),
				Links:      ExtractLinks(content),
			}, nil
		}
	}
	return DenoteOutline{}, fmt.Errorf("not found: %s", id)
}
