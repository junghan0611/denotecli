// rename_tag_test.go
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupRenameTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	notes := filepath.Join(dir, "notes")
	os.MkdirAll(notes, 0755)

	files := []struct {
		name    string
		content string
	}{
		{"20250101T100000--테스트-노트__apples_emacs_test.org",
			"#+title: 테스트 노트\n#+filetags: :apples:emacs:test:\n#+identifier: 20250101T100000\n\n* 본문\n"},
		{"20250102T100000--두번째-노트__apples_vim.org",
			"#+title: 두번째 노트\n#+filetags: :apples:vim:\n#+identifier: 20250102T100000\n\n* 내용\n"},
		{"20250103T100000--다른-노트__emacs_config.org",
			"#+title: 다른 노트\n#+filetags: :emacs:config:\n#+identifier: 20250103T100000\n\n* 설정\n"},
	}

	for _, f := range files {
		os.WriteFile(filepath.Join(notes, f.name), []byte(f.content), 0644)
	}
	return dir
}

func TestRenameTagDryRun(t *testing.T) {
	dir := setupRenameTestDir(t)
	files := ScanDirs([]string{dir})

	result, err := RenameTag(files, "apples", "apple", true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Modified != 2 {
		t.Errorf("expected 2 modified, got %d", result.Modified)
	}
	if !result.DryRun {
		t.Error("should be dry run")
	}

	// Verify files NOT actually changed
	files2 := ScanDirs([]string{dir})
	for _, f := range files2 {
		if hasTag(f.Tags, "apple") {
			t.Error("dry run should not modify files")
		}
	}
}

func TestRenameTagActual(t *testing.T) {
	dir := setupRenameTestDir(t)
	files := ScanDirs([]string{dir})

	result, err := RenameTag(files, "apples", "apple", false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Modified != 2 {
		t.Errorf("expected 2 modified, got %d", result.Modified)
	}

	// Re-scan and verify
	files2 := ScanDirs([]string{dir})
	appleCount := 0
	applesCount := 0
	for _, f := range files2 {
		for _, tag := range f.Tags {
			if tag == "apple" {
				appleCount++
			}
			if tag == "apples" {
				applesCount++
			}
		}
	}
	if appleCount != 2 {
		t.Errorf("expected 2 files with 'apple', got %d", appleCount)
	}
	if applesCount != 0 {
		t.Errorf("expected 0 files with 'apples', got %d", applesCount)
	}

	// Verify frontmatter also updated
	for _, f := range files2 {
		if hasTag(f.Tags, "apple") {
			data, _ := os.ReadFile(f.Path)
			content := string(data)
			if strings.Contains(content, ":apples:") {
				t.Errorf("frontmatter still has :apples: in %s", f.Path)
			}
			if !strings.Contains(content, ":apple:") {
				t.Errorf("frontmatter missing :apple: in %s", f.Path)
			}
		}
	}
}

func TestRenameTagSorted(t *testing.T) {
	dir := setupRenameTestDir(t)
	files := ScanDirs([]string{dir})

	// Rename 'apples' to 'aaa' — should sort before other tags
	result, err := RenameTag(files, "apples", "aaa", false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Modified != 2 {
		t.Fatal("expected 2")
	}

	// Verify tags are sorted in filename
	files2 := ScanDirs([]string{dir})
	for _, f := range files2 {
		if hasTag(f.Tags, "aaa") {
			// First tag should be 'aaa'
			if f.Tags[0] != "aaa" {
				t.Errorf("tags not sorted in filename: %v", f.Tags)
			}
		}
	}
}

func TestRenameTagNoMatch(t *testing.T) {
	dir := setupRenameTestDir(t)
	files := ScanDirs([]string{dir})

	result, err := RenameTag(files, "nonexistent", "something", false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Modified != 0 {
		t.Errorf("expected 0 modified, got %d", result.Modified)
	}
}

func TestRenameTagMerge(t *testing.T) {
	// If file already has newTag, rename should merge (not duplicate)
	dir := t.TempDir()
	notes := filepath.Join(dir, "notes")
	os.MkdirAll(notes, 0755)
	os.WriteFile(filepath.Join(notes, "20250101T100000--note__apple_apples.org"),
		[]byte("#+title: note\n#+filetags: :apple:apples:\n#+identifier: 20250101T100000\n"), 0644)

	files := ScanDirs([]string{dir})
	result, err := RenameTag(files, "apples", "apple", false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Modified != 1 {
		t.Fatalf("expected 1, got %d", result.Modified)
	}

	// Should have apple once, not twice
	files2 := ScanDirs([]string{dir})
	for _, f := range files2 {
		count := 0
		for _, tag := range f.Tags {
			if tag == "apple" {
				count++
			}
		}
		if count != 1 {
			t.Errorf("apple count = %d, want 1. tags: %v", count, f.Tags)
		}
	}
}

func TestRewriteFilenameTags(t *testing.T) {
	cases := []struct {
		filename string
		oldTag   string
		newTag   string
		want     string
	}{
		{"20250101T100000--note__apples_emacs.org", "apples", "apple", "20250101T100000--note__apple_emacs.org"},
		{"20250101T100000--note__vim.org", "vim", "neovim", "20250101T100000--note__neovim.org"},
		{"20250101T100000--note.org", "test", "foo", "20250101T100000--note.org"}, // no tags section
	}
	for _, c := range cases {
		got := rewriteFilenameTags(c.filename, c.oldTag, c.newTag)
		if got != c.want {
			t.Errorf("rewriteFilenameTags(%q, %q, %q) = %q, want %q", c.filename, c.oldTag, c.newTag, got, c.want)
		}
	}
}

func TestRewriteFiletags(t *testing.T) {
	content := "#+title: Test\n#+filetags: :apples:emacs:test:\n#+identifier: 123\n"
	got := rewriteFiletags(content, "apples", "apple")
	if !strings.Contains(got, ":apple:") {
		t.Error("missing :apple:")
	}
	if strings.Contains(got, ":apples:") {
		t.Error("still has :apples:")
	}
	// Should be sorted
	if !strings.Contains(got, ":apple:emacs:test:") {
		t.Errorf("not sorted: %s", got)
	}
}
