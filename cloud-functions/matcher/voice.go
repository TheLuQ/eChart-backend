package matcher

import (
	"regexp"
	"sort"
)

var voiceDigitsRegex = regexp.MustCompile(`[0-9]+`)

// ParseVoice extracts voice numbers (1-4) from a raw file name.
//
// It matches maximal digit runs and keeps only those whose exact text is
// "1", "2", "3" or "4". Comparing the whole run (not its integer value) means
// zero-padded track/order prefixes such as "01" or "02" are skipped, because
// "01" != "1". Timestamps ("20-55-41-767") and IMSLP catalogue numbers are
// multi-digit runs and are ignored for the same reason.
//
// Ranges like "1&2", "1, 2" or "1 i 2" naturally yield both numbers.
func ParseVoice(rawName string) []int {
	matches := voiceDigitsRegex.FindAllString(rawName, -1)

	seen := make(map[int]bool)
	voices := make([]int, 0, len(matches))
	for _, m := range matches {
		var voice int
		switch m {
		case "1":
			voice = 1
		case "2":
			voice = 2
		case "3":
			voice = 3
		case "4":
			voice = 4
		default:
			continue
		}
		if seen[voice] {
			continue
		}
		seen[voice] = true
		voices = append(voices, voice)
	}

	sort.Ints(voices)
	return voices
}
