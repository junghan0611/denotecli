// parser_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFilename(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		wantID        string
		wantSignature string
		wantTitle     string
		wantTags      []string
		wantOK        bool
	}{
		{
			name:      "standard file",
			filename:  "20251107T082610--제목-하이픈-구분__tag1_tag2_tag3.org",
			wantID:    "20251107T082610",
			wantTitle: "제목-하이픈-구분",
			wantTags:  []string{"tag1", "tag2", "tag3"},
			wantOK:    true,
		},
		{
			name:      "no tags",
			filename:  "20251021T105353--simple-title.org",
			wantID:    "20251021T105353",
			wantTitle: "simple-title",
			wantTags:  nil,
			wantOK:    true,
		},
		{
			name:      "korean title with tags",
			filename:  "20250314T152111--에릭-호퍼-방랑자의-철학__philosophy_creativity.org",
			wantID:    "20250314T152111",
			wantTitle: "에릭-호퍼-방랑자의-철학",
			wantTags:  []string{"philosophy", "creativity"},
			wantOK:    true,
		},
		{
			name:      "single tag",
			filename:  "20241127T161109--llm-대화-로그__llmlog.org",
			wantID:    "20241127T161109",
			wantTitle: "llm-대화-로그",
			wantTags:  []string{"llmlog"},
			wantOK:    true,
		},
		{
			name:      "NO-BREAK SPACE in title (U+00A0)",
			filename:  "20230521T215600--\u2021\u00a0이맥스__emacs_metameta_texteditor_ritual_workflow_productivity.org",
			wantID:    "20230521T215600",
			wantTitle: "\u2021\u00a0이맥스",
			wantTags:  []string{"emacs", "metameta", "texteditor", "ritual", "workflow", "productivity"},
			wantOK:    true,
		},
		{
			name:          "signature with tags",
			filename:      "20250904T075937==5a2--힣-ai-에이전트-편재성-기억-연결__agents_ai.org",
			wantID:        "20250904T075937",
			wantSignature: "5a2",
			wantTitle:     "힣-ai-에이전트-편재성-기억-연결",
			wantTags:      []string{"agents", "ai"},
			wantOK:        true,
		},
		{
			name:          "signature without tags",
			filename:      "20250424T225036==0za--†-운명-소명-사명.org",
			wantID:        "20250424T225036",
			wantSignature: "0za",
			wantTitle:     "†-운명-소명-사명",
			wantTags:      nil,
			wantOK:        true,
		},
		{
			name:          "signature with many tags",
			filename:      "20230825T162600==3--†-닷파일-설정파일__configuration_dotfiles_install_meta.org",
			wantID:        "20230825T162600",
			wantSignature: "3",
			wantTitle:     "†-닷파일-설정파일",
			wantTags:      []string{"configuration", "dotfiles", "install", "meta"},
			wantOK:        true,
		},
		{
			name:     "not a denote file",
			filename: "README.md",
			wantOK:   false,
		},
		{
			name:     "org but not denote",
			filename: "AGENTS.md",
			wantOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df, ok := ParseFilename(tt.filename)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if df.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", df.ID, tt.wantID)
			}
			if df.Signature != tt.wantSignature {
				t.Errorf("Signature = %q, want %q", df.Signature, tt.wantSignature)
			}
			if df.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", df.Title, tt.wantTitle)
			}
			if len(df.Tags) != len(tt.wantTags) {
				t.Fatalf("Tags = %v, want %v", df.Tags, tt.wantTags)
			}
			for i, tag := range df.Tags {
				if tag != tt.wantTags[i] {
					t.Errorf("Tags[%d] = %q, want %q", i, tag, tt.wantTags[i])
				}
			}
		})
	}
}

func TestParseFrontmatter(t *testing.T) {
	content := `#+title:      에릭 호퍼: 방랑자의 철학
#+date:       [2025-03-14 Fri 15:21]
#+filetags:   :philosophy:creativity:wisdom:
#+identifier: 20250314T152111
#+description: 에릭 호퍼의 사상과 창조성에 대한 노트

* 히스토리
- [2025-03-14] 초안 작성

* 본문
내용...
`
	fm := ParseFrontmatter(content)
	if fm.Title != "에릭 호퍼: 방랑자의 철학" {
		t.Errorf("Title = %q", fm.Title)
	}
	if fm.Date != "[2025-03-14 Fri 15:21]" {
		t.Errorf("Date = %q", fm.Date)
	}
	if fm.Identifier != "20250314T152111" {
		t.Errorf("Identifier = %q", fm.Identifier)
	}
	wantTags := []string{"philosophy", "creativity", "wisdom"}
	if len(fm.Filetags) != len(wantTags) {
		t.Fatalf("Filetags = %v, want %v", fm.Filetags, wantTags)
	}
	for i, tag := range fm.Filetags {
		if tag != wantTags[i] {
			t.Errorf("Filetags[%d] = %q, want %q", i, tag, wantTags[i])
		}
	}
}

func TestExtractLinks(t *testing.T) {
	content := `* 관련 메타
- [[denote:20240601T204208][링크 제목1]]
- [[denote:20250314T152111][링크 제목2]]
- [[https://example.com][외부 링크]]

* 본문
참고: [[denote:20241127T161109]]
`
	links := ExtractLinks(content)
	want := []string{"20240601T204208", "20250314T152111", "20241127T161109"}
	if len(links) != len(want) {
		t.Fatalf("links = %v, want %v", links, want)
	}
	for i, l := range links {
		if l != want[i] {
			t.Errorf("links[%d] = %q, want %q", i, l, want[i])
		}
	}
}

func BenchmarkScanDirs(b *testing.B) {
	home, err := os.UserHomeDir()
	if err != nil {
		b.Skip("cannot get home dir")
	}
	dir := filepath.Join(home, "org")
	if _, err := os.Stat(dir); err != nil {
		b.Skip("~/org not found")
	}
	for i := 0; i < b.N; i++ {
		files := ScanDirs([]string{dir})
		if len(files) < 1000 {
			b.Fatalf("expected >1000 files, got %d", len(files))
		}
	}
}
