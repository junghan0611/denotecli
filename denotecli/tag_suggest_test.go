// tag_suggest_test.go
package main

import (
	"testing"
)

func TestSuggestTagCleanup(t *testing.T) {
	files := []DenoteFile{
		{ID: "1", Tags: []string{"emacs", "vim", "apple", "apples"}},
		{ID: "2", Tags: []string{"emacs", "apple", "communication", "communicational"}},
		{ID: "3", Tags: []string{"vim", "apples", "note", "notes"}},
	}

	result := SuggestTagCleanup(files)
	if result.TotalClusters == 0 {
		t.Fatal("expected clusters")
	}

	// Should find apple/apples cluster
	found := false
	for _, c := range result.Clusters {
		for _, tag := range c.Tags {
			if tag.Name == "apple" || tag.Name == "apples" {
				found = true
			}
		}
	}
	if !found {
		t.Error("missing apple/apples cluster")
	}
}

func TestSuggestTagCleanupStemBased(t *testing.T) {
	// communicate/communication/communicational should all cluster
	files := []DenoteFile{
		{ID: "1", Tags: []string{"communicate"}},
		{ID: "2", Tags: []string{"communication"}},
		{ID: "3", Tags: []string{"communicational"}},
	}

	result := SuggestTagCleanup(files)
	found := false
	for _, c := range result.Clusters {
		if len(c.Tags) == 3 {
			found = true
		}
	}
	if !found {
		t.Error("expected 3-tag cluster for communicate variants")
	}
}

func TestSuggestTagCleanupIntentional(t *testing.T) {
	// emacs + emacslisp should be filtered as intentional
	files := []DenoteFile{
		{ID: "1", Tags: []string{"emacs"}},
		{ID: "2", Tags: []string{"emacs"}},
		{ID: "3", Tags: []string{"emacslisp"}},
	}

	result := SuggestTagCleanup(files)
	for _, c := range result.Clusters {
		for _, tag := range c.Tags {
			if tag.Name == "emacslisp" {
				t.Error("emacslisp should be filtered as intentional compound")
			}
		}
	}
}
