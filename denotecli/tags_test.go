// tags_test.go
package main

import (
	"testing"
)

func TestCollectTags(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	stats := CollectTags(files, "", 50)
	if stats.TotalFiles != 5 {
		t.Errorf("TotalFiles = %d, want 5", stats.TotalFiles)
	}
	if stats.TotalTags < 5 {
		t.Errorf("TotalTags = %d, want >= 5", stats.TotalTags)
	}
	// emacs should be in the list
	found := false
	for _, ts := range stats.Tags {
		if ts.Name == "emacs" {
			found = true
			if ts.Count != 2 {
				t.Errorf("emacs count = %d, want 2", ts.Count)
			}
		}
	}
	if !found {
		t.Error("tag 'emacs' not found")
	}
}

func TestCollectTagsPattern(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	stats := CollectTags(files, "ph", 50)
	if len(stats.Tags) != 1 {
		t.Fatalf("expected 1 tag matching 'ph', got %d", len(stats.Tags))
	}
	if stats.Tags[0].Name != "philosophy" {
		t.Errorf("tag = %q, want philosophy", stats.Tags[0].Name)
	}
}

func TestReadByID(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	dc, err := ReadByID(files, "20251107T082610", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if dc.ID != "20251107T082610" {
		t.Errorf("ID = %q", dc.ID)
	}
	if dc.Title != "에릭 호퍼: 방랑자의 철학" {
		t.Errorf("Title = %q", dc.Title)
	}
	if dc.Content == "" {
		t.Error("Content is empty")
	}
}

func TestReadByIDNotFound(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	_, err := ReadByID(files, "99999999T999999", 0, 0)
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}
