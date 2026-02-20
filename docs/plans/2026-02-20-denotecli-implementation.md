# denotecli Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Go CLI that parses Denote filenames and org frontmatter from ~/org/ (3,000+ files), providing search/read/tags for AI agents (JSON-only output).

**Architecture:** Single Go module at `denotecli/`, stdlib only. Regex-based Denote filename parser, filepath.WalkDir for file discovery. Subcommands route via os.Args.

**Tech Stack:** Go 1.21+, stdlib only (encoding/json, regexp, os, path/filepath, strings, fmt, bufio, sort, strconv)

**Critical Edge Case:** Some Denote filenames contain U+00A0 (NO-BREAK SPACE) instead of regular spaces. Example: `20230521T215600--‡\u00a0이맥스__emacs_metameta_texteditor_ritual_workflow_productivity.org`. Go's `strings.Fields()` and `strings.TrimSpace()` treat U+00A0 as whitespace via `unicode.IsSpace()`. The parser regex `(.+?)` handles this correctly, but any string splitting/trimming on titles must NOT use these functions. The Denote ID (`YYYYMMDDTHHMMSS`) is the unique key — if you have the ID, you can always find the file.

---

### Task 1: Go module + Denote filename parser

**Files:**
- Create: `denotecli/go.mod`
- Create: `denotecli/parser.go`
- Create: `denotecli/parser_test.go`

**Step 1: Initialize Go module**

```bash
mkdir -p denotecli && cd denotecli && go mod init github.com/junghan0611/org-mode-skills/denotecli
```

**Step 2: Write parser test**

```go
// parser_test.go
package main

import (
	"testing"
)

func TestParseFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantID   string
		wantTitle string
		wantTags []string
		wantOK   bool
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
```

**Step 3: Run test — expect FAIL**

```bash
cd denotecli && go test -v -run TestParse
```

Expected: compilation error (functions not defined).

**Step 4: Implement parser**

```go
// parser.go
package main

import (
	"regexp"
	"strings"
)

// DenoteFile represents parsed Denote file metadata.
type DenoteFile struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Date  string   `json:"date"`
	Path  string   `json:"path"`
}

// DenoteContent extends DenoteFile with file content.
type DenoteContent struct {
	DenoteFile
	Content string   `json:"content"`
	Links   []string `json:"links"`
}

// Frontmatter holds parsed org-mode frontmatter fields.
type Frontmatter struct {
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Filetags    []string `json:"filetags"`
	Identifier  string   `json:"identifier"`
	Description string   `json:"description"`
}

var denoteRe = regexp.MustCompile(`^(\d{8}T\d{6})--(.+?)(?:__(.+))?\.org$`)
var linkRe = regexp.MustCompile(`\[\[denote:(\d{8}T\d{6})\]`)

// ParseFilename extracts Denote metadata from a filename.
func ParseFilename(filename string) (DenoteFile, bool) {
	m := denoteRe.FindStringSubmatch(filename)
	if m == nil {
		return DenoteFile{}, false
	}

	df := DenoteFile{
		ID:    m[1],
		Title: m[2],
	}

	if m[3] != "" {
		df.Tags = strings.Split(m[3], "_")
	}

	// Derive date from ID: 20251107T082610 → 2025-11-07
	if len(df.ID) >= 8 {
		df.Date = df.ID[0:4] + "-" + df.ID[4:6] + "-" + df.ID[6:8]
	}

	return df, true
}

// ParseFrontmatter extracts org-mode frontmatter from file content.
func ParseFrontmatter(content string) Frontmatter {
	var fm Frontmatter
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '*' {
			break // stop at first heading or blank line after frontmatter
		}
		switch {
		case strings.HasPrefix(line, "#+title:"):
			fm.Title = strings.TrimSpace(line[8:])
		case strings.HasPrefix(line, "#+date:"):
			fm.Date = strings.TrimSpace(line[7:])
		case strings.HasPrefix(line, "#+filetags:"):
			raw := strings.TrimSpace(line[11:])
			raw = strings.Trim(raw, ":")
			if raw != "" {
				fm.Filetags = strings.Split(raw, ":")
			}
		case strings.HasPrefix(line, "#+identifier:"):
			fm.Identifier = strings.TrimSpace(line[13:])
		case strings.HasPrefix(line, "#+description:"):
			fm.Description = strings.TrimSpace(line[14:])
		}
	}
	return fm
}

// ExtractLinks finds all [[denote:ID]] links in content.
func ExtractLinks(content string) []string {
	matches := linkRe.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var links []string
	for _, m := range matches {
		id := m[1]
		if !seen[id] {
			seen[id] = true
			links = append(links, id)
		}
	}
	return links
}
```

**Step 5: Run test — expect PASS**

```bash
cd denotecli && go test -v -run "TestParse|TestExtract"
```

Expected: all PASS.

**Step 6: Commit**

```bash
git add denotecli/ && git commit -m "feat(denotecli): Go module + Denote filename/frontmatter parser with tests"
```

---

### Task 2: Search logic

**Files:**
- Create: `denotecli/search.go`
- Create: `denotecli/search_test.go`

**Step 1: Write search test**

```go
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
		subdir string
		name   string
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
			"#+title: ‡ 이맥스\n#+filetags: :emacs:metameta:texteditor:\n#+identifier: 20230521T215600\n\n* Emacs\n에디터 설정\n"},
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
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
}

func TestSearchTitleOnly(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	// "emacs" is in the title of one file
	results := Search(files, "emacs", "", true, 20)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
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

func TestSearchMaxResults(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	results := Search(files, "", "", false, 2)
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}
```

**Step 2: Run test — expect FAIL**

```bash
cd denotecli && go test -v -run TestScan
```

Expected: compilation error.

**Step 3: Implement search**

```go
// search.go
package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanDirs walks directories and returns all parsed Denote files.
func ScanDirs(dirs []string) []DenoteFile {
	var files []DenoteFile
	for _, dir := range dirs {
		dir = expandHome(dir)
		filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			df, ok := ParseFilename(d.Name())
			if !ok {
				return nil
			}
			df.Path = path
			files = append(files, df)
			return nil
		})
	}
	return files
}

// Search filters DenoteFiles by query, tag, and title-only mode.
// Multiple query words are AND-matched.
func Search(files []DenoteFile, query string, tagFilter string, titleOnly bool, max int) []DenoteFile {
	words := splitWords(query)
	var results []DenoteFile

	for i := range files {
		f := &files[i]

		// Tag filter
		if tagFilter != "" && !hasTag(f.Tags, tagFilter) {
			continue
		}

		// Query filter
		if len(words) > 0 {
			searchable := buildSearchable(f, titleOnly)
			if !matchAllWords(searchable, words) {
				continue
			}
		}

		results = append(results, *f)
		if len(results) >= max {
			break
		}
	}
	return results
}

func splitWords(q string) []string {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	return strings.Fields(strings.ToLower(q))
}

func buildSearchable(f *DenoteFile, titleOnly bool) string {
	if titleOnly {
		return strings.ToLower(f.Title)
	}
	var b strings.Builder
	b.WriteString(strings.ToLower(f.ID))
	b.WriteByte(' ')
	b.WriteString(strings.ToLower(f.Title))
	for _, tag := range f.Tags {
		b.WriteByte(' ')
		b.WriteString(strings.ToLower(tag))
	}
	return b.String()
}

func matchAllWords(searchable string, words []string) bool {
	for _, w := range words {
		if !strings.Contains(searchable, w) {
			return false
		}
	}
	return true
}

func hasTag(tags []string, filter string) bool {
	parts := strings.Split(strings.ToLower(filter), ",")
	for _, want := range parts {
		want = strings.TrimSpace(want)
		for _, tag := range tags {
			if strings.ToLower(tag) == want {
				return true
			}
		}
	}
	return false
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
```

**Step 4: Run test — expect PASS**

```bash
cd denotecli && go test -v -run "TestScan|TestSearch"
```

Expected: all PASS.

**Step 5: Commit**

```bash
git add denotecli/ && git commit -m "feat(denotecli): search logic with WalkDir, multi-word AND, tag filter"
```

---

### Task 3: Read + Tags

**Files:**
- Create: `denotecli/read.go`
- Create: `denotecli/tags.go`
- Create: `denotecli/tags_test.go`

**Step 1: Write tags test**

```go
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
			if ts.Count != 1 {
				t.Errorf("emacs count = %d, want 1", ts.Count)
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
```

**Step 2: Run test — expect FAIL**

```bash
cd denotecli && go test -v -run "TestCollect|TestRead"
```

**Step 3: Implement read.go**

```go
// read.go
package main

import (
	"fmt"
	"os"
	"strings"
)

// ReadByID finds a file by Denote ID and returns its full content.
func ReadByID(files []DenoteFile, id string, offset int, limit int) (DenoteContent, error) {
	for _, f := range files {
		if f.ID == id {
			data, err := os.ReadFile(f.Path)
			if err != nil {
				return DenoteContent{}, fmt.Errorf("read %s: %w", f.Path, err)
			}
			content := string(data)

			// Parse frontmatter for richer metadata
			fm := ParseFrontmatter(content)
			if fm.Title != "" {
				f.Title = fm.Title
			}
			if fm.Date != "" {
				f.Date = fm.Date
			}
			if len(fm.Filetags) > 0 {
				f.Tags = fm.Filetags
			}

			// Apply offset/limit
			if offset > 0 || limit > 0 {
				lines := strings.Split(content, "\n")
				if offset >= len(lines) {
					content = ""
				} else {
					end := len(lines)
					if limit > 0 && offset+limit < end {
						end = offset + limit
					}
					content = strings.Join(lines[offset:end], "\n")
				}
			}

			return DenoteContent{
				DenoteFile: f,
				Content:    content,
				Links:      ExtractLinks(string(data)),
			}, nil
		}
	}
	return DenoteContent{}, fmt.Errorf("not found: %s", id)
}
```

**Step 4: Implement tags.go**

```go
// tags.go
package main

import (
	"regexp"
	"sort"
	"strings"
)

// TagStat holds a tag name and its count.
type TagStat struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// TagsResult holds the tag statistics output.
type TagsResult struct {
	TotalFiles int       `json:"total_files"`
	TotalTags  int       `json:"total_tags"`
	Tags       []TagStat `json:"tags"`
}

// CollectTags aggregates tag counts from all files.
func CollectTags(files []DenoteFile, pattern string, top int) TagsResult {
	counts := make(map[string]int)
	for _, f := range files {
		for _, tag := range f.Tags {
			counts[strings.ToLower(tag)]++
		}
	}

	var patRe *regexp.Regexp
	if pattern != "" {
		patRe, _ = regexp.Compile(pattern)
	}

	var tags []TagStat
	for name, count := range counts {
		if patRe != nil && !patRe.MatchString(name) {
			continue
		}
		tags = append(tags, TagStat{Name: name, Count: count})
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Count > tags[j].Count
	})

	if top > 0 && len(tags) > top {
		tags = tags[:top]
	}

	return TagsResult{
		TotalFiles: len(files),
		TotalTags:  len(counts),
		Tags:       tags,
	}
}
```

**Step 5: Run test — expect PASS**

```bash
cd denotecli && go test -v -run "TestCollect|TestRead"
```

**Step 6: Commit**

```bash
git add denotecli/ && git commit -m "feat(denotecli): read by ID + tag statistics with tests"
```

---

### Task 4: CLI main + JSON output

**Files:**
- Create: `denotecli/main.go`

**Step 1: Implement main.go**

```go
// main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "search":
		cmdSearch()
	case "read":
		cmdRead()
	case "tags":
		cmdTags()
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func cmdSearch() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli search <query> [--tags TAG] [--dirs DIR,...] [--title-only] [--max N]")
	}
	query := os.Args[2]
	args := os.Args[3:]
	tagFilter := getFlag(args, "--tags", "")
	dirsStr := getFlag(args, "--dirs", "~/org")
	titleOnly := hasFlag(args, "--title-only")
	maxStr := getFlag(args, "--max", "20")
	max, _ := strconv.Atoi(maxStr)
	if max <= 0 {
		max = 20
	}

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	results := Search(files, query, tagFilter, titleOnly, max)
	printJSON(results)
}

func cmdRead() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli read <id> [--dirs DIR,...] [--offset N] [--limit N]")
	}
	id := os.Args[2]
	args := os.Args[3:]
	dirsStr := getFlag(args, "--dirs", "~/org")
	offsetStr := getFlag(args, "--offset", "0")
	limitStr := getFlag(args, "--limit", "0")
	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	dc, err := ReadByID(files, id, offset, limit)
	if err != nil {
		fatal(err.Error())
	}
	printJSON(dc)
}

func cmdTags() {
	args := os.Args[2:]
	dirsStr := getFlag(args, "--dirs", "~/org")
	pattern := getFlag(args, "--pattern", "")
	topStr := getFlag(args, "--top", "50")
	top, _ := strconv.Atoi(topStr)
	if top <= 0 {
		top = 50
	}

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	stats := CollectTags(files, pattern, top)
	printJSON(stats)
}

func getFlag(args []string, name string, def string) string {
	for i, arg := range args {
		if arg == name && i+1 < len(args) {
			return args[i+1]
		}
	}
	return def
}

func hasFlag(args []string, name string) bool {
	for _, arg := range args {
		if arg == name {
			return true
		}
	}
	return false
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	enc.Encode(v)
}

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, `denotecli - Denote knowledge base CLI for AI agents

Usage:
  denotecli search <query> [--tags TAG] [--dirs DIR,...] [--title-only] [--max N]
  denotecli read <id> [--dirs DIR,...] [--offset N] [--limit N]
  denotecli tags [--pattern PAT] [--top N] [--dirs DIR,...]

Options:
  --dirs DIR,...    Search directories, comma-separated (default: ~/org)
  --tags TAG        Filter by tag (comma-separated)
  --title-only      Search title only (search command)
  --max N           Max results (default: 20)
  --offset N        Start line (read command)
  --limit N         Lines to read (read command, 0=all)
  --pattern PAT     Tag name regex filter (tags command)
  --top N           Top N tags (default: 50)
`)
}
```

**Step 2: Build and test manually**

```bash
cd denotecli && go build -o denotecli .
```

```bash
./denotecli tags --dirs ~/org --top 10
./denotecli search "에릭 호퍼" --dirs ~/org
./denotecli search "emacs" --dirs ~/org --tags emacs --max 5
./denotecli read 20250314T152111 --dirs ~/org
./denotecli --help
```

**Step 3: Run all tests**

```bash
cd denotecli && go test -v ./...
```

**Step 4: Commit**

```bash
git add denotecli/ && git commit -m "feat(denotecli): CLI main with search/read/tags commands, JSON output"
```

---

### Task 5: Repo cleanup + SKILL.md + benchmark

**Files:**
- Move: legacy files → `docs/archive/`
- Create: `SKILL.md` (pi-skills format, replace existing)
- Modify: `README.md` (concise rewrite)
- Modify: `.gitignore`

**Step 1: Archive legacy files**

```bash
mkdir -p docs/archive
mv scripts/ docs/archive/
mv skill-creator/ docs/archive/
mv references/ docs/archive/
mv temp/ docs/archive/
mv tmp/ docs/archive/
mv .claude-plugin/ docs/archive/
mv README-KO.md docs/archive/README-KO-old.md
mv TODO.org docs/archive/
mv README.md docs/archive/README-old.md
mv SKILL.md docs/archive/SKILL-old.md
```

**Step 2: Write SKILL.md**

```markdown
---
name: denote-org
description: "Denote knowledge base CLI for searching, reading, and analyzing 3,000+ org-mode files. Use when working with ~/org/, Denote files (YYYYMMDDTHHMMSS--title__tags.org), or org-mode knowledge bases."
---

# Denote-Org Skill

CLI tool for navigating Denote/org-mode knowledge bases. Binary at `{baseDir}/denotecli/denotecli`.

## Commands

### Search notes

\`\`\`bash
{baseDir}/denotecli/denotecli search "<query>" [--tags TAG] [--dirs DIR,...] [--title-only] [--max N]
\`\`\`

- Multi-word queries use AND matching
- Searches filename (ID, title, tags)
- `--tags`: filter by tag (comma-separated)
- `--dirs`: search directories (default: ~/org)
- `--title-only`: match title only
- `--max`: limit results (default: 20)

### Read a note

\`\`\`bash
{baseDir}/denotecli/denotecli read <denote-id> [--dirs DIR,...] [--offset N] [--limit N]
\`\`\`

- Returns full content + parsed metadata + outgoing denote links
- `--offset`/`--limit`: for large files

### Tag statistics

\`\`\`bash
{baseDir}/denotecli/denotecli tags [--pattern PAT] [--top N] [--dirs DIR,...]
\`\`\`

- Aggregates tags from all Denote filenames
- `--pattern`: regex filter
- `--top`: limit (default: 50)

## Denote File Format

- **Filename**: `YYYYMMDDTHHMMSS--title-with-hyphens__tag1_tag2.org`
- **Frontmatter**: `#+title:`, `#+date:`, `#+filetags:`, `#+identifier:`
- **Links**: `[[denote:YYYYMMDDTHHMMSS]]`

## Knowledge Base Structure

| Directory | Purpose |
|-----------|---------|
| `~/org/notes/` | Main notes |
| `~/org/bib/` | Bibliography |
| `~/org/journal/` | Journals |
| `~/org/llmlog/` | LLM conversation logs |

## All output is JSON

Every command outputs JSON to stdout. Errors go to stderr.

## Build

\`\`\`bash
cd {baseDir}/denotecli && go build -o denotecli .
\`\`\`
```

**Step 3: Write README.md**

```markdown
# orgmode-skills

Denote knowledge base CLI (`denotecli`) for AI agents. Searches, reads, and analyzes 3,000+ org-mode files.

## Quick Start

\`\`\`bash
cd denotecli && go build -o denotecli .
./denotecli search "에릭 호퍼" --dirs ~/org
./denotecli read 20250314T152111 --dirs ~/org
./denotecli tags --dirs ~/org --top 10
\`\`\`

## As Claude Code Skill

\`\`\`bash
ln -s /path/to/orgmode-skills ~/.claude/skills/denote-org
\`\`\`

## License

Apache 2.0
```

**Step 4: Update .gitignore**

```
denotecli/denotecli
```

**Step 5: Benchmark with real data**

```bash
cd denotecli && go test -bench=. -benchtime=3s -v
```

Add benchmark to parser_test.go:

```go
func BenchmarkScanDirs(b *testing.B) {
	dir := os.ExpandEnv("$HOME/org")
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
```

Expected: < 100ms per scan for 3,000+ files.

**Step 6: Commit all cleanup**

```bash
git add -A && git commit -m "refactor: rebuild repo with Go denotecli, archive legacy files"
```

**Step 7: Final verification**

```bash
cd denotecli
./denotecli search "emacs" --dirs ~/org --max 5
./denotecli search "창조" --dirs ~/org --tags philosophy
./denotecli read 20251107T082610 --dirs ~/org
./denotecli tags --dirs ~/org --top 10
./denotecli --help
```
