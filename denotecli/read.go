// read.go
package main

import (
	"fmt"
	"os"
	"strings"
)

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
