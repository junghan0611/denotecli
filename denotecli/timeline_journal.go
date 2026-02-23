// timeline_journal.go — denotecli timeline-journal: monthly journal overview
package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TimelineJournalResult is the output of timeline-journal.
type TimelineJournalResult struct {
	Period     string              `json:"period"`
	TotalDays  int                 `json:"total_days"`
	ActiveDays int                 `json:"active_days"`
	Days       []TimelineJournalDay `json:"days"`
}

// TimelineJournalDay summarizes a single day.
type TimelineJournalDay struct {
	Date         string   `json:"date"`
	DayOfWeek    string   `json:"day_of_week"`
	Sources      []string `json:"sources"`       // "journal", "datetree"
	JournalCount int      `json:"journal_count"`  // entries from journal
	DatetreeCount int     `json:"datetree_count"` // entries from datetree
	NotesCount   int      `json:"notes_count"`    // denote notes created
}

var datetreeDayRe = regexp.MustCompile(`^\*{3,5}\s+(\d{4}-\d{2}-\d{2})\s+\w+`)

// QueryTimelineJournal returns a monthly overview of journal activity.
func QueryTimelineJournal(dirs []string, month string, fromStr, toStr string) TimelineJournalResult {
	var from, to time.Time

	if month != "" {
		t, err := time.ParseInLocation("2006-01", month, time.Local)
		if err != nil {
			fatal("invalid --month: " + month + " (use YYYY-MM)")
		}
		from = t
		to = t.AddDate(0, 1, -1)
	} else if fromStr != "" && toStr != "" {
		var err error
		from, err = time.ParseInLocation("2006-01-02", fromStr, time.Local)
		if err != nil {
			fatal("invalid --from: " + fromStr)
		}
		to, err = time.ParseInLocation("2006-01-02", toStr, time.Local)
		if err != nil {
			fatal("invalid --to: " + toStr)
		}
	} else {
		// Default: current month
		now := time.Now()
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		to = from.AddDate(0, 1, -1)
	}

	result := TimelineJournalResult{
		Period: from.Format("2006-01-02") + " ~ " + to.Format("2006-01-02"),
	}

	// Count total days in range
	totalDays := int(to.Sub(from).Hours()/24) + 1
	result.TotalDays = totalDays

	// Build day map
	dayMap := make(map[string]*TimelineJournalDay)

	for _, dir := range dirs {
		dir = expandHome(dir)

		// Scan journal/ files
		scanJournalTimeline(dir, from, to, dayMap)

		// Scan diary.org datetree
		scanDatetreeTimeline(dir, from, to, dayMap)

		// Scan notes created
		scanNotesTimeline(dir, from, to, dayMap)
	}

	// Sort and collect
	dates := make([]string, 0, len(dayMap))
	for d := range dayMap {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	for _, d := range dates {
		result.Days = append(result.Days, *dayMap[d])
	}
	result.ActiveDays = len(result.Days)

	return result
}

func getOrCreateDay(dayMap map[string]*TimelineJournalDay, dateStr string) *TimelineJournalDay {
	if d, ok := dayMap[dateStr]; ok {
		return d
	}
	t, _ := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	d := &TimelineJournalDay{
		Date:      dateStr,
		DayOfWeek: t.Weekday().String(),
	}
	dayMap[dateStr] = d
	return d
}

func addSource(d *TimelineJournalDay, source string) {
	for _, s := range d.Sources {
		if s == source {
			return
		}
	}
	d.Sources = append(d.Sources, source)
}

func scanJournalTimeline(dir string, from, to time.Time, dayMap map[string]*TimelineJournalDay) {
	journalDir := filepath.Join(dir, "journal")
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".org") {
			continue
		}
		if len(name) < 8 {
			continue
		}

		// Extract date from filename
		dateCompact := name[:8]
		if len(dateCompact) != 8 {
			continue
		}
		dateStr := dateCompact[:4] + "-" + dateCompact[4:6] + "-" + dateCompact[6:8]
		fileDate, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil {
			continue
		}

		// Weekly files: scan all days within
		if strings.Contains(name, "week") {
			path := filepath.Join(journalDir, name)
			scanWeeklyForTimeline(path, from, to, dayMap)
			continue
		}

		// Daily files: check date range
		if !strings.Contains(name, "journal") {
			continue
		}
		if fileDate.Before(from) || fileDate.After(to) {
			continue
		}

		// Count entries
		path := filepath.Join(journalDir, name)
		count := countTimeEntries(path)
		if count > 0 {
			d := getOrCreateDay(dayMap, dateStr)
			d.JournalCount = count
			addSource(d, "journal")
		}
	}
}

func scanWeeklyForTimeline(path string, from, to time.Time, dayMap map[string]*TimelineJournalDay) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	var currentDate string
	var count int

	flushDay := func() {
		if currentDate != "" && count > 0 {
			t, err := time.ParseInLocation("2006-01-02", currentDate, time.Local)
			if err == nil && !t.Before(from) && !t.After(to) {
				d := getOrCreateDay(dayMap, currentDate)
				d.JournalCount = count
				addSource(d, "journal")
			}
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "* ") && len(line) >= 13 {
			rest := line[2:]
			if len(rest) >= 10 && rest[4] == '-' && rest[7] == '-' {
				flushDay()
				currentDate = rest[:10]
				count = 0
			}
		} else if timeEntryRe.MatchString(line) {
			count++
		}
	}
	flushDay()
}

func countTimeEntries(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if timeEntryRe.MatchString(scanner.Text()) {
			count++
		}
	}
	return count
}

func scanDatetreeTimeline(dir string, from, to time.Time, dayMap map[string]*TimelineJournalDay) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		name := e.Name()
		if !strings.Contains(name, "diary") || !strings.HasSuffix(name, ".org") {
			continue
		}

		path := filepath.Join(dir, name)
		f, err := os.Open(path)
		if err != nil {
			continue
		}

		var currentDate string
		var count int

		flushDay := func() {
			if currentDate != "" && count > 0 {
				t, err := time.ParseInLocation("2006-01-02", currentDate, time.Local)
				if err == nil && !t.Before(from) && !t.After(to) {
					d := getOrCreateDay(dayMap, currentDate)
					d.DatetreeCount = count
					addSource(d, "datetree")
				}
			}
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			m := datetreeDayRe.FindStringSubmatch(line)
			if m != nil {
				flushDay()
				currentDate = m[1]
				count = 0
			} else if currentDate != "" && strings.HasPrefix(strings.TrimSpace(line), "*") {
				trimmed := strings.TrimSpace(line)
				stars := 0
				for _, c := range trimmed {
					if c == '*' {
						stars++
					} else {
						break
					}
				}
				rest := strings.TrimSpace(trimmed[stars:])
				if len(rest) >= 5 && rest[2] == ':' {
					// Looks like HH:MM
					h, _ := strconv.Atoi(rest[:2])
					m2, _ := strconv.Atoi(rest[3:5])
					if h >= 0 && h <= 23 && m2 >= 0 && m2 <= 59 {
						count++
					}
				}
			}
		}
		flushDay()
		f.Close()
	}
}

func scanNotesTimeline(dir string, from, to time.Time, dayMap map[string]*TimelineJournalDay) {
	files := ScanDirs([]string{dir})
	for _, f := range files {
		if len(f.ID) < 8 {
			continue
		}
		dateStr := f.ID[:4] + "-" + f.ID[4:6] + "-" + f.ID[6:8]
		t, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil || t.Before(from) || t.After(to) {
			continue
		}

		// Skip journal files
		base := filepath.Base(f.Path)
		if strings.Contains(base, "__journal") || strings.Contains(base, "__archive_journal") {
			continue
		}

		d := getOrCreateDay(dayMap, dateStr)
		d.NotesCount++
	}
}
