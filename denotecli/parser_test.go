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
#+hugo_lastmod: [2025-03-29 Sat 02:06]

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
	if fm.HugoLastmod != "[2025-03-29 Sat 02:06]" {
		t.Errorf("HugoLastmod = %q", fm.HugoLastmod)
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

func TestExtractDate(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"[2025-06-10]", "2025-06-10"},
		{"[2025-03-29 Sat 02:06]", "2025-03-29"},
		{"<2024-01-03 Wed 16:54>", "2024-01-03"},
		{"2023-06-19", "2023-06-19"},
		{"Time-stamp: <2025-01-31 12:54:06 junghan>", "2025-01-31"},
		{"no date here", ""},
	}
	for _, tt := range tests {
		if got := ExtractDate(tt.in); got != tt.want {
			t.Errorf("ExtractDate(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseFrontmatterHugoLastmod(t *testing.T) {
	cases := []struct {
		name        string
		content     string
		wantLastmod string
		wantDate    string // ExtractDate output
	}{
		{
			name:        "bracket date only",
			content:     "#+title: t\n#+identifier: 20250610T120000\n#+hugo_lastmod: [2025-06-10]\n",
			wantLastmod: "[2025-06-10]",
			wantDate:    "2025-06-10",
		},
		{
			name:        "bracket with day and time",
			content:     "#+title: t\n#+identifier: 20250329T020600\n#+hugo_lastmod: [2025-03-29 Sat 02:06]\n",
			wantLastmod: "[2025-03-29 Sat 02:06]",
			wantDate:    "2025-03-29",
		},
		{
			name:        "angle active timestamp",
			content:     "#+title: t\n#+identifier: 20240103T165400\n#+hugo_lastmod: <2024-01-03 Wed 16:54>\n",
			wantLastmod: "<2024-01-03 Wed 16:54>",
			wantDate:    "2024-01-03",
		},
		{
			name:        "bare date",
			content:     "#+title: t\n#+identifier: 20230619T000000\n#+hugo_lastmod: 2023-06-19\n",
			wantLastmod: "2023-06-19",
			wantDate:    "2023-06-19",
		},
		{
			name:        "emacs time-stamp form",
			content:     "#+title: t\n#+identifier: 20250131T125400\n#+hugo_lastmod: Time-stamp: <2025-01-31 12:54:06 junghan>\n",
			wantLastmod: "Time-stamp: <2025-01-31 12:54:06 junghan>",
			wantDate:    "2025-01-31",
		},
		{
			name:        "missing field",
			content:     "#+title: t\n#+identifier: 20250101T000000\n",
			wantLastmod: "",
			wantDate:    "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fm := ParseFrontmatter(tc.content)
			if fm.HugoLastmod != tc.wantLastmod {
				t.Errorf("HugoLastmod = %q, want %q", fm.HugoLastmod, tc.wantLastmod)
			}
			if got := ExtractDate(fm.HugoLastmod); got != tc.wantDate {
				t.Errorf("ExtractDate(%q) = %q, want %q", fm.HugoLastmod, got, tc.wantDate)
			}
		})
	}
}

func TestParseFrontmatterIgnoresBodyHugoLastmod(t *testing.T) {
	// A line that looks like #+hugo_lastmod inside the body (after first heading)
	// must NOT be parsed as a frontmatter value.
	content := `#+title:      t
#+identifier: 20260101T000000

* 본문
#+hugo_lastmod: [2099-12-31]
`
	fm := ParseFrontmatter(content)
	if fm.HugoLastmod != "" {
		t.Errorf("HugoLastmod leaked from body = %q", fm.HugoLastmod)
	}
}

func TestReadFrontmatterFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "20260508T100000--sample.org")
	body := `#+title:      sample
#+date:       [2026-05-08 Fri 10:00]
#+identifier: 20260508T100000
#+hugo_lastmod: [2026-05-08 Fri 22:24]

* 본문
real content goes here.
`
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	fm := ReadFrontmatterFile(path, 50)
	if fm.Title != "sample" {
		t.Errorf("Title = %q", fm.Title)
	}
	if got := ExtractDate(fm.HugoLastmod); got != "2026-05-08" {
		t.Errorf("ExtractDate = %q, want 2026-05-08", got)
	}

	// Missing file -> empty Frontmatter, no panic.
	missing := ReadFrontmatterFile(filepath.Join(tmp, "nope.org"), 50)
	if missing.Title != "" || missing.HugoLastmod != "" {
		t.Errorf("missing file should yield zero Frontmatter, got %+v", missing)
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
