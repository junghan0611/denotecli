// tag_suggest.go
package main

import (
	"sort"
	"strings"
)

// TagCluster represents a group of tags sharing the same stem.
type TagCluster struct {
	Stem    string        `json:"stem"`
	Tags    []TagWithCount `json:"tags"`
	Total   int           `json:"total"`
}

// TagWithCount is a tag with its file count.
type TagWithCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// TagSuggestResult holds tag cleanup suggestions.
type TagSuggestResult struct {
	TotalTags    int          `json:"total_tags"`
	TotalClusters int         `json:"total_clusters"`
	Clusters     []TagCluster `json:"clusters"`
}

// SuggestTagCleanup groups tags by stem to find duplicates.
func SuggestTagCleanup(files []DenoteFile) TagSuggestResult {
	// Count tags
	counts := make(map[string]int)
	for _, f := range files {
		for _, t := range f.Tags {
			counts[t]++
		}
	}

	// Group by stem
	stemGroups := make(map[string][]TagWithCount)
	for tag, count := range counts {
		stem := Stem(tag)
		stemGroups[stem] = append(stemGroups[stem], TagWithCount{Name: tag, Count: count})
	}

	// Only keep clusters with 2+ tags (these are the duplicates)
	var clusters []TagCluster
	for stem, tags := range stemGroups {
		if len(tags) < 2 {
			continue
		}
		// Sort tags within cluster by count desc
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Count > tags[j].Count
		})
		total := 0
		for _, t := range tags {
			total += t.Count
		}
		clusters = append(clusters, TagCluster{
			Stem:  stem,
			Tags:  tags,
			Total: total,
		})
	}

	// Sort clusters by total count desc (highest impact first)
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Total > clusters[j].Total
	})

	// Filter out clusters where all tags share obvious compound prefixes
	// (e.g. emacs/emacslisp are intentionally different)
	var filtered []TagCluster
	for _, c := range clusters {
		if isIntentionalCluster(c) {
			continue
		}
		filtered = append(filtered, c)
	}

	return TagSuggestResult{
		TotalTags:     len(counts),
		TotalClusters: len(filtered),
		Clusters:      filtered,
	}
}

// isIntentionalCluster checks if a cluster is likely intentional (not a mistake).
// Tags that are clearly compound words (emacs + emacslisp) are not duplicates.
func isIntentionalCluster(c TagCluster) bool {
	if len(c.Tags) != 2 {
		return false
	}
	a, b := c.Tags[0].Name, c.Tags[1].Name
	// If one fully contains the other AND the suffix is a meaningful word (>3 chars)
	if strings.HasPrefix(b, a) && len(b)-len(a) > 3 {
		return true
	}
	if strings.HasPrefix(a, b) && len(a)-len(b) > 3 {
		return true
	}
	return false
}
