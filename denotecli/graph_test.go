// graph_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func setupGraphTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	notes := filepath.Join(dir, "notes")
	os.MkdirAll(notes, 0755)

	files := []struct {
		name    string
		content string
	}{
		{"20250101T100000--노트-A__test.org",
			"#+title: 노트 A\n#+filetags: :test:\n#+identifier: 20250101T100000\n\n* 본문\n[[denote:20250101T200000][노트 B 링크]]\n[[denote:20250101T300000][노트 C 링크]]\n"},
		{"20250101T200000--노트-B__test.org",
			"#+title: 노트 B\n#+filetags: :test:\n#+identifier: 20250101T200000\n\n* 본문\n[[denote:20250101T100000][노트 A 역링크]]\n"},
		{"20250101T300000--노트-C__test.org",
			"#+title: 노트 C\n#+filetags: :test:\n#+identifier: 20250101T300000\n\n* 본문\n독립 노트\n"},
	}

	for _, f := range files {
		os.WriteFile(filepath.Join(notes, f.name), []byte(f.content), 0644)
	}
	return dir
}

func TestBuildGraphOutgoing(t *testing.T) {
	dir := setupGraphTestDir(t)
	files := ScanDirs([]string{dir})

	result, err := BuildGraph(files, "20250101T100000", 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Outgoing) != 2 {
		t.Fatalf("outgoing: expected 2, got %d", len(result.Outgoing))
	}
	// Should link to B and C
	targets := make(map[string]bool)
	for _, l := range result.Outgoing {
		targets[l.TargetID] = true
	}
	if !targets["20250101T200000"] || !targets["20250101T300000"] {
		t.Errorf("outgoing targets: %v", targets)
	}
}

func TestBuildGraphIncoming(t *testing.T) {
	dir := setupGraphTestDir(t)
	files := ScanDirs([]string{dir})

	// Note A should have incoming from Note B
	result, err := BuildGraph(files, "20250101T100000", 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Incoming) != 1 {
		t.Fatalf("incoming: expected 1, got %d", len(result.Incoming))
	}
	if result.Incoming[0].SourceID != "20250101T200000" {
		t.Errorf("incoming source: %q", result.Incoming[0].SourceID)
	}
}

func TestBuildGraphIsolated(t *testing.T) {
	dir := setupGraphTestDir(t)
	files := ScanDirs([]string{dir})

	// Note C has no outgoing links and incoming from A only
	result, err := BuildGraph(files, "20250101T300000", 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Outgoing) != 0 {
		t.Errorf("outgoing: expected 0, got %d", len(result.Outgoing))
	}
	if len(result.Incoming) != 1 {
		t.Fatalf("incoming: expected 1 (from A), got %d", len(result.Incoming))
	}
}

func TestBuildGraphNotFound(t *testing.T) {
	dir := setupGraphTestDir(t)
	files := ScanDirs([]string{dir})

	_, err := BuildGraph(files, "99999999T999999", 1)
	if err == nil {
		t.Error("expected error")
	}
}

func TestBuildGraphMetadata(t *testing.T) {
	dir := setupGraphTestDir(t)
	files := ScanDirs([]string{dir})

	result, err := BuildGraph(files, "20250101T100000", 1)
	if err != nil {
		t.Fatal(err)
	}

	if result.Title != "노트 A" {
		t.Errorf("Title = %q", result.Title)
	}
	if len(result.Tags) == 0 || result.Tags[0] != "test" {
		t.Errorf("Tags = %v", result.Tags)
	}
}
