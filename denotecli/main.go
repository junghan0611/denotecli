// main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const Version = "0.8.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "day":
		cmdDayQuery()
	case "timeline-journal":
		cmdTimelineJournal()
	case "search":
		cmdSearch()
	case "search-headings":
		cmdSearchHeadings()
	case "search-content":
		cmdSearchContent()
	case "read":
		cmdRead()
	case "create":
		cmdCreate()
	case "keyword-map":
		cmdKeywordMap()
	case "graph":
		cmdGraph()
	case "rename-tag":
		cmdRenameTag()
	case "tags":
		cmdTags()
	case "-V", "--version", "version":
		fmt.Println("denotecli " + Version)
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func cmdDayQuery() {
	args := os.Args[2:]
	date := ""
	if len(os.Args) >= 3 && !strings.HasPrefix(os.Args[2], "--") {
		date = os.Args[2]
		args = os.Args[3:]
	}
	if err := validateFlags(args,
		[]string{"--dirs", "--years-ago", "--days-ago"}, nil); err != nil {
		fatal(err.Error())
	}
	dirsStr := getFlag(args, "--dirs", "~/org")
	yearsAgo := getFlag(args, "--years-ago", "")
	daysAgo := getFlag(args, "--days-ago", "")

	resolved, err := resolveDayDate(date, yearsAgo, daysAgo)
	if err != nil {
		fatal(err.Error())
	}

	dirs := strings.Split(dirsStr, ",")
	result := QueryDay(dirs, resolved)
	printJSON(result)
}

func cmdTimelineJournal() {
	args := os.Args[2:]
	if err := validateFlags(args,
		[]string{"--dirs", "--month", "--from", "--to"}, nil); err != nil {
		fatal(err.Error())
	}
	dirsStr := getFlag(args, "--dirs", "~/org")
	month := getFlag(args, "--month", "")
	from := getFlag(args, "--from", "")
	to := getFlag(args, "--to", "")

	dirs := strings.Split(dirsStr, ",")
	result := QueryTimelineJournal(dirs, month, from, to)
	printJSON(result)
}

func cmdSearch() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli search <query> [--tags TAG] [--dirs DIR,...] [--title-only] [--max N]")
	}
	query := os.Args[2]
	args := os.Args[3:]
	if err := validateFlags(args,
		[]string{"--tags", "--dirs", "--max"},
		[]string{"--title-only"}); err != nil {
		fatal(err.Error())
	}
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

func cmdGraph() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli graph <id> [--dirs DIR,...]")
	}
	id := os.Args[2]
	args := os.Args[3:]
	if err := validateFlags(args, []string{"--dirs"}, nil); err != nil {
		fatal(err.Error())
	}
	dirsStr := getFlag(args, "--dirs", "~/org")

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	result, err := BuildGraph(files, id, 1)
	if err != nil {
		fatal(err.Error())
	}
	printJSON(result)
}

func cmdKeywordMap() {
	args := os.Args[2:]
	query := ""
	if len(os.Args) >= 3 && !strings.HasPrefix(os.Args[2], "--") {
		query = os.Args[2]
		args = os.Args[3:]
	}
	if err := validateFlags(args,
		[]string{"--dirs", "--max"}, nil); err != nil {
		fatal(err.Error())
	}
	dirsStr := getFlag(args, "--dirs", "~/org")
	maxStr := getFlag(args, "--max", "50")
	max, _ := strconv.Atoi(maxStr)
	if max <= 0 {
		max = 50
	}

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	result := BuildKeywordMap(files, query)
	if len(result.Entries) > max {
		result.Entries = result.Entries[:max]
	}
	printJSON(result)
}

func cmdCreate() {
	args := os.Args[2:]
	if err := validateFlags(args,
		[]string{"--title", "--tags", "--dir", "--content"}, nil); err != nil {
		fatal(err.Error())
	}
	title := getFlag(args, "--title", "")
	if title == "" {
		fatal("usage: denotecli create --title <title> [--tags TAG,...] [--dir DIR] [--content TEXT]")
	}
	tagsStr := getFlag(args, "--tags", "")
	dir := getFlag(args, "--dir", "~/org/notes")
	content := getFlag(args, "--content", "")

	var tags []string
	if tagsStr != "" {
		for _, t := range strings.Split(tagsStr, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, sanitizeTag(t))
			}
		}
	}

	path, err := CreateNote(dir, title, tags, content)
	if err != nil {
		fatal(err.Error())
	}

	// Read back the created file to return structured JSON
	df, ok := ParseFilename(filepath.Base(path))
	if !ok {
		fatal("parse created filename: " + path)
	}
	df.Path = path
	printJSON(df)
}

func cmdSearchContent() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli search-content <query> [--dirs DIR,...] [--tags TAG] [--max N] [--matches N]")
	}
	query := os.Args[2]
	args := os.Args[3:]
	if err := validateFlags(args,
		[]string{"--dirs", "--tags", "--max", "--matches"}, nil); err != nil {
		fatal(err.Error())
	}
	dirsStr := getFlag(args, "--dirs", "~/org")
	tagFilter := getFlag(args, "--tags", "")
	maxStr := getFlag(args, "--max", "20")
	matchesStr := getFlag(args, "--matches", "3")
	max, _ := strconv.Atoi(maxStr)
	matches, _ := strconv.Atoi(matchesStr)
	if max <= 0 {
		max = 20
	}
	if matches <= 0 {
		matches = 3
	}

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	results := SearchContent(files, query, tagFilter, max, matches)
	printJSON(results)
}

func cmdSearchHeadings() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli search-headings <query> [--dirs DIR,...] [--tags TAG] [--level N] [--max N]")
	}
	query := os.Args[2]
	args := os.Args[3:]
	if err := validateFlags(args,
		[]string{"--dirs", "--tags", "--level", "--max"}, nil); err != nil {
		fatal(err.Error())
	}
	dirsStr := getFlag(args, "--dirs", "~/org")
	tagFilter := getFlag(args, "--tags", "")
	levelStr := getFlag(args, "--level", "0")
	maxStr := getFlag(args, "--max", "20")
	level, _ := strconv.Atoi(levelStr)
	max, _ := strconv.Atoi(maxStr)
	if max <= 0 {
		max = 20
	}

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	results := SearchHeadings(files, query, tagFilter, level, max)
	printJSON(results)
}

func cmdRead() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli read <id> [--dirs DIR,...] [--offset N] [--limit N] [--outline]")
	}
	id := os.Args[2]
	args := os.Args[3:]
	if err := validateFlags(args,
		[]string{"--dirs", "--offset", "--limit", "--level"},
		[]string{"--outline"}); err != nil {
		fatal(err.Error())
	}
	dirsStr := getFlag(args, "--dirs", "~/org")
	outline := hasFlag(args, "--outline")

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)

	if outline {
		levelStr := getFlag(args, "--level", "0")
		level, _ := strconv.Atoi(levelStr)
		do, err := ReadOutlineByID(files, id)
		if err != nil {
			fatal(err.Error())
		}
		do.Outline = FilterOutlineByLevel(do.Outline, level)
		printJSON(do)
		return
	}

	offsetStr := getFlag(args, "--offset", "0")
	limitStr := getFlag(args, "--limit", "0")
	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	dc, err := ReadByID(files, id, offset, limit)
	if err != nil {
		fatal(err.Error())
	}
	printJSON(dc)
}

func cmdRenameTag() {
	args := os.Args[2:]
	if err := validateFlags(args,
		[]string{"--from", "--to", "--dirs"},
		[]string{"--dry-run"}); err != nil {
		fatal(err.Error())
	}
	oldTag := getFlag(args, "--from", "")
	newTag := getFlag(args, "--to", "")
	if oldTag == "" || newTag == "" {
		fatal("usage: denotecli rename-tag --from <old> --to <new> [--dirs DIR,...] [--dry-run]")
	}
	dirsStr := getFlag(args, "--dirs", "~/org")
	dryRun := hasFlag(args, "--dry-run")

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	result, err := RenameTag(files, oldTag, newTag, dryRun)
	if err != nil {
		fatal(err.Error())
	}
	printJSON(result)
}

func cmdTags() {
	args := os.Args[2:]
	if err := validateFlags(args,
		[]string{"--dirs", "--pattern", "--top"},
		[]string{"--suggest"}); err != nil {
		fatal(err.Error())
	}
	dirsStr := getFlag(args, "--dirs", "~/org")
	suggest := hasFlag(args, "--suggest")

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)

	if suggest {
		result := SuggestTagCleanup(files)
		printJSON(result)
		return
	}

	pattern := getFlag(args, "--pattern", "")
	topStr := getFlag(args, "--top", "50")
	top, _ := strconv.Atoi(topStr)
	if top <= 0 {
		top = 50
	}
	stats := CollectTags(files, pattern, top)
	printJSON(stats)
}

func resolveDayDate(date, yearsAgo, daysAgo string) (string, error) {
	now := time.Now()
	if yearsAgo != "" {
		n, err := strconv.Atoi(yearsAgo)
		if err != nil || n < 0 {
			return "", fmt.Errorf("invalid --years-ago: %s", yearsAgo)
		}
		return now.AddDate(-n, 0, 0).Format("2006-01-02"), nil
	}
	if daysAgo != "" {
		n, err := strconv.Atoi(daysAgo)
		if err != nil || n < 0 {
			return "", fmt.Errorf("invalid --days-ago: %s", daysAgo)
		}
		return now.AddDate(0, 0, -n).Format("2006-01-02"), nil
	}
	if date == "" {
		return now.Format("2006-01-02"), nil
	}
	// YYYYMMDD → YYYY-MM-DD
	if len(date) == 8 && !strings.Contains(date, "-") {
		return date[:4] + "-" + date[4:6] + "-" + date[6:8], nil
	}
	// YYYYMMDDT... (Denote ID)
	if len(date) >= 8 && strings.Contains(date, "T") {
		d := date[:strings.Index(date, "T")]
		if len(d) == 8 {
			return d[:4] + "-" + d[4:6] + "-" + d[6:8], nil
		}
	}
	if len(date) == 10 && date[4] == '-' && date[7] == '-' {
		return date, nil
	}
	return "", fmt.Errorf("unrecognized date format: %s", date)
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

// validateFlags errors if args contains any --flag not in valueFlags or boolFlags.
// valueFlags take the next token as value (e.g. --max 20).
// boolFlags are standalone (e.g. --dry-run).
// This catches typos like --tag (vs --tags) and --limit (vs --max)
// that would otherwise be silently ignored by getFlag/hasFlag, leading to
// wrong defaults masquerading as filtered results. See llmlog 20260512T102541.
func validateFlags(args []string, valueFlags, boolFlags []string) error {
	val := make(map[string]bool, len(valueFlags))
	for _, f := range valueFlags {
		val[f] = true
	}
	bf := make(map[string]bool, len(boolFlags))
	for _, f := range boolFlags {
		bf[f] = true
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "--") {
			continue
		}
		if val[a] {
			i++ // skip flag value
			continue
		}
		if bf[a] {
			continue
		}
		return fmt.Errorf("unknown flag: %s", a)
	}
	return nil
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
	fmt.Fprintf(os.Stderr, `denotecli %s - Denote knowledge base CLI for AI agents

Usage:
  denotecli day [DATE] [--dirs DIR,...] [--years-ago N] [--days-ago N]
  denotecli timeline-journal [--month YYYY-MM] [--from DATE] [--to DATE] [--dirs DIR,...]
  denotecli search <query> [--tags TAG] [--dirs DIR,...] [--title-only] [--max N]
  denotecli search-content <query> [--dirs DIR,...] [--max N] [--matches N]
  denotecli search-headings <query> [--dirs DIR,...] [--level N] [--max N]
  denotecli graph <id> [--dirs DIR,...]
  denotecli keyword-map [query] [--dirs DIR,...] [--max N]
  denotecli create --title <title> [--tags TAG,...] [--dir DIR] [--content TEXT]
  denotecli read <id> [--dirs DIR,...] [--offset N] [--limit N] [--outline]
  denotecli rename-tag --from <old> --to <new> [--dirs DIR,...] [--dry-run]
  denotecli tags [--pattern PAT] [--top N] [--suggest] [--dirs DIR,...]

Options:
  --dirs DIR,...    Search directories, comma-separated (default: ~/org)
  --tags TAG        Filter by tag (comma-separated)
  --title-only      Search title only (search command)
  --max N           Max results (default: 20)
  --outline         Show heading structure only (read command)
  --matches N       Max matches per file (search-content, default: 3)
  --level N         Max heading level (outline/search-headings, 0=all)
  --offset N        Start line (read command)
  --limit N         Lines to read (read command, 0=all)
  --pattern PAT     Tag name regex filter (tags command)
  --top N           Top N tags (default: 50)
`, Version)
}
