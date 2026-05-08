package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseDailyJournal(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "20230222T000000--2023-02-22__archive_journal.org")
	os.WriteFile(path, []byte(`* 2023-Week-08 23-02-22 (Wednesday)
  Everything is hard before it is easy.
** 06:20 기상 약 복용
** 07:02 준비하고 등원
** 12:23 시작하자
** 22:55 자고 일어나서 하자. 늦었다.
`), 0644)

	entries := parseDailyJournal(path)
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
	if entries[0].Time != "06:20" || entries[0].Text != "기상 약 복용" {
		t.Errorf("entry[0] = %+v", entries[0])
	}
	if entries[3].Time != "22:55" {
		t.Errorf("entry[3].Time = %q", entries[3].Time)
	}
}

func TestParseWeeklyJournal(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "20260216T000000--2026-02-16__journal_week07.org")
	os.WriteFile(path, []byte(`#+title: 2026-02-16
#+filetags: :journal:week07:

* 2026-02-16 Monday
** 02:32 자다깨다
** 09:06 좋은 아침이로세

* 2026-02-17 Tuesday
** 08:00 출근
** 12:00 점심

* 2026-02-18 Wednesday
** 09:30 회의
`), 0644)

	// Monday
	entries := parseWeeklyJournal(path, "2026-02-16", mustParseDate("2026-02-16"))
	if len(entries) != 2 {
		t.Fatalf("Monday: expected 2 entries, got %d", len(entries))
	}
	if entries[0].Time != "02:32" {
		t.Errorf("entries[0].Time = %q", entries[0].Time)
	}

	// Tuesday
	entries = parseWeeklyJournal(path, "2026-02-17", mustParseDate("2026-02-17"))
	if len(entries) != 2 {
		t.Fatalf("Tuesday: expected 2 entries, got %d", len(entries))
	}
	if entries[0].Text != "출근" {
		t.Errorf("entries[0].Text = %q", entries[0].Text)
	}

	// Wednesday
	entries = parseWeeklyJournal(path, "2026-02-18", mustParseDate("2026-02-18"))
	if len(entries) != 1 {
		t.Fatalf("Wednesday: expected 1 entry, got %d", len(entries))
	}

	// Thursday (no entries)
	entries = parseWeeklyJournal(path, "2026-02-19", mustParseDate("2026-02-19"))
	if len(entries) != 0 {
		t.Errorf("Thursday: expected 0 entries, got %d", len(entries))
	}
}

func TestParseDatetree(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "20220101T010100--diary.org")
	os.WriteFile(path, []byte(`** 2023-02 February
*** 2023-02-21 Tuesday
**** 16:59 - 작업
:LOGBOOK:
CLOCK: [2023-02-21 Tue 16:59]--[2023-02-21 Tue 21:36] =>  4:37
:END:
*** 2023-02-22 Wednesday
**** 07:45 - 등원 준비
:LOGBOOK:
CLOCK: [2023-02-22 Wed 07:45]--[2023-02-22 Wed 07:46] =>  0:01
:END:
**** 08:02 - 등원 가야지
:LOGBOOK:
CLOCK: [2023-02-22 Wed 08:02]--[2023-02-22 Wed 12:22] =>  4:20
:END:
**** 17:10 - 하원 간다. 늦었다.
:LOGBOOK:
CLOCK: [2023-02-22 Wed 17:10]--[2023-02-22 Wed 21:03] =>  3:53
:END:
*** 2023-02-23 Thursday
**** 06:14 - 대변
`), 0644)

	entries := parseDatetree(path, "2023-02-22 Wednesday")
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Time != "07:45" || entries[0].Text != "등원 준비" {
		t.Errorf("entries[0] = %+v", entries[0])
	}
	if entries[0].Clock != "0:01" {
		t.Errorf("entries[0].Clock = %q, want 0:01", entries[0].Clock)
	}
	if entries[1].Clock != "4:20" {
		t.Errorf("entries[1].Clock = %q, want 4:20", entries[1].Clock)
	}
	if entries[2].Clock != "3:53" {
		t.Errorf("entries[2].Clock = %q, want 3:53", entries[2].Clock)
	}
}

func TestParseDatetreeLevel4(t *testing.T) {
	// diary.org uses **** for dates, ***** for entries
	tmp := t.TempDir()
	path := filepath.Join(tmp, "diary.org")
	os.WriteFile(path, []byte(`*** 2023-02 February
**** 2023-02-22 Wednesday
***** 07:45 - 등원 준비
:LOGBOOK:
CLOCK: [2023-02-22 Wed 07:45]--[2023-02-22 Wed 07:46] =>  0:01
:END:
***** 08:02 - 등원 가야지
:LOGBOOK:
CLOCK: [2023-02-22 Wed 08:02]--[2023-02-22 Wed 12:22] =>  4:20
:END:
**** 2023-02-23 Thursday
***** 06:14 - 대변
`), 0644)

	entries := parseDatetree(path, "2023-02-22 Wednesday")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestQueryDayIntegration(t *testing.T) {
	// Create full test structure
	base := t.TempDir()

	// journal/ directory
	journalDir := filepath.Join(base, "journal")
	os.MkdirAll(journalDir, 0755)

	// Daily archive file
	os.WriteFile(filepath.Join(journalDir, "20230222T000000--2023-02-22__archive_journal.org"),
		[]byte(`* 2023-02-22
** 06:20 기상
** 12:00 점심
`), 0644)

	// diary.org
	os.WriteFile(filepath.Join(base, "20220101T010100--diary.org"),
		[]byte(`*** 2023-02-22 Wednesday
**** 07:45 - 등원
:LOGBOOK:
CLOCK: [2023-02-22 Wed 07:45]--[2023-02-22 Wed 07:46] =>  0:01
:END:
*** 2023-02-23 Thursday
**** 06:00 - next day
`), 0644)

	// A note created on this date
	notesDir := filepath.Join(base, "notes")
	os.MkdirAll(notesDir, 0755)
	os.WriteFile(filepath.Join(notesDir, "20230222T123300--테스트-노트__test.org"),
		[]byte(`#+title: 테스트 노트
`), 0644)

	result := QueryDay([]string{base}, "2023-02-22")

	if result.Date != "2023-02-22" {
		t.Errorf("date = %q", result.Date)
	}
	if result.DayOfWeek != "Wednesday" {
		t.Errorf("day = %q", result.DayOfWeek)
	}
	if result.Journal == nil {
		t.Fatal("journal is nil")
	}
	if len(result.Journal.Entries) != 2 {
		t.Errorf("journal entries = %d, want 2", len(result.Journal.Entries))
	}
	if result.Journal.Format != "daily-archive" {
		t.Errorf("journal format = %q", result.Journal.Format)
	}
	if result.Datetree == nil {
		t.Fatal("datetree is nil")
	}
	if len(result.Datetree.Entries) != 1 {
		t.Errorf("datetree entries = %d, want 1", len(result.Datetree.Entries))
	}
	if result.Datetree.Entries[0].Clock != "0:01" {
		t.Errorf("datetree clock = %q", result.Datetree.Entries[0].Clock)
	}
	if len(result.NotesCreated) != 1 {
		t.Errorf("notes = %d, want 1", len(result.NotesCreated))
	}
	if len(result.NotesModified) != 0 {
		t.Errorf("modified notes = %d, want 0", len(result.NotesModified))
	}
}

func TestQueryDayIncludesNotesModified(t *testing.T) {
	base := t.TempDir()
	notesDir := filepath.Join(base, "notes")
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(notesDir, "20230222T123300--created-today__test.org"), []byte(`#+title: Created Today
#+hugo_lastmod: [2023-02-22 Wed 18:00]
`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(notesDir, "20230220T090000--modified-today__test.org"), []byte(`#+title: Modified Today
#+hugo_lastmod: [2023-02-22 Wed 17:00]
`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(notesDir, "20230219T090000--untouched__test.org"), []byte(`#+title: Untouched
#+hugo_lastmod: [2023-02-21 Tue 17:00]
`), 0644); err != nil {
		t.Fatal(err)
	}

	result := QueryDay([]string{base}, "2023-02-22")
	if len(result.NotesCreated) != 1 {
		t.Fatalf("notes created = %d, want 1", len(result.NotesCreated))
	}
	if result.NotesCreated[0].ID != "20230222T123300" {
		t.Fatalf("created note ID = %q", result.NotesCreated[0].ID)
	}
	if len(result.NotesModified) != 1 {
		t.Fatalf("notes modified = %d, want 1", len(result.NotesModified))
	}
	if result.NotesModified[0].ID != "20230220T090000" {
		t.Fatalf("modified note ID = %q", result.NotesModified[0].ID)
	}
	if result.NotesModified[0].Lastmod != "2023-02-22" {
		t.Fatalf("modified note lastmod = %q", result.NotesModified[0].Lastmod)
	}
}

func TestCleanHeadingText(t *testing.T) {
	tests := []struct{ in, want string }{
		{"TODO task", "task"},
		{"DONE finished", "finished"},
		{"text :tag1:tag2:", "text"},
		{"plain text", "plain text"},
	}
	for _, tt := range tests {
		got := cleanHeadingText(tt.in)
		if got != tt.want {
			t.Errorf("cleanHeadingText(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFindJournalWeeklyFromMidweek(t *testing.T) {
	base := t.TempDir()
	journalDir := filepath.Join(base, "journal")
	os.MkdirAll(journalDir, 0755)

	// Week file starts on Monday 2026-02-16
	os.WriteFile(filepath.Join(journalDir, "20260216T000000--2026-02-16__journal_week07.org"),
		[]byte(`#+title: 2026-02-16
* 2026-02-16 Monday
** 09:00 monday work
* 2026-02-18 Wednesday
** 10:00 wednesday meeting
* 2026-02-20 Friday
** 14:00 friday review
`), 0644)

	// Query Wednesday (mid-week)
	result := QueryDay([]string{base}, "2026-02-18")
	if result.Journal == nil {
		t.Fatal("journal is nil for Wednesday")
	}
	if result.Journal.Format != "weekly" {
		t.Errorf("format = %q, want weekly", result.Journal.Format)
	}
	if len(result.Journal.Entries) != 1 {
		t.Errorf("entries = %d, want 1", len(result.Journal.Entries))
	}
	if result.Journal.Entries[0].Text != "wednesday meeting" {
		t.Errorf("text = %q", result.Journal.Entries[0].Text)
	}
}

// TestQueryDayNotesModified exercises the hugo_lastmod axis end-to-end.
//
// Layout:
//   - notes/A: created on the queried date (goes to NotesCreated only)
//   - notes/B: older file, hugo_lastmod matches today (goes to NotesModified)
//   - notes/C: older file, no hugo_lastmod (skipped from both)
//   - notes/D: older file, hugo_lastmod = different date (skipped)
//   - notes/E: created on the queried date AND hugo_lastmod = today
//     → must appear in NotesCreated only, NotesModified must NOT duplicate
//   - notes/F: hugo_lastmod inside body (after first heading) → must be ignored
func TestQueryDayNotesModified(t *testing.T) {
	base := t.TempDir()
	notesDir := filepath.Join(base, "notes")
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		t.Fatal(err)
	}
	target := "2026-05-08"

	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(notesDir, name), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("20260508T093000--note-a__test.org",
		"#+title: A\n#+identifier: 20260508T093000\n")
	write("20240906T154822--note-b__test.org",
		"#+title: B\n#+identifier: 20240906T154822\n#+hugo_lastmod: [2026-05-08 Fri 22:24]\n")
	write("20240906T154823--note-c__test.org",
		"#+title: C\n#+identifier: 20240906T154823\n")
	write("20240906T154824--note-d__test.org",
		"#+title: D\n#+identifier: 20240906T154824\n#+hugo_lastmod: [2025-12-25]\n")
	write("20260508T100000--note-e__test.org",
		"#+title: E\n#+identifier: 20260508T100000\n#+hugo_lastmod: [2026-05-08]\n")
	write("20240906T154825--note-f__test.org",
		"#+title: F\n#+identifier: 20240906T154825\n\n* body\n#+hugo_lastmod: [2026-05-08]\n")

	result := QueryDay([]string{base}, target)

	createdIDs := map[string]bool{}
	for _, f := range result.NotesCreated {
		createdIDs[f.ID] = true
	}
	if !createdIDs["20260508T093000"] || !createdIDs["20260508T100000"] {
		t.Errorf("NotesCreated missing expected IDs: got %+v", createdIDs)
	}

	modifiedIDs := map[string]string{}
	for _, f := range result.NotesModified {
		modifiedIDs[f.ID] = f.Lastmod
	}
	if got, ok := modifiedIDs["20240906T154822"]; !ok || got != target {
		t.Errorf("note B should be in NotesModified with lastmod=%s, got %q (ok=%v)", target, got, ok)
	}
	if _, dup := modifiedIDs["20260508T100000"]; dup {
		t.Errorf("note E was created today; must NOT duplicate into NotesModified")
	}
	if _, leak := modifiedIDs["20240906T154823"]; leak {
		t.Errorf("note C has no hugo_lastmod; must NOT appear in NotesModified")
	}
	if _, leak := modifiedIDs["20240906T154824"]; leak {
		t.Errorf("note D has different hugo_lastmod; must NOT appear in NotesModified")
	}
	if _, leak := modifiedIDs["20240906T154825"]; leak {
		t.Errorf("note F's hugo_lastmod sits in body; must NOT appear in NotesModified")
	}
}

func TestQueryDayNotesModifiedSkipsJournal(t *testing.T) {
	// A journal file that happens to carry hugo_lastmod must not pollute the
	// modified axis — journal entries already have their own surface.
	base := t.TempDir()
	journalDir := filepath.Join(base, "journal")
	if err := os.MkdirAll(journalDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(journalDir, "20240101T000000--2024-01-01__journal.org"),
		[]byte("#+title: 2024-01-01\n#+identifier: 20240101T000000\n#+hugo_lastmod: [2026-05-08]\n** 09:00 something\n"),
		0644,
	); err != nil {
		t.Fatal(err)
	}
	result := QueryDay([]string{base}, "2026-05-08")
	if len(result.NotesModified) != 0 {
		t.Errorf("journal hugo_lastmod must be ignored, got %+v", result.NotesModified)
	}
}

func mustParseDate(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}
