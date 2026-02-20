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
