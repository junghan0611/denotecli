// day.go — denotecli day: journal timeline for a specific date
package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DayResult is the top-level output of `denotecli day`.
//
// notes_created and notes_modified are complementary, not overlapping —
// a file created on the queried date is excluded from notes_modified even if
// its #+hugo_lastmod also matches, to avoid duplication. notes_modified
// reports only files whose creation date differs from the queried date but
// whose explicit #+hugo_lastmod field resolves to it.
type DayResult struct {
	Date          string        `json:"date"`
	DayOfWeek     string        `json:"day_of_week"`
	YearsAgo      int           `json:"years_ago"`
	Journal       *JournalData  `json:"journal"`
	Datetree      *DatetreeData `json:"datetree"`
	NotesCreated  []DenoteFile  `json:"notes_created"`
	NotesModified []DenoteFile  `json:"notes_modified"`
}

// JournalData holds parsed journal entries for a day.
type JournalData struct {
	Source  string         `json:"source"`
	Format  string         `json:"format"` // "daily", "daily-archive", "weekly"
	Entries []JournalEntry `json:"entries"`
}

// JournalEntry is a single time entry in a journal.
type JournalEntry struct {
	Time string `json:"time"`
	Text string `json:"text"`
}

// DatetreeData holds parsed diary.org entries for a day.
type DatetreeData struct {
	Source  string          `json:"source"`
	Entries []DatetreeEntry `json:"entries"`
}

// DatetreeEntry is a single time entry in diary.org datetree.
type DatetreeEntry struct {
	Time  string `json:"time"`
	Text  string `json:"text"`
	Clock string `json:"clock,omitempty"`
}

var timeEntryRe = regexp.MustCompile(`^\*{2,}\s+(\d{2}:\d{2})\s+(.*)`)
var clockRe = regexp.MustCompile(`CLOCK:\s+\[.*?\]--\[.*?\]\s+=>\s+(\S+)`)

// QueryDay builds a complete day view from all sources.
func QueryDay(dirs []string, dateStr string) DayResult {
	t, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		fatal("invalid date: " + dateStr)
	}

	result := DayResult{
		Date:          dateStr,
		DayOfWeek:     t.Weekday().String(),
		NotesCreated:  []DenoteFile{},
		NotesModified: []DenoteFile{},
	}

	// Calculate years ago
	now := time.Now()
	if t.Before(now) {
		years := now.Year() - t.Year()
		if years > 0 {
			result.YearsAgo = years
		}
	}

	// 1. Journal
	result.Journal = findJournal(dirs, dateStr, t)

	// 2. Diary.org datetree
	result.Datetree = findDatetree(dirs, dateStr)

	// 3. Notes created on this date
	datePrefix := strings.ReplaceAll(dateStr, "-", "")
	files := ScanDirs(dirs)
	createdIDs := make(map[string]bool)
	for _, f := range files {
		if strings.HasPrefix(f.ID, datePrefix) {
			// Skip journal files themselves
			base := filepath.Base(f.Path)
			if strings.Contains(base, "__journal") || strings.Contains(base, "__archive_journal") {
				continue
			}
			result.NotesCreated = append(result.NotesCreated, f)
			createdIDs[f.ID] = true
		}
	}

	// 4. Notes modified on this date — explicit #+hugo_lastmod only.
	//    Skips files already in NotesCreated (axes are complementary).
	//    Reads only frontmatter head per file (~50 lines) to keep the day
	//    command interactive even on 3000+ note corpora.
	result.NotesModified = collectNotesModified(files, dateStr, createdIDs)

	return result
}

// collectNotesModified scans Denote files for #+hugo_lastmod matching dateStr.
// Files with IDs in skip are excluded (they are reported via NotesCreated).
// Journal files are also skipped. The returned DenoteFile entries have their
// Lastmod field populated with the resolved YYYY-MM-DD value.
func collectNotesModified(files []DenoteFile, dateStr string, skip map[string]bool) []DenoteFile {
	out := []DenoteFile{}
	for _, f := range files {
		if skip[f.ID] {
			continue
		}
		base := filepath.Base(f.Path)
		if strings.Contains(base, "__journal") || strings.Contains(base, "__archive_journal") {
			continue
		}
		fm := ReadFrontmatterFile(f.Path, 50)
		if fm.HugoLastmod == "" {
			continue
		}
		d := ExtractDate(fm.HugoLastmod)
		if d != dateStr {
			continue
		}
		f.Lastmod = d
		out = append(out, f)
	}
	return out
}

// findJournal searches journal/ directory for a specific date.
func findJournal(dirs []string, dateStr string, t time.Time) *JournalData {
	dateCompact := strings.ReplaceAll(dateStr, "-", "")

	for _, dir := range dirs {
		dir = expandHome(dir)
		journalDir := filepath.Join(dir, "journal")
		if _, err := os.Stat(journalDir); err != nil {
			continue
		}

		// Try 1: daily file (archive or regular)
		// Pattern: YYYYMMDDT000000--YYYY-MM-DD__*journal*.org
		entries, _ := os.ReadDir(journalDir)
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, dateCompact+"T000000") {
				continue
			}
			if !strings.HasSuffix(name, ".org") {
				continue
			}

			path := filepath.Join(journalDir, name)
			format := ""
			if strings.Contains(name, "week") {
				continue // skip weekly files in daily pass
			} else if strings.Contains(name, "__archive_journal") {
				format = "daily-archive"
			} else if strings.Contains(name, "__journal") {
				format = "daily"
			} else {
				continue // not a journal file
			}

			entries := parseDailyJournal(path)
			if len(entries) > 0 {
				return &JournalData{
					Source:  path,
					Format:  format,
					Entries: entries,
				}
			}
		}

		// Try 2: weekly file
		// Find the Monday of the week containing this date
		weekday := t.Weekday()
		daysFromMonday := int(weekday) - int(time.Monday)
		if daysFromMonday < 0 {
			daysFromMonday += 7
		}
		monday := t.AddDate(0, 0, -daysFromMonday)
		mondayCompact := monday.Format("20060102")

		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, mondayCompact+"T000000") {
				continue
			}
			if !strings.Contains(name, "week") {
				continue
			}
			if !strings.HasSuffix(name, ".org") {
				continue
			}

			path := filepath.Join(journalDir, name)
			dayEntries := parseWeeklyJournal(path, dateStr, t)
			if len(dayEntries) > 0 {
				return &JournalData{
					Source:  path,
					Format:  "weekly",
					Entries: dayEntries,
				}
			}
		}
	}
	return nil
}

// parseDailyJournal extracts time entries from a daily journal file.
// Looks for level 2 headings: ** HH:MM text
func parseDailyJournal(path string) []JournalEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var entries []JournalEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		m := timeEntryRe.FindStringSubmatch(line)
		if m != nil {
			// Count asterisks to determine level
			text := strings.TrimSpace(m[2])
			entries = append(entries, JournalEntry{
				Time: m[1],
				Text: cleanHeadingText(text),
			})
		}
	}
	return entries
}

// parseWeeklyJournal extracts entries for a specific date from a weekly file.
// Weekly format: * YYYY-MM-DD DayName → ** HH:MM text
func parseWeeklyJournal(path, targetDate string, t time.Time) []JournalEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	dayHeader := "* " + targetDate
	inDay := false
	var entries []JournalEntry

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for day header: * YYYY-MM-DD DayName
		if strings.HasPrefix(line, "* ") {
			if strings.HasPrefix(line, dayHeader) {
				inDay = true
				continue
			} else if inDay {
				// Next day section, stop
				break
			}
			continue
		}

		if !inDay {
			continue
		}

		// Inside our day, look for ** HH:MM entries
		m := timeEntryRe.FindStringSubmatch(line)
		if m != nil {
			entries = append(entries, JournalEntry{
				Time: m[1],
				Text: cleanHeadingText(strings.TrimSpace(m[2])),
			})
		}
	}
	return entries
}

// findDatetree searches for diary.org datetree entries for a date.
func findDatetree(dirs []string, dateStr string) *DatetreeData {
	t, _ := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	dayName := t.Weekday().String()
	// Pattern to match: **** YYYY-MM-DD DayName
	dateHeader := dateStr + " " + dayName

	for _, dir := range dirs {
		dir = expandHome(dir)
		// Look for diary.org files (Denote-named)
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			name := e.Name()
			if !strings.Contains(name, "diary") || !strings.HasSuffix(name, ".org") {
				continue
			}

			path := filepath.Join(dir, name)
			dtEntries := parseDatetree(path, dateHeader)
			if len(dtEntries) > 0 {
				return &DatetreeData{
					Source:  path,
					Entries: dtEntries,
				}
			}
		}
	}
	return nil
}

// parseDatetree extracts entries from a datetree file for a specific date.
// Structure: *** YYYY-MM-DD DayName → **** HH:MM - text + CLOCK blocks
func parseDatetree(path, dateHeader string) []DatetreeEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	inDay := false
	var entries []DatetreeEntry
	var currentClock string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Datetree level varies: could be **** or *** depending on structure
		// Look for the date header at any heading level
		if strings.Contains(line, dateHeader) && strings.HasPrefix(trimmed, "*") {
			inDay = true
			continue
		}

		if !inDay {
			continue
		}

		// Any heading line
		if strings.HasPrefix(trimmed, "*") {
			// Count stars
			stars := 0
			for _, c := range trimmed {
				if c == '*' {
					stars++
				} else {
					break
				}
			}
			rest := strings.TrimSpace(trimmed[stars:])

			// If this heading looks like a date (next day), stop
			if len(rest) >= 10 && rest[4] == '-' && rest[7] == '-' && !strings.Contains(rest, dateHeader) {
				break
			}

			// Time entry: HH:MM - text (any level deeper than date heading)
			parts := strings.SplitN(rest, " ", 2)
			if len(parts) >= 1 && len(parts[0]) == 5 && parts[0][2] == ':' {
				text := ""
				if len(parts) >= 2 {
					text = strings.TrimPrefix(parts[1], "- ")
					text = strings.TrimPrefix(text, "-")
					text = strings.TrimSpace(text)
				}
				// Save previous entry's clock if any
				if len(entries) > 0 && currentClock != "" {
					entries[len(entries)-1].Clock = currentClock
				}
				currentClock = ""
				entries = append(entries, DatetreeEntry{
					Time: parts[0],
					Text: cleanHeadingText(text),
				})
			}
		}

		// CLOCK line
		m := clockRe.FindStringSubmatch(line)
		if m != nil {
			currentClock = m[1]
		}
	}

	// Last entry's clock
	if len(entries) > 0 && currentClock != "" {
		entries[len(entries)-1].Clock = currentClock
	}

	return entries
}

// cleanHeadingText removes org-mode markup from heading text.
func cleanHeadingText(s string) string {
	// Remove TODO/DONE keywords
	for _, kw := range []string{"TODO ", "DONE "} {
		s = strings.TrimPrefix(s, kw)
	}
	// Remove trailing tags :tag1:tag2:
	if idx := strings.LastIndex(s, " :"); idx > 0 {
		rest := s[idx+1:]
		if strings.HasSuffix(rest, ":") && !strings.Contains(rest, " ") {
			s = s[:idx]
		}
	}
	return strings.TrimSpace(s)
}
