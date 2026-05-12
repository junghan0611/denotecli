// search.go
package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanDirs walks directories and returns all parsed Denote files.
func ScanDirs(dirs []string) []DenoteFile {
	var files []DenoteFile
	for _, dir := range dirs {
		dir = expandHome(dir)
		filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			df, ok := ParseFilename(d.Name())
			if !ok {
				return nil
			}
			df.Path = path
			files = append(files, df)
			return nil
		})
	}
	return files
}

// Search filters DenoteFiles by query, tag, and title-only mode.
// Multiple query words are AND-matched.
func Search(files []DenoteFile, query string, tagFilter string, titleOnly bool, max int) []DenoteFile {
	words := splitWords(query)
	// Initialize as empty slice (not nil) so JSON marshals to [] not null.
	// Agent callers rely on a valid array contract even when no matches.
	results := make([]DenoteFile, 0)

	for i := range files {
		f := &files[i]

		// Tag filter
		if tagFilter != "" && !hasTag(f.Tags, tagFilter) {
			continue
		}

		// Query filter
		if len(words) > 0 {
			searchable := buildSearchable(f, titleOnly)
			if !matchAllWords(searchable, words) {
				continue
			}
		}

		results = append(results, *f)
		if len(results) >= max {
			break
		}
	}
	return results
}

func splitWords(q string) []string {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	return strings.Fields(strings.ToLower(q))
}

func buildSearchable(f *DenoteFile, titleOnly bool) string {
	if titleOnly {
		return strings.ToLower(f.Title)
	}
	var b strings.Builder
	b.WriteString(strings.ToLower(f.ID))
	b.WriteByte(' ')
	b.WriteString(strings.ToLower(f.Title))
	for _, tag := range f.Tags {
		b.WriteByte(' ')
		b.WriteString(strings.ToLower(tag))
	}
	return b.String()
}

func matchAllWords(searchable string, words []string) bool {
	for _, w := range words {
		if !strings.Contains(searchable, w) {
			return false
		}
	}
	return true
}

func hasTag(tags []string, filter string) bool {
	parts := strings.Split(strings.ToLower(filter), ",")
	for _, want := range parts {
		want = strings.TrimSpace(want)
		for _, tag := range tags {
			if strings.ToLower(tag) == want {
				return true
			}
		}
	}
	return false
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		path = filepath.Join(home, path[2:])
	}
	// Resolve symlinks so WalkDir traverses the real directory
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return resolved
}
