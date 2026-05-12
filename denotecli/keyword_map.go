// keyword_map.go
package main

import (
	"os"
	"regexp"
	"strings"
)

// KeywordEntry represents a Korean keyword ↔ English tag mapping from a meta note.
type KeywordEntry struct {
	Keyword string `json:"keyword"`
	Tags    []string `json:"tags"`
	NoteID  string `json:"note_id"`
	Title   string `json:"title"`
}

// KeywordMapResult holds the full keyword map.
type KeywordMapResult struct {
	TotalEntries int            `json:"total_entries"`
	Entries      []KeywordEntry `json:"entries"`
}

var hashKeywordRe = regexp.MustCompile(`#(\S+)`)

// BuildKeywordMap extracts Korean keyword → English tag mappings from meta notes.
// Meta notes use #+title with #한글키워드 and filename tags in English.
func BuildKeywordMap(files []DenoteFile, query string) KeywordMapResult {
	words := splitWords(strings.ToLower(query))
	entries := make([]KeywordEntry, 0)

	for _, f := range files {
		// Only process meta-tagged files
		if !hasTag(f.Tags, "meta") && !hasTag(f.Tags, "metameta") {
			continue
		}

		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		fm := ParseFrontmatter(string(data))
		if fm.Title == "" {
			continue
		}

		// Extract #keyword patterns from title
		// Replace NO-BREAK SPACE (U+00A0) with regular space before matching
		cleanTitle := strings.ReplaceAll(fm.Title, "\u00a0", " ")
		matches := hashKeywordRe.FindAllStringSubmatch(cleanTitle, -1)
		if len(matches) == 0 {
			continue
		}

		// Get non-meta tags from filename (these are the English equivalents)
		var engTags []string
		for _, t := range f.Tags {
			if t != "meta" && t != "metameta" {
				engTags = append(engTags, t)
			}
		}

		for _, m := range matches {
			kw := m[1]
			entry := KeywordEntry{
				Keyword: kw,
				Tags:    engTags,
				NoteID:  f.ID,
				Title:   fm.Title,
			}

			if len(words) == 0 {
				entries = append(entries, entry)
			} else {
				searchable := strings.ToLower(kw + " " + strings.Join(engTags, " "))
				if matchAllWords(searchable, words) {
					entries = append(entries, entry)
				}
			}
		}
	}

	return KeywordMapResult{
		TotalEntries: len(entries),
		Entries:      entries,
	}
}
