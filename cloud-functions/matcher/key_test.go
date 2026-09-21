package matcher

import "testing"

func TestParseKey(t *testing.T) {
	testCases := []struct {
		rawName  string
		expected string
	}{
		{"Klarnet 1 in B.pdf", "Bb"},
		{"Trąbka 2 in E.pdf", "E"},
		{"Waltornia 1 in F.pdf", "F"},
		{"Waltornia 1 in G.pdf", "G"},
		{"Klarnet_1 in A.pdf", "A"},
		{"horn in f.pdf", "F"},
		{"05  Ding Dong KAYAH - Clarinet 1 in Bb.pdf", "Bb"},
		{"Bb_Trumpet_1.pdf", "Bb"},
		{"_Gloria_-Saksofon_alt_1_Eb.pdf", "Eb"},
		{"_A_wczora_z_wieczora_-Clarinetto_in_Bb_1.pdf", "Bb"},
		{"Horn 1 (F).pdf", "F"},
		{"Clarinet 1 (C).pdf", "C"},
		{"happy-xmas_f-horn.pdf", "F"},
		{"_Gloria_-Puzon_1_in_C.pdf", "C"},
		// No key present.
		{"Gore Gwiazda Jezusowi.pdf", ""},
		{"bassoon.pdf", ""},
		{"drum.pdf", ""},
		{"_Gloria_A3.pdf", ""},
		{"_A_wczora_z_wieczora_A4.pdf", ""},
		{"_A_wczora_z_wieczora_-Piano.pdf", ""},
		{"Radetzky_Score.pdf", ""},
		{"Skrzypce 1.pdf", ""},
		{"Kontrabas.pdf", ""},
	}

	for _, tc := range testCases {
		got := ParseKey(tc.rawName)
		if got != tc.expected {
			t.Errorf("ParseKey(%q) = %q, want %q", tc.rawName, got, tc.expected)
		}
	}
}
