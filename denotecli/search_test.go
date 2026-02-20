// search_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir creates a temp directory with fake Denote files.
func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	files := []struct {
		subdir  string
		name    string
		content string
	}{
		{"notes", "20251107T082610--에릭-호퍼-방랑자의-철학__philosophy_creativity.org",
			"#+title: 에릭 호퍼: 방랑자의 철학\n#+filetags: :philosophy:creativity:\n#+identifier: 20251107T082610\n\n* 본문\n창조성에 대한 노트\n"},
		{"notes", "20250901T100000--emacs-설정-가이드__emacs_config.org",
			"#+title: Emacs 설정 가이드\n#+filetags: :emacs:config:\n#+identifier: 20250901T100000\n\n* 설정\ndoom emacs 설정\n"},
		{"bib", "20240601T204208--지식-관리-시스템__pkm_knowledge.org",
			"#+title: 지식 관리 시스템\n#+filetags: :pkm:knowledge:\n#+identifier: 20240601T204208\n\n* PKM\n개인 지식 관리\n"},
		{"llmlog", "20241127T161109--claude-대화-로그__llmlog.org",
			"#+title: Claude 대화 로그\n#+filetags: :llmlog:\n#+identifier: 20241127T161109\n\n* 대화\nAI 대화 내용\n"},
		{"meta", "20230521T215600--\u2021\u00a0이맥스__emacs_metameta_texteditor_ritual_workflow_productivity.org",
			"#+title: \u2021 이맥스\n#+filetags: :emacs:metameta:texteditor:\n#+identifier: 20230521T215600\n\n* Emacs\n에디터 설정\n"},
		{"notes", "README.md", "# Not a Denote file\n"},
	}

	for _, f := range files {
		subdir := filepath.Join(dir, f.subdir)
		os.MkdirAll(subdir, 0755)
		os.WriteFile(filepath.Join(subdir, f.name), []byte(f.content), 0644)
	}
	return dir
}

func TestScanDir(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})
	if len(files) != 5 {
		t.Fatalf("expected 5 denote files, got %d", len(files))
	}
}

func TestSearchByQuery(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := Search(files, "에릭", "", false, 20)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].ID != "20251107T082610" {
		t.Errorf("ID = %q", results[0].ID)
	}
}

func TestSearchMultiWord(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := Search(files, "에릭 철학", "", false, 20)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
}

func TestSearchByTag(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := Search(files, "", "emacs", false, 20)
	if len(results) != 2 {
		t.Fatalf("expected 2 (emacs-설정-가이드 + ‡이맥스), got %d", len(results))
	}
}

func TestSearchNBSPTitle(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// Search by Korean part of title containing U+00A0
	results := Search(files, "이맥스", "", false, 20)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].ID != "20230521T215600" {
		t.Errorf("ID = %q", results[0].ID)
	}
}

func TestSearchTitleOnly(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// "emacs" is in the title of one file (emacs-설정-가이드)
	// the other emacs file has ‡이맥스 title (no "emacs" in title, only in tags)
	results := Search(files, "emacs", "", true, 20)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
}

func TestSearchMaxResults(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := Search(files, "", "", false, 2)
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}
