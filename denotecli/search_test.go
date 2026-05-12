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
		{"notes", "20260101T120000--다단계-문서-테스트__test_outline.org",
			"#+title: 다단계 문서 테스트\n#+filetags: :test:outline:\n#+identifier: 20260101T120000\n\n* 1장 서론\n서론 본문\n** 1.1 배경\n배경 설명\n** 1.2 목적 :IMPORTANT:\n목적 설명\n* 2장 본론\n본론 시작\n** 2.1 방법론\n방법론 설명\n*** 2.1.1 세부 방법\n세부 내용\n** 2.2 결과\n결과 설명\n* 3장 결론\n결론 본문\n[[denote:20251107T082610]]\n"},
		// Bug 4 fixture: slug is short Korean, #+title: has richer name variants.
		// "KimSungdo" and "책과삶" are ONLY in the header — slug-only search misses them.
		{"bib", "20240620T132208--소쉬르__linguistics_semiotics.org",
			"#+title: @페르디낭드소쉬르 @FerdinanddeSaussure 언어의표면 @김성도 @KimSungdo 책과삶\n" +
				"#+filetags: :linguistics:semiotics:\n#+identifier: 20240620T132208\n\n* 본문\n소쉬르 노트\n"},
		// Bug 1 fixture: #+filetags: has :fleeting: but filename has no tag slot.
		// --tags fleeting should find this after the filetags union fix.
		{"notes", "20260201T090000--fleeting-note-without-filename-tag.org",
			"#+title: 플리팅 노트 파일명에 태그 없음\n" +
				"#+filetags: :fleeting:notes:\n#+identifier: 20260201T090000\n\n* 내용\n임시 메모\n"},
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
	if len(files) != 8 {
		t.Fatalf("expected 8 denote files, got %d", len(files))
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

// Bug 4 regression: search must find words present in #+title: but absent from filename slug.
// Forensic: ~/sync/org/llmlog/20260512T102541 — 결함 4번 (search filename-slot-only).
func TestSearchHeaderTitleUnion(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// "KimSungdo" is only in #+title: of the 소쉬르 fixture; not in slug.
	results := Search(files, "KimSungdo", "", false, 20)
	if len(results) != 1 {
		t.Fatalf("Bug 4: expected 1 hit for KimSungdo (header-only word), got %d", len(results))
	}
	if results[0].ID != "20240620T132208" {
		t.Errorf("Bug 4: wrong ID: %q", results[0].ID)
	}

	// "책과삶" also only in #+title:.
	results2 := Search(files, "책과삶", "", false, 20)
	if len(results2) != 1 {
		t.Fatalf("Bug 4: expected 1 hit for 책과삶 (header-only word), got %d", len(results2))
	}

	// HeaderTitle field must be populated.
	if results[0].HeaderTitle == "" {
		t.Error("Bug 4: HeaderTitle field is empty — not populated from #+title:")
	}

	// Slug-based search still works (regression guard).
	results3 := Search(files, "소쉬르", "", false, 20)
	if len(results3) != 1 {
		t.Fatalf("Bug 4 regression: slug-based search broke, expected 1, got %d", len(results3))
	}
}

// Bug 1 regression: --tags filter must find notes whose #+filetags: contains the tag
// even when the filename tag slot does NOT.
// Forensic: ~/sync/org/llmlog/20260512T102541 — 결함 1번 (filetags-header union).
func TestSearchTagsFiletagsUnion(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// The "fleeting-note-without-filename-tag" fixture has :fleeting: only in #+filetags:
	// (filename has no __tag slot at all). After the fix it must appear in --tags fleeting.
	results := Search(files, "", "fleeting", false, 20)
	if len(results) != 1 {
		t.Fatalf("Bug 1: expected 1 hit for --tags fleeting (header-only tag), got %d", len(results))
	}
	if results[0].ID != "20260201T090000" {
		t.Errorf("Bug 1: wrong ID: %q", results[0].ID)
	}

	// Tags union must include both filename and header tags.
	hasNotes := false
	hasFleeting := false
	for _, tag := range results[0].Tags {
		if tag == "notes" {
			hasNotes = true
		}
		if tag == "fleeting" {
			hasFleeting = true
		}
	}
	if !hasNotes || !hasFleeting {
		t.Errorf("Bug 1: Tags union missing; notes=%v fleeting=%v; got: %v", hasNotes, hasFleeting, results[0].Tags)
	}
}
