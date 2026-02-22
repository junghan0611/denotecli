// tag_suggest.go
package main

import (
	"sort"
	"strings"
)

// TagPair represents two similar tags that might be duplicates.
type TagPair struct {
	Tag1  string `json:"tag1"`
	Count1 int   `json:"count1"`
	Tag2  string `json:"tag2"`
	Count2 int   `json:"count2"`
	Reason string `json:"reason"`
}

// TagSuggestResult holds tag cleanup suggestions.
type TagSuggestResult struct {
	TotalTags   int       `json:"total_tags"`
	Suggestions []TagPair `json:"suggestions"`
}

// SuggestTagCleanup finds potentially duplicate/similar tags.
func SuggestTagCleanup(files []DenoteFile) TagSuggestResult {
	// Count tags
	counts := make(map[string]int)
	for _, f := range files {
		for _, t := range f.Tags {
			counts[t]++
		}
	}

	// Collect sorted tag list
	var tags []string
	for t := range counts {
		tags = append(tags, t)
	}
	sort.Strings(tags)

	var suggestions []TagPair

	for i := 0; i < len(tags); i++ {
		for j := i + 1; j < len(tags); j++ {
			a, b := tags[i], tags[j]
			if reason := isSimilar(a, b); reason != "" {
				suggestions = append(suggestions, TagPair{
					Tag1: a, Count1: counts[a],
					Tag2: b, Count2: counts[b],
					Reason: reason,
				})
			}
		}
	}

	// Sort by combined count (higher = more impactful to fix)
	sort.Slice(suggestions, func(i, j int) bool {
		ci := suggestions[i].Count1 + suggestions[i].Count2
		cj := suggestions[j].Count1 + suggestions[j].Count2
		return ci > cj
	})

	return TagSuggestResult{
		TotalTags:   len(tags),
		Suggestions: suggestions,
	}
}

// isSimilar checks if two tags are likely duplicates.
func isSimilar(a, b string) string {
	// Exact plural: tag vs tags
	if a+"s" == b || b+"s" == a {
		return "plural"
	}
	// -es plural: process vs processes
	if a+"es" == b || b+"es" == a {
		return "plural"
	}
	// -ies plural: strategy vs strategies (after stripping)
	if strings.HasSuffix(a, "y") && strings.TrimSuffix(a, "y")+"ies" == b {
		return "plural"
	}
	if strings.HasSuffix(b, "y") && strings.TrimSuffix(b, "y")+"ies" == a {
		return "plural"
	}

	// Common suffix variations: -tion/-tional, -ment/-mental, -ity/-ive
	stems := [][2]string{
		{"tion", "tional"}, {"ment", "mental"},
		{"ity", "ive"}, {"ism", "ist"},
		{"ing", ""}, {"tion", "te"},
		{"ence", "ent"}, {"ance", "ant"},
	}
	for _, pair := range stems {
		if swapSuffix(a, b, pair[0], pair[1]) || swapSuffix(b, a, pair[0], pair[1]) {
			return "derivation"
		}
	}

	// Prefix containment: if one is prefix of other and diff <= 3 chars
	if len(a) > 3 && len(b) > 3 {
		if strings.HasPrefix(b, a) && len(b)-len(a) <= 3 {
			return "prefix"
		}
		if strings.HasPrefix(a, b) && len(a)-len(b) <= 3 {
			return "prefix"
		}
	}

	return ""
}

func swapSuffix(a, b, suf1, suf2 string) bool {
	if strings.HasSuffix(a, suf1) {
		stem := strings.TrimSuffix(a, suf1)
		if stem+suf2 == b {
			return true
		}
	}
	return false
}
