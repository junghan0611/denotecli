// create.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

// CreateNote creates a new Denote note with proper naming and header.
func CreateNote(dir string, title string, tags []string, content string) (string, error) {
	dir = expandHome(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", dir, err)
	}

	now := time.Now()
	id := now.Format("20060102T150405")
	date := now.Format("[2006-01-02 Mon 15:04]")

	slug := slugify(title)
	sort.Strings(tags)
	tagSuffix := ""
	if len(tags) > 0 {
		tagSuffix = "__" + strings.Join(tags, "_")
	}

	filename := fmt.Sprintf("%s--%s%s.org", id, slug, tagSuffix)
	path := filepath.Join(dir, filename)

	// Build header
	var header strings.Builder
	fmt.Fprintf(&header, "#+title:      %s\n", title)
	fmt.Fprintf(&header, "#+date:       %s\n", date)
	if len(tags) > 0 {
		fmt.Fprintf(&header, "#+filetags:   :%s:\n", strings.Join(tags, ":"))
	}
	fmt.Fprintf(&header, "#+identifier: %s\n", id)
	fmt.Fprintf(&header, "#+export_file_name: %s.md\n", id)

	body := header.String()
	if content != "" {
		body += "\n" + content + "\n"
	}

	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return path, nil
}

var nonSlugRe = regexp.MustCompile(`[^\p{L}\p{N}]+`)

// slugify converts title to Denote-style slug (lowercase, hyphens).
func slugify(title string) string {
	s := strings.ToLower(title)
	// Replace non-letter/non-digit runs with hyphens
	s = nonSlugRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	// Collapse multiple hyphens
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return s
}

// sanitizeTag ensures tag is valid (lowercase alphanum+hyphen).
func sanitizeTag(tag string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(tag) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
