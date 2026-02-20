// tags.go
package main

import (
	"regexp"
	"sort"
	"strings"
)

// TagStat holds a tag name and its count.
type TagStat struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// TagsResult holds the tag statistics output.
type TagsResult struct {
	TotalFiles int       `json:"total_files"`
	TotalTags  int       `json:"total_tags"`
	Tags       []TagStat `json:"tags"`
}

// CollectTags aggregates tag counts from all files.
func CollectTags(files []DenoteFile, pattern string, top int) TagsResult {
	counts := make(map[string]int)
	for _, f := range files {
		for _, tag := range f.Tags {
			counts[strings.ToLower(tag)]++
		}
	}

	var patRe *regexp.Regexp
	if pattern != "" {
		patRe, _ = regexp.Compile(pattern)
	}

	var tags []TagStat
	for name, count := range counts {
		if patRe != nil && !patRe.MatchString(name) {
			continue
		}
		tags = append(tags, TagStat{Name: name, Count: count})
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Count > tags[j].Count
	})

	if top > 0 && len(tags) > top {
		tags = tags[:top]
	}

	return TagsResult{
		TotalFiles: len(files),
		TotalTags:  len(counts),
		Tags:       tags,
	}
}
