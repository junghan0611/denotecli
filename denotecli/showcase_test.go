// showcase_test.go — Human-readable showcase of all search patterns.
// Run: go test -v -run TestShowcase
// Purpose: see actual search behavior at a glance, spot edge cases visually.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupShowcaseDir(t *testing.T) (string, []DenoteFile) {
	t.Helper()
	dir := t.TempDir()
	subdirs := []string{"notes", "bib", "journal", "llmlog", "meta"}
	for _, d := range subdirs {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}

	testFiles := []struct {
		subdir  string
		name    string
		content string
	}{
		{"notes", "20250101T100000--emacs-설정-가이드__config_emacs.org",
			"#+title: Emacs 설정 가이드\n#+filetags: :config:emacs:\n#+identifier: 20250101T100000\n\n" +
				"* 기본 설정\n이맥스 기본 설정을 다룹니다.\n** use-package\nuse-package로 패키지 관리\n" +
				"** 키바인딩\nC-c C-c 같은 키바인딩\n* 고급 설정\n** LSP 연동\nlsp-mode 설정\n*** eglot\neglot 대안\n"},
		{"notes", "20250201T100000--둠이맥스-설정__config_doomemacs_emacs.org",
			"#+title: 둠이맥스 설정\n#+filetags: :config:doomemacs:emacs:\n#+identifier: 20250201T100000\n\n" +
				"* doom 모듈\n모듈 구성\n** init.el\n(doom! :completion company)\n" +
				"** config.el\n설정 커스텀\n[[denote:20250101T100000]]\n"},
		{"bib", "20240601T204208--지식-관리-시스템__bib_knowledge_pkm.org",
			"#+title: 지식 관리 시스템\n#+filetags: :bib:knowledge:pkm:\n#+identifier: 20240601T204208\n" +
				"#+reference: zettelkasten2020\n\n* 제텔카스텐\n루만의 슬립박스 방법론\n" +
				"** 원자적 노트\n하나의 노트 = 하나의 아이디어\n** 연결\n노트 간 링크가 핵심\n" +
				"[[denote:20250101T100000]]\n"},
		{"bib", "20250314T125213--에릭호퍼-방랑자의-철학__autobiography_bib_philosophy.org",
			"#+title: @에릭호퍼 길 위의 철학자\n#+filetags: :autobiography:bib:philosophy:\n#+identifier: 20250314T125213\n\n" +
				"* 맹신자들 :1951:\n대중운동의 본질에 관한 125가지 단상\n" +
				"** 1부 대중운동의 매력\n변화를 향한 갈망\n" +
				"** 2부 잠재적 전향자\n가난한 사람, 부적응자\n" +
				"* 인간의 조건 :1973:\n창조자, 예언자, 인간\n"},
		{"journal", "20250310T000000--2025-03-10__journal.org",
			"#+title: 2025-03-10\n#+filetags: :journal:\n#+identifier: 20250310T000000\n\n" +
				"* 오늘 한 일\n이맥스 설정 정리, 에릭 호퍼 독서\n* 내일 할 일\n양자역학 노트 정리\n"},
		{"llmlog", "20250315T120000--llm-emacs-질문__emacs_llmlog.org",
			"#+title: LLM Emacs 질문\n#+filetags: :emacs:llmlog:\n#+identifier: 20250315T120000\n\n" +
				"* 대화 :LLMLOG:\n** @user LSP 설정 어떻게 하나요?\nlsp-mode와 eglot 차이가 뭔가요?\n" +
				"** @assistant\nlsp-mode는 기능이 많고 eglot은 가볍습니다.\n"},
		{"meta", "20230521T215600--‡\u00a0이맥스__emacs_metameta_texteditor_workflow.org",
			"#+title: ‡ #이맥스\n#+filetags: :emacs:metameta:texteditor:workflow:\n#+identifier: 20230521T215600\n\n" +
				"* 관련메타\n- [[denote:20250101T100000][Emacs 설정 가이드]]\n- [[denote:20250201T100000][둠이맥스 설정]]\n" +
				"* 도구 철학\n이맥스는 단순 에디터가 아니라 사고 도구\n"},
		{"meta", "20230723T061300--†-깃허브-리포-저장소-마깃__action_github_magit_meta_version.org",
			"#+title: † #깃허브 #리포 #저장소 #마깃\n#+filetags: :action:github:magit:meta:version:\n#+identifier: 20230723T061300\n\n" +
				"* 깃허브 관리\n리포 관리 방법\n"},
	}

	for _, f := range testFiles {
		os.WriteFile(filepath.Join(dir, f.subdir, f.name), []byte(f.content), 0644)
	}

	files := ScanDirs([]string{dir})
	return dir, files
}

func p(format string, args ...interface{}) {
	fmt.Printf("    "+format+"\n", args...)
}

func TestShowcaseSearch(t *testing.T) {
	_, files := setupShowcaseDir(t)

	patterns := []struct {
		query     string
		tags      string
		titleOnly bool
		desc      string
	}{
		{"에릭 호퍼", "", false, "한글 다중단어 검색"},
		{"emacs", "", false, "영어 검색 (제목+태그 모두)"},
		{"emacs", "", true, "영어 title-only (태그 제외)"},
		{"설정", "emacs", false, "한글 + 태그 필터"},
		{"doom", "", false, "영어 단어"},
		{"20250314T125213", "", false, "ID로 직접 검색"},
		{"관리", "bib", false, "bib 태그로 범위 제한"},
		{"zzz없는검색어", "", false, "결과 없음"},
	}

	fmt.Println("\n=== search patterns ===")
	for _, pat := range patterns {
		results := Search(files, pat.query, pat.tags, pat.titleOnly, 20)
		fmt.Printf("  [%s] query=%q tags=%q titleOnly=%v → %d건\n", pat.desc, pat.query, pat.tags, pat.titleOnly, len(results))
		for _, r := range results {
			p("→ %s %s [%s]", r.ID, r.Title, strings.Join(r.Tags, ","))
		}
	}
}

func TestShowcaseSearchHeadings(t *testing.T) {
	_, files := setupShowcaseDir(t)

	patterns := []struct {
		query string
		tags  string
		level int
		desc  string
	}{
		{"설정", "", 0, "한글 헤딩 전체 레벨"},
		{"설정", "emacs", 0, "emacs 태그 필터"},
		{"설정", "", 1, "레벨1만"},
		{"LSP", "", 0, "영어 기술용어"},
		{"대중운동", "", 0, "한글 학술용어"},
		{"도구", "metameta", 0, "meta 노트만"},
	}

	fmt.Println("\n=== search-headings patterns ===")
	for _, pat := range patterns {
		results := SearchHeadings(files, pat.query, pat.tags, pat.level, 20)
		fmt.Printf("  [%s] query=%q tags=%q level=%d → %d건\n", pat.desc, pat.query, pat.tags, pat.level, len(results))
		for _, r := range results {
			p("→ %s L%d \"%s\" (line %d) in %s", r.ID, r.Heading.Level, r.Heading.Title, r.Heading.Line, r.Title)
		}
	}
}

func TestShowcaseSearchContent(t *testing.T) {
	_, files := setupShowcaseDir(t)

	patterns := []struct {
		query string
		tags  string
		desc  string
	}{
		{"이맥스", "", "한글 본문 검색"},
		{"lsp-mode", "", "영어 패키지명"},
		{"eglot", "emacs", "emacs 태그 필터 + 영어"},
		{"갈망", "bib", "bib만 본문 검색"},
		{"에릭 호퍼", "journal", "일지에서만 검색"},
		{"양자역학", "", "본문에만 있는 키워드"},
	}

	fmt.Println("\n=== search-content patterns ===")
	for _, pat := range patterns {
		results := SearchContent(files, pat.query, pat.tags, 20, 3)
		fmt.Printf("  [%s] query=%q tags=%q → %d건\n", pat.desc, pat.query, pat.tags, len(results))
		for _, r := range results {
			for _, m := range r.Matches {
				p("→ %s line %d: %s", r.ID, m.Line, truncate(m.Snippet, 60))
			}
		}
	}
}

func TestShowcaseOutline(t *testing.T) {
	_, files := setupShowcaseDir(t)

	ids := []struct {
		id    string
		level int
		desc  string
	}{
		{"20250314T125213", 0, "에릭호퍼 전체 outline"},
		{"20250314T125213", 1, "에릭호퍼 레벨1만"},
		{"20250101T100000", 0, "Emacs 설정 전체"},
		{"20250315T120000", 0, "LLM 대화 로그"},
	}

	fmt.Println("\n=== read --outline patterns ===")
	for _, item := range ids {
		do, err := ReadOutlineByID(files, item.id)
		if err != nil {
			continue
		}
		outline := FilterOutlineByLevel(do.Outline, item.level)
		fmt.Printf("  [%s] id=%s level=%d → %d headings\n", item.desc, item.id, item.level, len(outline))
		for _, h := range outline {
			indent := strings.Repeat("  ", h.Level-1)
			tags := ""
			if h.Tags != "" {
				tags = " " + h.Tags
			}
			p("%s%s (L%d)%s", indent, h.Title, h.Line, tags)
		}
	}
}

func TestShowcaseKeywordMap(t *testing.T) {
	_, files := setupShowcaseDir(t)

	queries := []struct {
		query string
		desc  string
	}{
		{"", "전체 매핑 덤프"},
		{"이맥스", "한글로 검색"},
		{"emacs", "영어로 검색"},
		{"github", "영어 태그"},
		{"없는키워드", "결과 없음"},
	}

	fmt.Println("\n=== keyword-map patterns ===")
	for _, q := range queries {
		result := BuildKeywordMap(files, q.query)
		fmt.Printf("  [%s] query=%q → %d entries\n", q.desc, q.query, result.TotalEntries)
		for _, e := range result.Entries {
			p("→ #%s ↔ [%s]", e.Keyword, strings.Join(e.Tags, ", "))
		}
	}
}

func TestShowcaseGraph(t *testing.T) {
	_, files := setupShowcaseDir(t)

	ids := []struct {
		id   string
		desc string
	}{
		{"20250101T100000", "Emacs 설정 (허브 노트)"},
		{"20250314T125213", "에릭호퍼 (독립 노트)"},
		{"20250310T000000", "일지 (leaf 노트)"},
	}

	fmt.Println("\n=== graph patterns ===")
	for _, item := range ids {
		result, err := BuildGraph(files, item.id, 1)
		if err != nil {
			continue
		}
		fmt.Printf("  [%s] id=%s → out:%d in:%d\n", item.desc, item.id, len(result.Outgoing), len(result.Incoming))
		for _, l := range result.Outgoing {
			p("→ out: %s", l.TargetID)
		}
		for _, l := range result.Incoming {
			p("← in:  %s (%s)", l.SourceID, l.SourceTitle)
		}
	}
}

func TestShowcaseTagSuggest(t *testing.T) {
	_, files := setupShowcaseDir(t)

	result := SuggestTagCleanup(files)
	fmt.Printf("\n=== tags --suggest ===\n")
	fmt.Printf("  total tags: %d, clusters: %d\n", result.TotalTags, result.TotalClusters)
	for _, c := range result.Clusters {
		tags := make([]string, len(c.Tags))
		for i, tag := range c.Tags {
			tags[i] = fmt.Sprintf("%s(%d)", tag.Name, tag.Count)
		}
		p("[%s] %s", c.Stem, strings.Join(tags, " ↔ "))
	}
}

func TestShowcaseRenameTag(t *testing.T) {
	dir, files := setupShowcaseDir(t)

	fmt.Println("\n=== rename-tag dry-run ===")
	result, _ := RenameTag(files, "emacs", "editor", true)
	fmt.Printf("  --from %q --to %q → %d files affected\n", result.OldTag, result.NewTag, result.Modified)
	for _, f := range result.Files {
		p("→ %s", filepath.Base(f))
	}

	// Actual rename on a copy
	fmt.Println("\n=== rename-tag actual (config → cfg) ===")
	files2 := ScanDirs([]string{dir})
	result2, _ := RenameTag(files2, "config", "cfg", false)
	fmt.Printf("  --from %q --to %q → %d files modified\n", result2.OldTag, result2.NewTag, result2.Modified)
	// Re-scan to show result
	files3 := ScanDirs([]string{dir})
	for _, f := range files3 {
		if hasTag(f.Tags, "cfg") {
			p("→ %s [%s]", filepath.Base(f.Path), strings.Join(f.Tags, ","))
		}
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
