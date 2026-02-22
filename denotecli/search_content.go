// search_content.go
package main

import (
	"os"
	"strings"
)

// ContentMatch represents a matched line in a file.
type ContentMatch struct {
	Line    int    `json:"line"`
	Snippet string `json:"snippet"`
}

// ContentResult represents a file with matched content lines.
type ContentResult struct {
	DenoteFile
	Matches []ContentMatch `json:"matches"`
}

// SearchContent searches inside file content across all files.
func SearchContent(files []DenoteFile, query string, max int, matchesPerFile int) []ContentResult {
	words := splitWords(strings.ToLower(query))
	if len(words) == 0 {
		return nil
	}

	var results []ContentResult

	for _, f := range files {
		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		var matches []ContentMatch

		for i, line := range lines {
			if matchAllWords(strings.ToLower(line), words) {
				snippet := line
				if len(snippet) > 200 {
					snippet = snippet[:200] + "..."
				}
				matches = append(matches, ContentMatch{
					Line:    i + 1,
					Snippet: strings.TrimSpace(snippet),
				})
				if matchesPerFile > 0 && len(matches) >= matchesPerFile {
					break
				}
			}
		}

		if len(matches) > 0 {
			results = append(results, ContentResult{
				DenoteFile: f,
				Matches:    matches,
			})
			if len(results) >= max {
				break
			}
		}
	}
	return results
}
