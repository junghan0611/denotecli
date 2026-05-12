// search_content_test.go
package main

import (
	"testing"
)

func TestSearchContentBasic(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := SearchContent(files, "창조성", "", 20, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 file, got %d", len(results))
	}
	if results[0].ID != "20251107T082610" {
		t.Errorf("ID = %q", results[0].ID)
	}
	if len(results[0].Matches) != 1 {
		t.Errorf("matches = %d", len(results[0].Matches))
	}
	if results[0].Matches[0].Line != 6 {
		t.Errorf("line = %d, want 6", results[0].Matches[0].Line)
	}
}

func TestSearchContentMultiWord(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// "doom emacs" appears on one line in emacs config file
	results := SearchContent(files, "doom emacs", "", 20, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 file, got %d", len(results))
	}
	if results[0].ID != "20250901T100000" {
		t.Errorf("ID = %q", results[0].ID)
	}
}

func TestSearchContentKorean(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := SearchContent(files, "개인 지식", "", 20, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 file, got %d", len(results))
	}
	if results[0].ID != "20240601T204208" {
		t.Errorf("ID = %q", results[0].ID)
	}
}

func TestSearchContentMaxFiles(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// Empty query → empty slice (NOT nil). Agent callers rely on a valid
	// array contract; nil would marshal to JSON null and break len() calls.
	results := SearchContent(files, "", "", 20, 3)
	if results == nil {
		t.Errorf("expected non-nil empty slice for empty query, got nil")
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty query, got %d", len(results))
	}
}

func TestSearchContentMatchesPerFile(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// Search for something that appears in frontmatter AND body
	// "emacs" appears in filetags line and heading in emacs config file
	results := SearchContent(files, "emacs", "", 20, 1)
	for _, r := range results {
		if len(r.Matches) > 1 {
			t.Errorf("file %s has %d matches, want <= 1", r.ID, len(r.Matches))
		}
	}
}

func TestSearchContentNoResults(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := SearchContent(files, "존재하지않는내용xyz", "", 20, 3)
	if len(results) != 0 {
		t.Errorf("expected 0, got %d", len(results))
	}
}

func TestSearchContentMaxResultsCap(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// Search broad term, cap at 2 files
	results := SearchContent(files, "설정", "", 2, 3)
	if len(results) > 2 {
		t.Errorf("expected <= 2, got %d", len(results))
	}
}

func TestSearchContentTagFilter(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// Without tag filter: "설정" matches in multiple files
	all := SearchContent(files, "설정", "", 20, 3)

	// With tag filter: only emacs-tagged files
	filtered := SearchContent(files, "설정", "emacs", 20, 3)
	if len(filtered) > len(all) {
		t.Errorf("filtered (%d) should be <= all (%d)", len(filtered), len(all))
	}
	for _, r := range filtered {
		found := false
		for _, tag := range r.Tags {
			if tag == "emacs" {
				found = true
			}
		}
		if !found {
			t.Errorf("file %s missing emacs tag", r.ID)
		}
	}
}
