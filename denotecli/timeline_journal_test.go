package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestQueryTimelineJournalDaily(t *testing.T) {
	base := t.TempDir()
	journalDir := filepath.Join(base, "journal")
	os.MkdirAll(journalDir, 0755)

	// Create 3 daily files
	os.WriteFile(filepath.Join(journalDir, "20230201T000000--2023-02-01__archive_journal.org"),
		[]byte("* 2023-02-01\n** 06:00 기상\n** 12:00 점심\n"), 0644)
	os.WriteFile(filepath.Join(journalDir, "20230202T000000--2023-02-02__archive_journal.org"),
		[]byte("* 2023-02-02\n** 07:00 기상\n"), 0644)
	os.WriteFile(filepath.Join(journalDir, "20230203T000000--2023-02-03__journal.org"),
		[]byte("* 2023-02-03\n** 09:00 출근\n** 18:00 퇴근\n** 22:00 독서\n"), 0644)

	result := QueryTimelineJournal([]string{base}, "2023-02", "", "")

	if result.ActiveDays != 3 {
		t.Errorf("active_days = %d, want 3", result.ActiveDays)
	}
	if result.TotalDays != 28 {
		t.Errorf("total_days = %d, want 28", result.TotalDays)
	}
	if len(result.Days) != 3 {
		t.Fatalf("days = %d, want 3", len(result.Days))
	}
	if result.Days[0].JournalCount != 2 {
		t.Errorf("day[0] journal_count = %d, want 2", result.Days[0].JournalCount)
	}
	if result.Days[2].JournalCount != 3 {
		t.Errorf("day[2] journal_count = %d, want 3", result.Days[2].JournalCount)
	}
}

func TestQueryTimelineJournalWeekly(t *testing.T) {
	base := t.TempDir()
	journalDir := filepath.Join(base, "journal")
	os.MkdirAll(journalDir, 0755)

	os.WriteFile(filepath.Join(journalDir, "20260216T000000--2026-02-16__journal_week07.org"),
		[]byte("* 2026-02-16 Monday\n** 09:00 work\n** 12:00 lunch\n* 2026-02-17 Tuesday\n** 10:00 meeting\n"), 0644)

	result := QueryTimelineJournal([]string{base}, "2026-02", "", "")

	// Should find Monday (2 entries) and Tuesday (1 entry)
	found := map[string]int{}
	for _, d := range result.Days {
		found[d.Date] = d.JournalCount
	}
	if found["2026-02-16"] != 2 {
		t.Errorf("Monday = %d, want 2", found["2026-02-16"])
	}
	if found["2026-02-17"] != 1 {
		t.Errorf("Tuesday = %d, want 1", found["2026-02-17"])
	}
}

func TestQueryTimelineJournalDatetree(t *testing.T) {
	base := t.TempDir()

	os.WriteFile(filepath.Join(base, "20220101T010100--diary.org"),
		[]byte(`** 2023-02 February
*** 2023-02-22 Wednesday
**** 07:45 - 등원
**** 08:02 - 등원 가야지
**** 17:10 - 하원
*** 2023-02-23 Thursday
**** 06:14 - 기상
`), 0644)

	result := QueryTimelineJournal([]string{base}, "2023-02", "", "")

	found := map[string]int{}
	for _, d := range result.Days {
		found[d.Date] = d.DatetreeCount
	}
	if found["2023-02-22"] != 3 {
		t.Errorf("02-22 datetree = %d, want 3", found["2023-02-22"])
	}
	if found["2023-02-23"] != 1 {
		t.Errorf("02-23 datetree = %d, want 1", found["2023-02-23"])
	}
}

func TestQueryTimelineJournalNotes(t *testing.T) {
	base := t.TempDir()
	notesDir := filepath.Join(base, "notes")
	os.MkdirAll(notesDir, 0755)

	os.WriteFile(filepath.Join(notesDir, "20230222T123300--테스트__test.org"),
		[]byte("#+title: 테스트\n"), 0644)
	os.WriteFile(filepath.Join(notesDir, "20230222T150000--다른노트__note.org"),
		[]byte("#+title: 다른노트\n"), 0644)
	os.WriteFile(filepath.Join(notesDir, "20230225T100000--월말__note.org"),
		[]byte("#+title: 월말\n"), 0644)

	result := QueryTimelineJournal([]string{base}, "2023-02", "", "")

	found := map[string]int{}
	for _, d := range result.Days {
		found[d.Date] = d.NotesCount
	}
	if found["2023-02-22"] != 2 {
		t.Errorf("02-22 notes = %d, want 2", found["2023-02-22"])
	}
	if found["2023-02-25"] != 1 {
		t.Errorf("02-25 notes = %d, want 1", found["2023-02-25"])
	}
}

func TestQueryTimelineJournalFromTo(t *testing.T) {
	base := t.TempDir()
	journalDir := filepath.Join(base, "journal")
	os.MkdirAll(journalDir, 0755)

	os.WriteFile(filepath.Join(journalDir, "20230210T000000--2023-02-10__journal.org"),
		[]byte("* 2023-02-10\n** 09:00 work\n"), 0644)
	os.WriteFile(filepath.Join(journalDir, "20230220T000000--2023-02-20__journal.org"),
		[]byte("* 2023-02-20\n** 10:00 work\n"), 0644)

	// Only 10th~15th range
	result := QueryTimelineJournal([]string{base}, "", "2023-02-10", "2023-02-15")
	if result.ActiveDays != 1 {
		t.Errorf("active_days = %d, want 1 (only 02-10)", result.ActiveDays)
	}
}

func TestQueryTimelineJournalCombined(t *testing.T) {
	base := t.TempDir()
	journalDir := filepath.Join(base, "journal")
	os.MkdirAll(journalDir, 0755)

	// Journal + datetree on same day
	os.WriteFile(filepath.Join(journalDir, "20230222T000000--2023-02-22__archive_journal.org"),
		[]byte("* 2023-02-22\n** 06:00 기상\n** 12:00 점심\n"), 0644)
	os.WriteFile(filepath.Join(base, "20220101T010100--diary.org"),
		[]byte("*** 2023-02-22 Wednesday\n**** 07:45 - 등원\n*** 2023-02-23 Thursday\n**** 08:00 - x\n"), 0644)

	result := QueryTimelineJournal([]string{base}, "2023-02", "", "")

	for _, d := range result.Days {
		if d.Date == "2023-02-22" {
			if d.JournalCount != 2 {
				t.Errorf("journal = %d, want 2", d.JournalCount)
			}
			if d.DatetreeCount != 1 {
				t.Errorf("datetree = %d, want 1", d.DatetreeCount)
			}
			if len(d.Sources) != 2 {
				t.Errorf("sources = %v, want [journal, datetree]", d.Sources)
			}
			return
		}
	}
	t.Error("2023-02-22 not found")
}
