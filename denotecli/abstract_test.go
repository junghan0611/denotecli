// abstract_test.go
//
// Regression tests for the two output defects GLG named on 2026-09-04:
//
//   1. `date` carried creation time only — `#+hugo_lastmod:` never reached
//      the JSON, so callers (sorge's sweep.py) re-parsed org files with their
//      own regex and dropped HH:MM, inventing debt that did not exist.
//   2. The leading `[!abstract] 이 노트에 대하여` block and `#+description:`
//      never reached the JSON, so every caller opened the file to learn what
//      a note actually is.

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const abstractNote = `#+title:      §denotecli: #담당자 day-query 설계 검토 통합 타임라인 스펙
#+date:       [2026-02-22 Sun 09:00]
#+filetags:   :botlog:denotecli:
#+hugo_lastmod: [2026-05-18 Mon 09:09]
#+identifier: 20260222T090000
#+description: 하루 단위로 엮어 보는 통합 타임라인 스펙.

#+begin_quote
[!abstract] 이 노트에 대하여

첫 문단이다.
둘째 줄이다.
#+end_quote

* 히스토리
- [2026-05-18 Mon 09:09] @junghan — lifetract와 연결
`

func TestExtractAbstract(t *testing.T) {
	a := ExtractAbstract(abstractNote)
	if a == nil {
		t.Fatal("ExtractAbstract returned nil for a note that carries [!abstract]")
	}
	if a.Kind != "abstract" {
		t.Errorf("Kind = %q, want %q", a.Kind, "abstract")
	}
	if a.Title != "이 노트에 대하여" {
		t.Errorf("Title = %q", a.Title)
	}
	want := "첫 문단이다.\n둘째 줄이다."
	if a.Body != want {
		t.Errorf("Body = %q, want %q", a.Body, want)
	}
}

func TestExtractAbstractAbsent(t *testing.T) {
	if a := ExtractAbstract("#+title: 없음\n\n* 헤딩\n본문\n"); a != nil {
		t.Errorf("want nil, got %+v", a)
	}
}

// A quote block that is not a callout must not be mistaken for an abstract.
func TestExtractAbstractPlainQuoteIgnored(t *testing.T) {
	src := "#+title: t\n\n#+begin_quote\n그냥 인용이다.\n#+end_quote\n\n* 헤딩\n"
	if a := ExtractAbstract(src); a != nil {
		t.Errorf("plain quote should not be an abstract, got %+v", a)
	}
}

// A callout that only appears below the first heading is body content, not
// the note's own abstract.
func TestExtractAbstractStopsAtFirstHeading(t *testing.T) {
	src := "#+title: t\n\n* 헤딩\n#+begin_quote\n[!note] 안쪽\n본문\n#+end_quote\n"
	if a := ExtractAbstract(src); a != nil {
		t.Errorf("callout after first heading should be ignored, got %+v", a)
	}
}

// The whole point of keeping the raw stamp: HH:MM survives, so a caller can
// compare a commit at 18:32 against a stamp at 21:55 without re-reading org.
func TestApplyFrontmatterKeepsLastmodTime(t *testing.T) {
	var f DenoteFile
	applyFrontmatter(&f, ParseFrontmatter(abstractNote))

	if f.HugoLastmod != "[2026-05-18 Mon 09:09]" {
		t.Errorf("HugoLastmod = %q, want raw org timestamp", f.HugoLastmod)
	}
	if f.Lastmod != "2026-05-18" {
		t.Errorf("Lastmod = %q, want normalized date", f.Lastmod)
	}
	if f.Description == "" {
		t.Error("Description empty — #+description: did not reach the struct")
	}
	if f.Date != "[2026-02-22 Sun 09:00]" {
		t.Errorf("Date = %q — creation stamp must stay intact", f.Date)
	}
}

// A note with no #+hugo_lastmod: must not gain a bogus stamp.
func TestApplyFrontmatterNoLastmod(t *testing.T) {
	f := DenoteFile{Date: "2026-01-01"}
	applyFrontmatter(&f, ParseFrontmatter("#+title: t\n#+identifier: 20260101T000000\n"))
	if f.Lastmod != "" || f.HugoLastmod != "" {
		t.Errorf("lastmod should stay empty, got %q / %q", f.Lastmod, f.HugoLastmod)
	}
}

// --- Frontmatter region: shapes that actually occur in the corpus ---
//
// Measured 2026-09-04 over ~/sync/org (3,478 notes): the old "stop at the first
// blank line" rule dropped #+description: on 4 notes. Each shape below is one of
// those files reduced to its skeleton, plus the guard that keeps body keywords out.

func TestParseFrontmatterBlankLineGap(t *testing.T) {
	// 20230904T144600 shape — a blank line splits the keyword block.
	fm := ParseFrontmatter("#+title:      t\n#+hugo_lastmod: [2025-05-28]\n#+identifier: 20230904T144600\n#+export_file_name: 20230904T144600.md\n\n#+description: 블록 사이의 빈 줄은 끝이 아니다.\n#+reference:  web-750wordswrite\n\n#+begin_quote\n[!abstract] 이 노트에 대하여\n\n본문.\n#+end_quote\n")
	if fm.Description == "" {
		t.Error("Description dropped at the blank-line gap")
	}
	if fm.HugoLastmod != "[2025-05-28]" {
		t.Errorf("HugoLastmod = %q", fm.HugoLastmod)
	}
}

func TestParseFrontmatterPropertiesDrawer(t *testing.T) {
	// 20241216T152928 shape — a gptel :PROPERTIES: drawer sits above the keywords
	// and its rows are arbitrary prose.
	fm := ParseFrontmatter(":PROPERTIES:\n:GPTEL_MODEL: grok-beta\n:GPTEL_SYSTEM: You are a large language model living in Emacs.\n:END:\n#+title:      t\n#+hugo_lastmod: [2024-12-16 Mon 15:29]\n#+identifier: 20241216T152928\n\n\n#+description: 드로어 뒤의 키워드도 읽어야 한다.\n")
	if fm.Title != "t" {
		t.Errorf("Title = %q — drawer rows must not end the region", fm.Title)
	}
	if fm.Description == "" || fm.HugoLastmod == "" {
		t.Errorf("lost fields after drawer: desc=%q lastmod=%q", fm.Description, fm.HugoLastmod)
	}
}

func TestParseFrontmatterNamedDrawerAndFileLocal(t *testing.T) {
	// 20231030T093000 shape — a :DOC-CONFIG: drawer (not :PROPERTIES:), and
	// 20220101T010100 shape — a bare `-*- … -*-` file-local line as line 1.
	fm := ParseFrontmatter("# -*- coding: utf-8-unix; -*-\n:DOC-CONFIG:\n#+property: header-args :tangle no\n#+startup: fold\n:END:\n#+title:      @norang #어젠다\n#+hugo_lastmod: [2025-02-18]\n#+identifier: 20231030T093000\n")
	if fm.Title == "" || fm.HugoLastmod == "" {
		t.Errorf("named drawer ended the region early: %+v", fm)
	}

	fm2 := ParseFrontmatter("-*- mode: Org; coding: utf-8-unix; -*-\n:PROPERTIES:\n:ID:       20220101T010100\n:END:\n#+title: diary\n#+identifier: 20220101T010100\n\n#+description: 파일 로컬 줄이 첫 줄이어도 읽는다.\n")
	if fm2.Title != "diary" || fm2.Description == "" {
		t.Errorf("file-local first line ended the region early: %+v", fm2)
	}
}

// The blank-line rule used to be the body-leak guard. Its replacement must hold
// the same line: a keyword inside a src block is body, not frontmatter.
func TestParseFrontmatterIgnoresKeywordInsideBlock(t *testing.T) {
	fm := ParseFrontmatter("#+title:      t\n#+identifier: 20260508T094035\n\n#+begin_src text\n#+hugo_lastmod: [2099-12-31]\n#+description: 본문 예시일 뿐이다.\n#+end_src\n")
	if fm.HugoLastmod != "" || fm.Description != "" {
		t.Errorf("body block leaked: lastmod=%q desc=%q", fm.HugoLastmod, fm.Description)
	}
}

func TestParseFrontmatterStopsAtProse(t *testing.T) {
	fm := ParseFrontmatter("#+title:      t\n#+identifier: 20260101T000000\n\n그냥 본문 문장이다.\n#+description: 여기는 본문이다.\n")
	if fm.Description != "" {
		t.Errorf("prose did not end the region: %q", fm.Description)
	}
}

// 62 corpus notes open the block as #+BEGIN_QUOTE.
func TestExtractAbstractUppercaseBlock(t *testing.T) {
	a := ExtractAbstract("#+title: t\n\n#+BEGIN_QUOTE\n[!abstract] 이 노트에 대하여\n\n본문.\n#+END_QUOTE\n\n* 헤딩\n")
	if a == nil || a.Body != "본문." {
		t.Errorf("uppercase block not extracted: %+v", a)
	}
}

func TestExtractAbstractUnclosedBlock(t *testing.T) {
	if a := ExtractAbstract("#+title: t\n\n#+begin_quote\n[!abstract] 제목\n본문\n"); a != nil {
		t.Errorf("unclosed block must not yield an abstract, got %+v", a)
	}
}

// A plain quote first, the real callout second — keep looking, do not give up.
func TestExtractAbstractSkipsPlainQuoteThenFinds(t *testing.T) {
	a := ExtractAbstract("#+title: t\n\n#+begin_quote\n평범한 인용.\n#+end_quote\n\n#+begin_quote\n[!note] 두 번째\n\n본문.\n#+end_quote\n\n* 헤딩\n")
	if a == nil || a.Kind != "note" || a.Title != "두 번째" {
		t.Errorf("second block not found: %+v", a)
	}
}

func TestExtractAbstractCalloutWithoutTitle(t *testing.T) {
	a := ExtractAbstract("#+title: t\n\n#+begin_quote\n[!abstract]\n\n본문만 있다.\n#+end_quote\n")
	if a == nil || a.Title != "" || a.Body != "본문만 있다." {
		t.Errorf("titleless callout = %+v", a)
	}
}

// Corpus variant: #+hugo_lastmod: Time-stamp: <2024-12-10 04:56:08 junghan>
func TestApplyFrontmatterTimeStampVariant(t *testing.T) {
	var f DenoteFile
	applyFrontmatter(&f, ParseFrontmatter("#+title: t\n#+identifier: 20240603T125251\n#+hugo_lastmod: Time-stamp: <2024-12-10 04:56:08 junghan>\n"))
	if f.Lastmod != "2024-12-10" {
		t.Errorf("Lastmod = %q, want 2024-12-10", f.Lastmod)
	}
	if !strings.HasPrefix(f.HugoLastmod, "Time-stamp:") {
		t.Errorf("HugoLastmod must stay raw, got %q", f.HugoLastmod)
	}
}

// --- Wiring: the fields must survive the trip to each command's JSON ---

func writeAbstractCorpus(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	stamped := "#+title:      도장 있는 노트\n#+date:       [2026-02-22 Sun 09:00]\n#+filetags:   :botlog:denotecli:\n#+hugo_lastmod: [2026-05-18 Mon 09:09]\n#+identifier: 20260222T090000\n#+description: 한 줄 요약이다.\n\n#+begin_quote\n[!abstract] 이 노트에 대하여\n\n무엇인지 말하는 자리.\n#+end_quote\n\n* 히스토리\n- 한 줄\n"
	bare := "#+title:      맨 노트\n#+identifier: 20260223T090000\n\n* 본문\n내용\n"
	must := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	must("20260222T090000--도장-있는-노트__botlog_denotecli.org", stamped)
	must("20260223T090000--맨-노트.org", bare)
	return dir
}

func TestScanDirsCarriesLastmodAndDescription(t *testing.T) {
	files := ScanDirs([]string{writeAbstractCorpus(t)})
	byID := map[string]DenoteFile{}
	for _, f := range files {
		byID[f.ID] = f
	}
	got, ok := byID["20260222T090000"]
	if !ok {
		t.Fatalf("note missing from scan: %+v", files)
	}
	if got.Lastmod != "2026-05-18" {
		t.Errorf("Lastmod = %q", got.Lastmod)
	}
	if got.HugoLastmod != "[2026-05-18 Mon 09:09]" {
		t.Errorf("HugoLastmod = %q — the raw stamp must reach search/list JSON", got.HugoLastmod)
	}
	if got.Description != "한 줄 요약이다." {
		t.Errorf("Description = %q", got.Description)
	}
	// A note without the keys must not gain invented values.
	if bare := byID["20260223T090000"]; bare.Lastmod != "" || bare.HugoLastmod != "" || bare.Description != "" {
		t.Errorf("bare note gained fields: %+v", bare)
	}
}

func TestReadByIDCarriesAbstract(t *testing.T) {
	files := ScanDirs([]string{writeAbstractCorpus(t)})
	dc, err := ReadByID(files, "20260222T090000", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if dc.Abstract == nil || dc.Abstract.Kind != "abstract" || dc.Abstract.Body != "무엇인지 말하는 자리." {
		t.Errorf("Abstract = %+v", dc.Abstract)
	}
	if dc.HugoLastmod == "" || dc.Description == "" {
		t.Errorf("read lost stamp/description: %+v", dc.DenoteFile)
	}
	// read keeps the org-original #+date:, unlike the ID-derived date in search.
	if dc.Date != "[2026-02-22 Sun 09:00]" {
		t.Errorf("Date = %q", dc.Date)
	}
}

func TestReadOutlineCarriesAbstract(t *testing.T) {
	files := ScanDirs([]string{writeAbstractCorpus(t)})
	do, err := ReadOutlineByID(files, "20260222T090000")
	if err != nil {
		t.Fatal(err)
	}
	if do.Abstract == nil {
		t.Fatal("--outline dropped the abstract — the first call must say what the note is")
	}
	if do.Abstract.Title != "이 노트에 대하여" || do.HugoLastmod == "" || do.Description == "" {
		t.Errorf("outline metadata incomplete: abstract=%+v file=%+v", do.Abstract, do.DenoteFile)
	}

	bare, err := ReadOutlineByID(files, "20260223T090000")
	if err != nil {
		t.Fatal(err)
	}
	if bare.Abstract != nil {
		t.Errorf("note without a callout must report no abstract, got %+v", bare.Abstract)
	}
}

// The JSON contract: optional keys are absent, not null/empty, when the note
// does not carry them — callers test presence, not value.
func TestOptionalFieldsOmittedFromJSON(t *testing.T) {
	files := ScanDirs([]string{writeAbstractCorpus(t)})
	do, err := ReadOutlineByID(files, "20260223T090000")
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(do)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"lastmod", "hugo_lastmod", "description", "abstract", "signature"} {
		if strings.Contains(string(b), "\""+key+"\"") {
			t.Errorf("key %q present for a note that lacks it: %s", key, b)
		}
	}

	full, err := ReadOutlineByID(files, "20260222T090000")
	if err != nil {
		t.Fatal(err)
	}
	b2, err := json.Marshal(full)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"lastmod", "hugo_lastmod", "description", "abstract"} {
		if !strings.Contains(string(b2), "\""+key+"\"") {
			t.Errorf("key %q missing for a note that carries it: %s", key, b2)
		}
	}
}
