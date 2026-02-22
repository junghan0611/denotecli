// keyword_map_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func setupMetaTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	meta := filepath.Join(dir, "meta")
	os.MkdirAll(meta, 0755)

	files := []struct {
		name    string
		content string
	}{
		{"20230723T061300--†-깃허브-리포-저장소-마깃__action_github_magit_meta_version.org",
			"#+title: † #깃허브 #리포 #저장소 #마깃\n#+filetags: :action:github:magit:meta:version:\n#+identifier: 20230723T061300\n\n* 관련\n"},
		{"20230521T215600--‡-이맥스__emacs_metameta_productivity_texteditor_workflow.org",
			"#+title: ‡ #이맥스\n#+filetags: :emacs:metameta:productivity:texteditor:workflow:\n#+identifier: 20230521T215600\n\n* Emacs\n"},
		{"20220519T101300--†-통계-확률__bib_inference_meta_theory_statistics_probability.org",
			"#+title: † #통계 #확률\n#+filetags: :bib:inference:meta:theory:statistics:probability:\n#+identifier: 20220519T101300\n\n* 통계\n"},
	}

	for _, f := range files {
		os.WriteFile(filepath.Join(meta, f.name), []byte(f.content), 0644)
	}

	// Non-meta file (should be excluded)
	notes := filepath.Join(dir, "notes")
	os.MkdirAll(notes, 0755)
	os.WriteFile(filepath.Join(notes, "20251107T082610--에릭-호퍼__philosophy.org"),
		[]byte("#+title: 에릭 호퍼\n#+filetags: :philosophy:\n#+identifier: 20251107T082610\n"), 0644)

	return dir
}

func TestBuildKeywordMapAll(t *testing.T) {
	dir := setupMetaTestDir(t)
	files := ScanDirs([]string{dir})

	result := BuildKeywordMap(files, "")
	if result.TotalEntries < 5 {
		t.Errorf("expected >= 5 entries, got %d", result.TotalEntries)
	}

	// Verify 깃허브 → github mapping exists
	found := false
	for _, e := range result.Entries {
		if e.Keyword == "깃허브" {
			found = true
			hasGithub := false
			for _, tag := range e.Tags {
				if tag == "github" {
					hasGithub = true
				}
			}
			if !hasGithub {
				t.Error("깃허브 should map to github tag")
			}
		}
	}
	if !found {
		t.Error("깃허브 keyword not found")
	}
}

func TestBuildKeywordMapQuery(t *testing.T) {
	dir := setupMetaTestDir(t)
	files := ScanDirs([]string{dir})

	// Search by Korean keyword
	result := BuildKeywordMap(files, "깃허브")
	if result.TotalEntries != 1 {
		t.Errorf("expected 1, got %d", result.TotalEntries)
	}

	// Search by English tag
	result = BuildKeywordMap(files, "emacs")
	if result.TotalEntries != 1 {
		t.Errorf("expected 1 (이맥스→emacs), got %d", result.TotalEntries)
	}
}

func TestBuildKeywordMapExcludesNonMeta(t *testing.T) {
	dir := setupMetaTestDir(t)
	files := ScanDirs([]string{dir})

	result := BuildKeywordMap(files, "호퍼")
	if result.TotalEntries != 0 {
		t.Errorf("non-meta file should be excluded, got %d", result.TotalEntries)
	}
}

func TestBuildKeywordMapNBSP(t *testing.T) {
	// Title with NO-BREAK SPACE should split keywords correctly
	dir := t.TempDir()
	meta := filepath.Join(dir, "meta")
	os.MkdirAll(meta, 0755)
	// U+00A0 between #이맥스 and #지식그래프
	os.WriteFile(filepath.Join(meta, "20230223T053200--†-이맥스-지식그래프__ekg_emacs_knowledgegraph_meta.org"),
		[]byte("#+title: † #이맥스\u00a0#지식그래프 §ekg\n#+filetags: :ekg:emacs:knowledgegraph:meta:\n#+identifier: 20230223T053200\n"), 0644)

	files := ScanDirs([]string{dir})
	result := BuildKeywordMap(files, "")

	// Should be 2 separate keywords, not 1 merged
	if result.TotalEntries != 2 {
		t.Errorf("expected 2 keywords, got %d", result.TotalEntries)
		for _, e := range result.Entries {
			t.Logf("  keyword=%q", e.Keyword)
		}
	}
}

func TestBuildKeywordMapNoResults(t *testing.T) {
	dir := setupMetaTestDir(t)
	files := ScanDirs([]string{dir})

	result := BuildKeywordMap(files, "존재하지않는키워드xyz")
	if result.TotalEntries != 0 {
		t.Errorf("expected 0, got %d", result.TotalEntries)
	}
}
