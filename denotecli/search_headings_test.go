// search_headings_test.go
package main

import (
	"testing"
)

func TestSearchHeadingsBasic(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := SearchHeadings(files, "서론", 0, 20)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Heading.Title != "1장 서론" {
		t.Errorf("Title = %q", results[0].Heading.Title)
	}
	if results[0].ID != "20260101T120000" {
		t.Errorf("ID = %q", results[0].ID)
	}
}

func TestSearchHeadingsMultiWord(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := SearchHeadings(files, "세부 방법", 0, 20)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Heading.Level != 3 {
		t.Errorf("Level = %d, want 3", results[0].Heading.Level)
	}
}

func TestSearchHeadingsLevelFilter(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// "방법" matches "2.1 방법론" (level 2) and "2.1.1 세부 방법" (level 3)
	allResults := SearchHeadings(files, "방법", 0, 20)
	if len(allResults) != 2 {
		t.Fatalf("no level filter: expected 2, got %d", len(allResults))
	}

	// Level 2 filter should exclude level 3
	filteredResults := SearchHeadings(files, "방법", 2, 20)
	if len(filteredResults) != 1 {
		t.Fatalf("level 2 filter: expected 1, got %d", len(filteredResults))
	}
	if filteredResults[0].Heading.Title != "2.1 방법론" {
		t.Errorf("Title = %q", filteredResults[0].Heading.Title)
	}
}

func TestSearchHeadingsAcrossFiles(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// "본문" appears as heading in 에릭호퍼 file, "본론" in 다단계문서
	results := SearchHeadings(files, "본", 0, 20)
	if len(results) < 2 {
		t.Fatalf("expected >= 2, got %d", len(results))
	}
	// Verify they come from different files
	ids := make(map[string]bool)
	for _, r := range results {
		ids[r.ID] = true
	}
	if len(ids) < 2 {
		t.Errorf("expected headings from multiple files, got %d unique IDs", len(ids))
	}
}

func TestSearchHeadingsMaxResults(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := SearchHeadings(files, "", 0, 3)
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
}

func TestSearchHeadingsKorean(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := SearchHeadings(files, "설정", 0, 20)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].ID != "20250901T100000" {
		t.Errorf("ID = %q, want emacs config file", results[0].ID)
	}
}

func TestSearchHeadingsNoResults(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := SearchHeadings(files, "존재하지않는검색어xyz", 0, 20)
	if len(results) != 0 {
		t.Fatalf("expected 0, got %d", len(results))
	}
}

func TestFilterOutlineByLevel(t *testing.T) {
	outline := []OutlineEntry{
		{Level: 1, Title: "H1", Line: 1},
		{Level: 2, Title: "H2", Line: 3},
		{Level: 3, Title: "H3", Line: 5},
		{Level: 1, Title: "H1b", Line: 7},
	}

	// No filter
	r := FilterOutlineByLevel(outline, 0)
	if len(r) != 4 {
		t.Errorf("level 0: expected 4, got %d", len(r))
	}

	// Level 1
	r = FilterOutlineByLevel(outline, 1)
	if len(r) != 2 {
		t.Errorf("level 1: expected 2, got %d", len(r))
	}

	// Level 2
	r = FilterOutlineByLevel(outline, 2)
	if len(r) != 3 {
		t.Errorf("level 2: expected 3, got %d", len(r))
	}
}
