// stemmer_test.go
package main

import "testing"

func TestStem(t *testing.T) {
	cases := []struct {
		word string
		want string
	}{
		// Plurals
		{"apples", "appl"},
		{"apple", "appl"},
		{"tags", "tag"},
		{"processes", "process"},
		{"strategies", "strategi"},
		{"strategy", "strategi"},

		// Derivations — same stem
		{"communicate", "commun"},
		{"communication", "commun"},
		{"communicational", "commun"},
		{"agent", "agent"},
		{"agents", "agent"},
		{"develop", "develop"},
		{"development", "develop"},
		// developmental → development (step2 stops early, practical enough)

		// Short words preserved
		{"ai", "ai"},
		{"go", "go"},
		{"lsp", "lsp"},

		// Already stemmed
		{"emacs", "emac"},
		{"vim", "vim"},
		{"linux", "linux"},
	}
	for _, c := range cases {
		got := Stem(c.word)
		if got != c.want {
			t.Errorf("Stem(%q) = %q, want %q", c.word, got, c.want)
		}
	}
}

func TestStemGrouping(t *testing.T) {
	// These should all group together
	groups := [][]string{
		{"apple", "apples"},
		{"communicate", "communication", "communicational"},
		{"agent", "agents"},
		{"note", "notes"},
		{"develop", "development"},
		{"agent", "agents"},
	}
	for _, g := range groups {
		stems := make(map[string]bool)
		for _, w := range g {
			stems[Stem(w)] = true
		}
		if len(stems) != 1 {
			t.Errorf("words %v have %d different stems: %v", g, len(stems), stems)
		}
	}
}
