// stemmer.go — Minimal Porter Stemmer for English tag normalization.
// Based on Martin Porter's algorithm (1980).
// Only handles the most common cases relevant to tag cleanup.
package main

import "strings"

// Stem returns the stem of an English word using simplified Porter rules.
func Stem(word string) string {
	w := strings.ToLower(word)
	if len(w) <= 2 {
		return w
	}

	// Step 1a: plurals
	if strings.HasSuffix(w, "sses") {
		w = w[:len(w)-2]
	} else if strings.HasSuffix(w, "ies") {
		w = w[:len(w)-2]
	} else if strings.HasSuffix(w, "ss") {
		// keep as is
	} else if strings.HasSuffix(w, "s") && len(w) > 3 {
		w = w[:len(w)-1]
	}

	// Step 1b: -ed, -ing
	if strings.HasSuffix(w, "eed") {
		if measure(w[:len(w)-3]) > 0 {
			w = w[:len(w)-1] // eed → ee
		}
	} else if strings.HasSuffix(w, "ed") && hasVowel(w[:len(w)-2]) {
		w = w[:len(w)-2]
		w = step1bCleanup(w)
	} else if strings.HasSuffix(w, "ing") && hasVowel(w[:len(w)-3]) {
		w = w[:len(w)-3]
		w = step1bCleanup(w)
	}

	// Step 1c: y → i
	if strings.HasSuffix(w, "y") && hasVowel(w[:len(w)-1]) && len(w) > 2 {
		w = w[:len(w)-1] + "i"
	}

	// Step 2: common suffixes
	step2 := [][2]string{
		{"ational", "ate"}, {"tional", "tion"}, {"enci", "ence"},
		{"anci", "ance"}, {"izer", "ize"}, {"abli", "able"},
		{"alli", "al"}, {"entli", "ent"}, {"eli", "e"},
		{"ousli", "ous"}, {"ization", "ize"}, {"ation", "ate"},
		{"ator", "ate"}, {"alism", "al"}, {"iveness", "ive"},
		{"fulness", "ful"}, {"ousness", "ous"}, {"aliti", "al"},
		{"iviti", "ive"}, {"biliti", "ble"}, {"logi", "log"},
	}
	for _, pair := range step2 {
		if strings.HasSuffix(w, pair[0]) {
			stem := w[:len(w)-len(pair[0])]
			if measure(stem) > 0 {
				w = stem + pair[1]
			}
			break
		}
	}

	// Step 3: more suffixes
	step3 := [][2]string{
		{"icate", "ic"}, {"ative", ""}, {"alize", "al"},
		{"iciti", "ic"}, {"ical", "ic"}, {"ful", ""},
		{"ness", ""},
	}
	for _, pair := range step3 {
		if strings.HasSuffix(w, pair[0]) {
			stem := w[:len(w)-len(pair[0])]
			if measure(stem) > 0 {
				w = stem + pair[1]
			}
			break
		}
	}

	// Step 4: remove known suffixes if measure > 1
	step4suffixes := []string{
		"al", "ance", "ence", "er", "ic", "able", "ible", "ant",
		"ement", "ment", "ent", "ion", "ou", "ism", "ate", "iti",
		"ous", "ive", "ize",
	}
	for _, suf := range step4suffixes {
		if strings.HasSuffix(w, suf) {
			stem := w[:len(w)-len(suf)]
			if measure(stem) > 1 {
				if suf == "ion" {
					if len(stem) > 0 && (stem[len(stem)-1] == 's' || stem[len(stem)-1] == 't') {
						w = stem
					}
				} else {
					w = stem
				}
			}
			break
		}
	}

	// Step 5a: remove trailing e
	if strings.HasSuffix(w, "e") {
		stem := w[:len(w)-1]
		if measure(stem) > 1 || (measure(stem) == 1 && !endsCVC(stem)) {
			w = stem
		}
	}

	// Step 5b: double l
	if strings.HasSuffix(w, "ll") && measure(w[:len(w)-1]) > 1 {
		w = w[:len(w)-1]
	}

	return w
}

func isVowel(c byte) bool {
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u'
}

func hasVowel(s string) bool {
	for i := 0; i < len(s); i++ {
		if isVowel(s[i]) {
			return true
		}
	}
	return false
}

// measure counts the number of VC (vowel-consonant) sequences.
func measure(s string) int {
	if len(s) == 0 {
		return 0
	}
	m := 0
	i := 0
	// skip initial consonants
	for i < len(s) && !isVowel(s[i]) {
		i++
	}
	for i < len(s) {
		// skip vowels
		for i < len(s) && isVowel(s[i]) {
			i++
		}
		if i < len(s) {
			m++
			// skip consonants
			for i < len(s) && !isVowel(s[i]) {
				i++
			}
		}
	}
	return m
}

func endsCVC(s string) bool {
	if len(s) < 3 {
		return false
	}
	c1, v, c2 := s[len(s)-3], s[len(s)-2], s[len(s)-1]
	if !isVowel(c1) && isVowel(v) && !isVowel(c2) {
		return c2 != 'w' && c2 != 'x' && c2 != 'y'
	}
	return false
}

func step1bCleanup(w string) string {
	if strings.HasSuffix(w, "at") || strings.HasSuffix(w, "bl") || strings.HasSuffix(w, "iz") {
		return w + "e"
	}
	if len(w) >= 2 && w[len(w)-1] == w[len(w)-2] {
		last := w[len(w)-1]
		if last != 'l' && last != 's' && last != 'z' {
			return w[:len(w)-1]
		}
	}
	if measure(w) == 1 && endsCVC(w) {
		return w + "e"
	}
	return w
}
