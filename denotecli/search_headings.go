// search_headings.go
package main

import (
	"os"
	"strings"
)

// HeadingMatch represents a matched heading in a file.
type HeadingMatch struct {
	DenoteFile
	Heading OutlineEntry `json:"heading"`
}

// SearchHeadings searches org headings across all files.
// tagFilter filters files by tag before heading search (comma-separated, OR).
func SearchHeadings(files []DenoteFile, query string, tagFilter string, maxLevel int, max int) []HeadingMatch {
	words := splitWords(strings.ToLower(query))
	var results []HeadingMatch

	for _, f := range files {
		if tagFilter != "" && !hasTag(f.Tags, tagFilter) {
			continue
		}
		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		outline := ExtractOutline(string(data))
		for _, entry := range outline {
			if maxLevel > 0 && entry.Level > maxLevel {
				continue
			}
			searchable := strings.ToLower(entry.Title)
			if len(words) == 0 || matchAllWords(searchable, words) {
				results = append(results, HeadingMatch{
					DenoteFile: f,
					Heading:    entry,
				})
				if len(results) >= max {
					return results
				}
			}
		}
	}
	return results
}
