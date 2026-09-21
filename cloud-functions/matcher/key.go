package matcher

import (
	"regexp"
	"strings"
)

var keyTokenSeparatorRegex = regexp.MustCompile(`[^A-Za-z0-9]+`)

// keyNormalization maps a recognized key token (lower-cased) to its canonical
// English form. German/Polish spellings are converted: "b" -> "Bb" (flat B),
// "es" -> "Eb", "h" -> "B" (natural B).
var keyNormalization = map[string]string{
	"a":  "A",
	"b":  "Bb",
	"bb": "Bb",
	"h":  "B",
	"c":  "C",
	"d":  "D",
	"e":  "E",
	"eb": "Eb",
	"es": "Eb",
	"f":  "F",
	"g":  "G",
}

// ParseKey extracts the musical key from a raw file name.
//
// The name is split into alphanumeric tokens (every non-alphanumeric run
// becomes a separator). Whole-token matching avoids substring false positives
// such as "Gore" -> "G" or "bassoon" -> "B", and keeps runs like "A3" or
// "IMSLP172451" intact so they never match a bare key letter.
//
// Selection priority:
//  1. the token immediately after an "in" token (e.g. "in Bb", "horn in f");
//  2. any two-letter key (Bb, Eb, Es);
//  3. a single-letter key, ignoring one that appears as the very first token
//     (which is typically a title word, e.g. the leading "A" in "A wczora...").
//
// The recognized key is normalized to its canonical English form. An empty
// string is returned when no key is found.
func ParseKey(rawName string) string {
	spaced := keyTokenSeparatorRegex.ReplaceAllString(rawName, " ")
	tokens := strings.Fields(spaced)

	// 1. Token following an "in" marker.
	for i := 0; i < len(tokens)-1; i++ {
		if strings.EqualFold(tokens[i], "in") {
			if key, ok := keyNormalization[strings.ToLower(tokens[i+1])]; ok {
				return key
			}
		}
	}

	// 2. Any two-letter key.
	for _, token := range tokens {
		if len(token) != 2 {
			continue
		}
		if key, ok := keyNormalization[strings.ToLower(token)]; ok {
			return key
		}
	}

	// 3. Single-letter key, skipping the very first token.
	for i, token := range tokens {
		if i == 0 || len(token) != 1 {
			continue
		}
		if key, ok := keyNormalization[strings.ToLower(token)]; ok {
			return key
		}
	}

	return ""
}
