// create_test.go
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSlugify(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"에릭 호퍼: 방랑자의 철학", "에릭-호퍼-방랑자의-철학"},
		{"Emacs 설정 가이드", "emacs-설정-가이드"},
		{"#LLM: 대화 로그", "llm-대화-로그"},
		{"simple", "simple"},
		{"  spaces  ", "spaces"},
		{"a--b--c", "a-b-c"},
	}
	for _, c := range cases {
		got := slugify(c.input)
		if got != c.want {
			t.Errorf("slugify(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestSanitizeTag(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"llmlog", "llmlog"},
		{"LLM-Log", "llmlog"},
		{"tag_1", "tag1"},
		{"한글태그", "한글태그"},
	}
	for _, c := range cases {
		got := sanitizeTag(c.input)
		if got != c.want {
			t.Errorf("sanitizeTag(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestCreateNote(t *testing.T) {
	dir := t.TempDir()

	path, err := CreateNote(dir, "테스트 노트 생성", []string{"llmlog", "test"}, "* 본문 내용\n테스트입니다.")
	if err != nil {
		t.Fatal(err)
	}

	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not found: %s", path)
	}

	// Verify filename pattern
	base := filepath.Base(path)
	if !strings.HasSuffix(base, ".org") {
		t.Errorf("not .org: %s", base)
	}
	if !strings.Contains(base, "--테스트-노트-생성__") {
		t.Errorf("missing slug: %s", base)
	}
	if !strings.Contains(base, "__llmlog_test.org") {
		t.Errorf("tags not sorted: %s", base)
	}

	// Verify content
	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "#+title:      테스트 노트 생성") {
		t.Error("missing title")
	}
	if !strings.Contains(content, "#+filetags:   :llmlog:test:") {
		t.Error("missing filetags")
	}
	if !strings.Contains(content, "#+identifier:") {
		t.Error("missing identifier")
	}
	if !strings.Contains(content, "#+export_file_name:") {
		t.Error("missing export_file_name")
	}
	if !strings.Contains(content, "* 본문 내용") {
		t.Error("missing body content")
	}

	// Verify parseable by our own parser
	df, ok := ParseFilename(base)
	if !ok {
		t.Fatalf("ParseFilename failed: %s", base)
	}
	if df.Title != "테스트-노트-생성" {
		t.Errorf("Title = %q", df.Title)
	}
}

func TestCreateNoteNoTags(t *testing.T) {
	dir := t.TempDir()

	path, err := CreateNote(dir, "태그 없는 노트", nil, "")
	if err != nil {
		t.Fatal(err)
	}

	base := filepath.Base(path)
	if strings.Contains(base, "__") {
		t.Errorf("should have no tags: %s", base)
	}

	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "#+filetags:") {
		t.Error("should have no filetags line")
	}
}

func TestCreateNoteSubdir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "sub", "dir")

	path, err := CreateNote(dir, "서브디렉토리", []string{"test"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("file not created in subdir")
	}
}

func TestCreateNoteEmptyContent(t *testing.T) {
	dir := t.TempDir()

	path, err := CreateNote(dir, "빈 노트", []string{"test"}, "")
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	// Should end with export_file_name line, no extra blank lines
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	last := lines[len(lines)-1]
	if !strings.HasPrefix(last, "#+export_file_name:") {
		t.Errorf("last line = %q, want export_file_name", last)
	}
}
