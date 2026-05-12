// rename_tag.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RenameTagResult holds the result of a tag rename operation.
type RenameTagResult struct {
	OldTag   string   `json:"old_tag"`
	NewTag   string   `json:"new_tag"`
	Modified int      `json:"modified"`
	Files    []string `json:"files"`
	DryRun   bool     `json:"dry_run"`
}

// RenameTag replaces oldTag with newTag across all files (filename + frontmatter).
func RenameTag(files []DenoteFile, oldTag, newTag string, dryRun bool) (RenameTagResult, error) {
	if oldTag == "" || newTag == "" {
		return RenameTagResult{}, fmt.Errorf("old_tag and new_tag required")
	}
	oldTag = strings.ToLower(oldTag)
	newTag = strings.ToLower(newTag)

	modified := make([]string, 0)

	for _, f := range files {
		if !hasTag(f.Tags, oldTag) {
			continue
		}

		if dryRun {
			modified = append(modified, f.Path)
			continue
		}

		// 1. Update frontmatter in file content
		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		content := string(data)
		newContent := rewriteFiletags(content, oldTag, newTag)
		if newContent != content {
			os.WriteFile(f.Path, []byte(newContent), 0644)
		}

		// 2. Rename file (update tags in filename)
		dir := filepath.Dir(f.Path)
		base := filepath.Base(f.Path)
		newBase := rewriteFilenameTags(base, oldTag, newTag)
		if newBase != base {
			newPath := filepath.Join(dir, newBase)
			if err := os.Rename(f.Path, newPath); err != nil {
				return RenameTagResult{}, fmt.Errorf("rename %s → %s: %w", base, newBase, err)
			}
			modified = append(modified, newPath)
		} else {
			modified = append(modified, f.Path)
		}
	}

	return RenameTagResult{
		OldTag:   oldTag,
		NewTag:   newTag,
		Modified: len(modified),
		Files:    modified,
		DryRun:   dryRun,
	}, nil
}

// rewriteFiletags replaces a tag in #+filetags: line.
func rewriteFiletags(content, oldTag, newTag string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "#+filetags:") {
			continue
		}
		raw := strings.TrimSpace(trimmed[11:])
		raw = strings.Trim(raw, ":")
		if raw == "" {
			continue
		}
		tags := strings.Split(raw, ":")
		var newTags []string
		replaced := false
		for _, t := range tags {
			if t == oldTag {
				// Only add newTag if not already present
				if !replaced {
					newTags = append(newTags, newTag)
					replaced = true
				}
			} else if t == newTag {
				// newTag already exists, skip duplicate
				if !replaced {
					newTags = append(newTags, t)
					replaced = true
				}
			} else {
				newTags = append(newTags, t)
			}
		}
		if !replaced {
			continue
		}
		sort.Strings(newTags)
		// Preserve original indentation
		prefix := line[:len(line)-len(trimmed)]
		lines[i] = prefix + "#+filetags:   :" + strings.Join(newTags, ":") + ":"
		break
	}
	return strings.Join(lines, "\n")
}

// rewriteFilenameTags replaces a tag in the filename's tag section.
func rewriteFilenameTags(filename, oldTag, newTag string) string {
	// Pattern: YYYYMMDDTHHMMSS--slug__tag1_tag2_tag3.org
	idx := strings.Index(filename, "__")
	if idx < 0 {
		return filename
	}
	ext := filepath.Ext(filename)
	tagPart := filename[idx+2 : len(filename)-len(ext)]
	tags := strings.Split(tagPart, "_")

	var newTags []string
	replaced := false
	for _, t := range tags {
		if t == oldTag {
			if !replaced {
				newTags = append(newTags, newTag)
				replaced = true
			}
		} else if t == newTag {
			if !replaced {
				newTags = append(newTags, t)
				replaced = true
			}
		} else {
			newTags = append(newTags, t)
		}
	}
	if !replaced {
		return filename
	}
	sort.Strings(newTags)
	return filename[:idx] + "__" + strings.Join(newTags, "_") + ext
}
