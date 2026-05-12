// json_contract_test.go
//
// Regression tests for two silent-failure bugs found while cleaning up
// the fleeting tag pattern on 2026-05-12:
//
//   Bug 2 — unknown flag silently ignored (getFlag returns default)
//   Bug 3 — empty result slice marshals to JSON null (not [])
//
// Forensic: ~/sync/org/llmlog/20260512T102541
//
// Lock the new contract:
//   - All search-like commands return non-nil slices even when empty.
//   - validateFlags rejects any --flag not in its known set.

package main

import (
	"encoding/json"
	"testing"
)

// --- Bug 3 contract: search results are always a valid array ---

func TestSearchReturnsEmptyArrayNotNil(t *testing.T) {
	// No files → no matches.
	results := Search(nil, "anything", "", false, 20)
	assertEmptyArrayContract(t, "Search", results)
}

func TestSearchHeadingsReturnsEmptyArrayNotNil(t *testing.T) {
	results := SearchHeadings(nil, "anything", "", 0, 20)
	assertEmptyArrayContract(t, "SearchHeadings", results)
}

func TestSearchContentReturnsEmptyArrayNotNil(t *testing.T) {
	results := SearchContent(nil, "anything", "", 20, 3)
	assertEmptyArrayContract(t, "SearchContent (no files)", results)

	// Also: empty query path — historically returned nil, now must be [].
	results2 := SearchContent(nil, "", "", 20, 3)
	assertEmptyArrayContract(t, "SearchContent (empty query)", results2)
}

func TestExtractOutlineEmptyReturnsEmptyArray(t *testing.T) {
	outline := ExtractOutline("")
	assertEmptyArrayContract(t, "ExtractOutline", outline)
}

// assertEmptyArrayContract verifies a result satisfies the agent JSON contract:
// (1) non-nil, (2) marshals to "[]" not "null".
func assertEmptyArrayContract[T any](t *testing.T, label string, results []T) {
	t.Helper()
	if results == nil {
		t.Errorf("%s: result is nil; agent callers will see JSON null and len() crashes", label)
	}
	b, err := json.Marshal(results)
	if err != nil {
		t.Fatalf("%s: marshal failed: %v", label, err)
	}
	if string(b) != "[]" {
		t.Errorf("%s: marshalled to %q, want []", label, string(b))
	}
}

// --- Bug 2 contract: unknown flags are rejected ---

func TestValidateFlagsAcceptsKnown(t *testing.T) {
	args := []string{"--tags", "fleeting", "--max", "200", "--title-only"}
	err := validateFlags(args,
		[]string{"--tags", "--max"},
		[]string{"--title-only"})
	if err != nil {
		t.Errorf("known flags rejected: %v", err)
	}
}

func TestValidateFlagsRejectsTypo(t *testing.T) {
	// Real-world repro: --tag (no s) and --limit instead of --max.
	// Before fix: silently ignored, default --max 20 applied, caller misled.
	args := []string{"--tag", "fleeting", "--limit", "200"}
	err := validateFlags(args,
		[]string{"--tags", "--dirs", "--max"},
		[]string{"--title-only"})
	if err == nil {
		t.Errorf("expected error for --tag/--limit typos, got nil (silent ignore)")
	}
}

func TestValidateFlagsAllowsPositionals(t *testing.T) {
	// Non-flag positional args (e.g. query, date, id) must pass through.
	args := []string{"2026-05-12", "--dirs", "~/org"}
	err := validateFlags(args, []string{"--dirs"}, nil)
	if err != nil {
		t.Errorf("positional rejected: %v", err)
	}
}

func TestValidateFlagsEmpty(t *testing.T) {
	if err := validateFlags(nil, []string{"--x"}, nil); err != nil {
		t.Errorf("empty args rejected: %v", err)
	}
}
