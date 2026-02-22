// main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const Version = "0.4.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
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

func cmdKeywordMap() {
	args := os.Args[2:]
	query := ""
	if len(os.Args) >= 3 && !strings.HasPrefix(os.Args[2], "--") {
		query = os.Args[2]
		args = os.Args[3:]
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
		fatal("usage: denotecli search-content <query> [--dirs DIR,...] [--max N] [--matches N]")
	}
	query := os.Args[2]
	args := os.Args[3:]
	dirsStr := getFlag(args, "--dirs", "~/org")
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
	results := SearchContent(files, query, max, matches)
	printJSON(results)
}

func cmdSearchHeadings() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli search-headings <query> [--dirs DIR,...] [--level N] [--max N]")
	}
	query := os.Args[2]
	args := os.Args[3:]
	dirsStr := getFlag(args, "--dirs", "~/org")
	levelStr := getFlag(args, "--level", "0")
	maxStr := getFlag(args, "--max", "20")
	level, _ := strconv.Atoi(levelStr)
	max, _ := strconv.Atoi(maxStr)
	if max <= 0 {
		max = 20
	}

	dirs := strings.Split(dirsStr, ",")
	files := ScanDirs(dirs)
	results := SearchHeadings(files, query, level, max)
	printJSON(results)
}

func cmdRead() {
	if len(os.Args) < 3 {
		fatal("usage: denotecli read <id> [--dirs DIR,...] [--offset N] [--limit N] [--outline]")
	}
	id := os.Args[2]
	args := os.Args[3:]
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
	fmt.Fprintf(os.Stderr, `denotecli %s - Denote knowledge base CLI for AI agents

Usage:
  denotecli search <query> [--tags TAG] [--dirs DIR,...] [--title-only] [--max N]
  denotecli search-content <query> [--dirs DIR,...] [--max N] [--matches N]
  denotecli search-headings <query> [--dirs DIR,...] [--level N] [--max N]
  denotecli keyword-map [query] [--dirs DIR,...] [--max N]
  denotecli create --title <title> [--tags TAG,...] [--dir DIR] [--content TEXT]
  denotecli read <id> [--dirs DIR,...] [--offset N] [--limit N] [--outline]
  denotecli tags [--pattern PAT] [--top N] [--dirs DIR,...]

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
