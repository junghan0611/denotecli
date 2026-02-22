// tag_suggest_test.go
package main

import (
	"testing"
)

func TestIsSimilarPlural(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"apple", "apples", "plural"},
		{"process", "processes", "plural"},
		{"strategy", "strategies", "plural"},
		{"emacs", "vim", ""},
		{"tag", "tags", "plural"},
	}
	for _, c := range cases {
		got := isSimilar(c.a, c.b)
		if got != c.want {
			t.Errorf("isSimilar(%q, %q) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestIsSimilarDerivation(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"communication", "communicational", "derivation"},
		{"development", "developmental", "derivation"},
		// {"creative", "creativity"} — ive→ivity 패턴은 너무 드물어서 미지원
		{"socialism", "socialist", "derivation"},
	}
	for _, c := range cases {
		got := isSimilar(c.a, c.b)
		if got != c.want {
			t.Errorf("isSimilar(%q, %q) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestIsSimilarPrefix(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"commit", "commits", "plural"}, // caught by plural first
		{"graph", "graphdb", "prefix"},
		{"note", "notes", "plural"},
	}
	for _, c := range cases {
		got := isSimilar(c.a, c.b)
		if got != c.want {
			t.Errorf("isSimilar(%q, %q) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestSuggestTagCleanup(t *testing.T) {
	files := []DenoteFile{
		{ID: "1", Tags: []string{"emacs", "vim", "apple", "apples"}},
		{ID: "2", Tags: []string{"emacs", "apple", "communication", "communicational"}},
		{ID: "3", Tags: []string{"vim", "apples", "note", "notes"}},
	}

	result := SuggestTagCleanup(files)
	if len(result.Suggestions) == 0 {
		t.Fatal("expected suggestions")
	}

	// Should find apple/apples and note/notes at minimum
	found := map[string]bool{}
	for _, s := range result.Suggestions {
		found[s.Tag1+"/"+s.Tag2] = true
	}
	if !found["apple/apples"] {
		t.Error("missing apple/apples suggestion")
	}
	if !found["note/notes"] {
		t.Error("missing note/notes suggestion")
	}
}
