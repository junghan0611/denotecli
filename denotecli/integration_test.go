// integration_test.go
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupIntegrationDir creates a realistic mini knowledge base.
func setupIntegrationDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	dirs := []string{"notes", "bib", "journal", "llmlog", "meta"}
	for _, d := range dirs {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}

	files := []struct {
		subdir  string
		name    string
		content string
	}{
		// Meta note with keyword mapping
		{"meta", "20230521T215600--‡-이맥스__emacs_metameta_texteditor_workflow.org",
			"#+title: ‡ #이맥스\n#+filetags: :emacs:metameta:texteditor:workflow:\n#+identifier: 20230521T215600\n\n" +
				"* 관련메타\n- [[denote:20250101T100000][Emacs 설정 가이드]]\n- [[denote:20250201T100000][둠이맥스 설정]]\n"},

		// Notes
		{"notes", "20250101T100000--emacs-설정-가이드__emacs_config.org",
			"#+title: Emacs 설정 가이드\n#+filetags: :emacs:config:\n#+identifier: 20250101T100000\n\n" +
				"* 기본 설정\n기본적인 이맥스 설정을 다룹니다.\n** use-package\nuse-package로 패키지를 관리합니다.\n" +
				"** 키바인딩\n키바인딩 설정 방법\n* 고급 설정\n** LSP 연동\nlsp-mode 설정 방법\n[[denote:20250201T100000]]\n"},
		{"notes", "20250201T100000--둠이맥스-설정__doomemacs_emacs_config.org",
			"#+title: 둠이맥스 설정\n#+filetags: :doomemacs:emacs:config:\n#+identifier: 20250201T100000\n\n" +
				"* doom 모듈\n둠이맥스 모듈 구성\n** init.el\n모듈 선언\n** config.el\n설정 커스텀\n" +
				"[[denote:20250101T100000][Emacs 설정 가이드]]\n"},

		// Bib note with reference
		{"bib", "20240601T204208--지식-관리-시스템__bib_pkm_knowledge.org",
			"#+title: 지식 관리 시스템\n#+filetags: :bib:pkm:knowledge:\n#+identifier: 20240601T204208\n" +
				"#+reference: zettelkasten2020\n\n* 제텔카스텐\n루만의 지식 관리 방법론\n" +
				"** 슬립박스\n메모를 연결하는 시스템\n[[denote:20250101T100000]]\n"},

		// Journal
		{"journal", "20250310T000000--2025-03-10__journal.org",
			"#+title: 2025-03-10\n#+filetags: :journal:\n#+identifier: 20250310T000000\n\n" +
				"* 오늘 한 일\n이맥스 설정 정리함\n** 커밋 로그\n둠이맥스 config 업데이트\n"},

		// LLM log
		{"llmlog", "20250315T120000--llm-emacs-질문__llmlog_emacs.org",
			"#+title: LLM Emacs 질문\n#+filetags: :llmlog:emacs:\n#+identifier: 20250315T120000\n\n" +
				"* 대화 :LLMLOG:\n** @user emacs에서 LSP 설정 어떻게 하나요?\n** @assistant lsp-mode를 설치하세요.\n"},
	}

	for _, f := range files {
		os.WriteFile(filepath.Join(dir, f.subdir, f.name), []byte(f.content), 0644)
	}
	return dir
}

// Scenario 1: Agent searches for a topic, narrows down, reads
func TestWorkflowSearchToRead(t *testing.T) {
	dir := setupIntegrationDir(t)
	files := ScanDirs([]string{dir})

	// Step 1: search by title
	results := Search(files, "이맥스", "", false, 20)
	if len(results) == 0 {
		t.Fatal("search found nothing")
	}

	// Step 2: read outline of first result
	do, err := ReadOutlineByID(files, results[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(do.Outline) == 0 {
		t.Fatal("outline empty")
	}

	// Step 3: read specific section using line from outline
	firstHeading := do.Outline[0]
	dc, err := ReadByID(files, results[0].ID, firstHeading.Line-1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if dc.Content == "" {
		t.Fatal("content empty")
	}
}

// Scenario 2: Agent uses keyword-map to translate, then searches
func TestWorkflowKeywordMapToSearch(t *testing.T) {
	dir := setupIntegrationDir(t)
	files := ScanDirs([]string{dir})

	// Step 1: user says "이맥스" → find English tag
	km := BuildKeywordMap(files, "이맥스")
	if km.TotalEntries == 0 {
		t.Fatal("keyword-map found nothing")
	}
	engTag := km.Entries[0].Tags[0] // should be "emacs" or similar

	// Step 2: search by English tag
	results := Search(files, "", engTag, false, 20)
	if len(results) < 2 {
		t.Fatalf("expected >= 2 emacs-tagged notes, got %d", len(results))
	}
}

// Scenario 3: Agent searches headings, then reads that section
func TestWorkflowHeadingSearch(t *testing.T) {
	dir := setupIntegrationDir(t)
	files := ScanDirs([]string{dir})

	// Step 1: search headings for "LSP"
	headings := SearchHeadings(files, "LSP", "", 0, 20)
	if len(headings) == 0 {
		t.Fatal("no headings found for LSP")
	}

	// Step 2: read around that heading
	h := headings[0]
	dc, err := ReadByID(files, h.ID, h.Heading.Line-1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(dc.Content), "lsp") {
		t.Errorf("content doesn't contain LSP: %q", dc.Content)
	}
}

// Scenario 4: Agent searches content with tag filter
func TestWorkflowContentSearchFiltered(t *testing.T) {
	dir := setupIntegrationDir(t)
	files := ScanDirs([]string{dir})

	// Search "설정" in bib notes only
	bibResults := SearchContent(files, "설정", "bib", 20, 3)
	// Search "설정" in all notes
	allResults := SearchContent(files, "설정", "", 20, 3)

	if len(bibResults) > len(allResults) {
		t.Error("filtered should be <= all")
	}
	// bib results should only contain bib-tagged files
	for _, r := range bibResults {
		if !hasTag(r.Tags, "bib") {
			t.Errorf("non-bib file in filtered results: %s", r.ID)
		}
	}
}

// Scenario 5: Agent explores graph from a note
func TestWorkflowGraphExploration(t *testing.T) {
	dir := setupIntegrationDir(t)
	files := ScanDirs([]string{dir})

	// Start from Emacs 설정 가이드
	graph, err := BuildGraph(files, "20250101T100000", 1)
	if err != nil {
		t.Fatal(err)
	}

	// Should have outgoing to 둠이맥스
	if len(graph.Outgoing) == 0 {
		t.Fatal("no outgoing links")
	}
	// Should have incoming from meta, 둠이맥스, bib, journal
	if len(graph.Incoming) < 2 {
		t.Fatalf("expected >= 2 incoming, got %d", len(graph.Incoming))
	}
}

// Scenario 6: Agent creates a note, then finds it
func TestWorkflowCreateAndFind(t *testing.T) {
	dir := setupIntegrationDir(t)
	llmlogDir := filepath.Join(dir, "llmlog")

	// Create a new llmlog note
	path, err := CreateNote(llmlogDir, "새 대화 기록", []string{"llmlog", "emacs"}, "* 대화 :LLMLOG:\n새 내용")
	if err != nil {
		t.Fatal(err)
	}

	// Re-scan and find it
	files := ScanDirs([]string{dir})
	df, ok := ParseFilename(filepath.Base(path))
	if !ok {
		t.Fatal("parse failed")
	}

	results := Search(files, df.ID, "", false, 5)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
}

// Scenario 7: Tag suggest finds issues
func TestWorkflowTagSuggest(t *testing.T) {
	dir := setupIntegrationDir(t)
	files := ScanDirs([]string{dir})

	result := SuggestTagCleanup(files)
	// Should have total_tags > 0
	if result.TotalTags == 0 {
		t.Fatal("no tags found")
	}
}
