package matcher

import (
	"testing"
)

func TestRatioMultipleCases(t *testing.T) {
	testCases := []struct {
		rawName               string
		other                 string
		expectedRatioTreshold float64
	}{
		{"pianino", "grand piano", 0.6},
		{"fortepian", "concert piano", 0.5},
		{"flet", "flauto", 0.5},
		{"skrzypce", "solo violin", 0.6},
		{"wiolonczela", "cello part", 0.55},
		{"gitara basowa", "bass guitar", 0.8},
		{"gitara", "acoustic guitar solo", 0.55},
		{"klarnet basowy", "bass clarinet", 0.8},
		{"klarnet", "clarinet section", 0.55},
		{"rożek angielski", "english horn", 0.8},
		{"obój", "oboe solo", 0.6},
		{"waltornia", "french horn", 0.5},
		{"tuba", "tuba player", 0.7},
		{"perkusja", "drum set", 0.5},
		{"trombone", "20 sen o Warszawie - Tromboneee 1", 0.7},
		{"saksofon alt", "sax altowy", 0.5},
		{"saksofon tenor", "tenor saxophone", 0.8},
		{"saksofon baryton", "baritone sax", 0.8},
	}

	for _, tc := range testCases {
		err, instrument := NewInstrument(tc.rawName)
		if err != nil {
			t.Errorf("Expected no error for '%s', got '%s'", tc.rawName, err)
			continue
		}
		ratio := instrument.GetSimilarityRatio(tc.other)
		t.Logf("Similarity ratio between '%s' and '%s': %f", instrument.NameEng, tc.other, ratio)
		if ratio < tc.expectedRatioTreshold {
			t.Errorf("Expected ratio > %f for '%s' and '%s', got %f", tc.expectedRatioTreshold, instrument.NameEng, tc.other, ratio)
		}
	}
}

func TestGetPairs(t *testing.T) {
	rawName := "saksofon alt in B"
	pairs := getAllPairs(rawName)
	expectedPairs := []string{"saksofon alt", "alt in", "in B"}
	if len(pairs) != len(expectedPairs) {
		t.Errorf("Expected %d pairs, got %d", len(expectedPairs), len(pairs))
	}
	for i, pair := range pairs {
		if pair != expectedPairs[i] {
			t.Errorf("Expected pair '%s', got '%s'", expectedPairs[i], pair)
		}
	}
}

func TestGetPairsSingleWord(t *testing.T) {
	rawName := "piano"
	pairs := getAllPairs(rawName)
	if len(pairs) != 1 || pairs[0] != "piano" {
		t.Errorf("Expected single pair 'piano', got %v", pairs)
	}
}

func TestGetMostSimilarInstrumentCoversAllInstruments(t *testing.T) {
	for _, expected := range instruments {
		err, got := GetMostSimilarInstrument(expected.NamePol)
		if err != nil {
			t.Fatalf("expected no error for polish input '%s', got '%v'", expected.NamePol, err)
		}

		if got.NamePol != expected.NamePol || got.NameEng != expected.NameEng {
			t.Errorf(
				"for polish input '%s' expected (%s, %s), got (%s, %s)",
				expected.NamePol,
				expected.NamePol,
				expected.NameEng,
				got.NamePol,
				got.NameEng,
			)
		}

		err, got = GetMostSimilarInstrument(expected.NameEng)
		if err != nil {
			t.Fatalf("expected no error for english input '%s', got '%v'", expected.NameEng, err)
		}

		if got.NameEng != expected.NameEng {
			t.Errorf(
				"for english input '%s' expected english name '%s', got '%s'",
				expected.NameEng,
				expected.NameEng,
				got.NameEng,
			)
		}
	}
}

func TestGetMostSimilarInstrumentWithTypos(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"gitara", "guitar"},
		{"gitra", "guitar"},
		{"saksofon altowy", "alto saxophone"},
		{"saksofon alt", "alto saxophone"},
		{"saksofon tenorowy", "tenor saxophone"},
		{"saksofon tenor", "tenor saxophone"},
		{"sax alt", "alto saxophone"},
		{"sax tenor", "tenor saxophone"},
		{"sax sopran", "soprano saxophone"},
		{"sax baryton", "baritone saxophone"},
		{"gitara basowa", "bass guitar"},
	}

	for _, tc := range testCases {
		err, got := GetMostSimilarInstrument(tc.input)
		if err != nil {
			t.Fatalf("expected no error for input '%s', got '%v'", tc.input, err)
		}

		if got.NameEng != tc.expected {
			t.Errorf(
				"for input '%s' expected english name '%s', got '%s'",
				tc.input,
				tc.expected,
				got.NameEng,
			)
		}
	}
}
