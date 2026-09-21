package matcher

import (
	"reflect"
	"testing"
)

func TestParseVoice(t *testing.T) {
	testCases := []struct {
		rawName  string
		expected []int
	}{
		{"Flet 2.pdf", []int{2}},
		{"Habanera, fl1.pdf", []int{1}},
		{"Habanera, Cl2.pdf", []int{2}},
		{"Klarnet 1 in B.pdf", []int{1}},
		{"Waltornia 4 in F.pdf", []int{4}},
		{"Radetzky_Trombone1&2.pdf", []int{1, 2}},
		{"Oboe 1, 2.pdf", []int{1, 2}},
		{"Bassoons (1, 2).pdf", []int{1, 2}},
		{"Flety 1i 2.pdf", []int{1, 2}},
		{"1 i 2 Flute, Piccolo.pdf", []int{1, 2}},
		{"01  Ding Dong KAYAH - Flute 1.pdf", []int{1}},
		{"08  Ding Dong KAYAH - Alto Saxophone 1.pdf", []int{1}},
		{"25  Ding Dong KAYAH - Drums.pdf", []int{}},
		{"Cellos.pdf 20-55-41-767.pdf", []int{}},
		{"Partytura.pdf 22-29-25-508.pdf", []int{}},
		{"IMSLP172459-PMLP21292-Waldteufel-Patineurs-Violin1.pdf", []int{1}},
		{"IMSLP172453-PMLP21292-Waldteufel-Patineurs-Clarinets.pdf", []int{}},
		{"happy-xmas_percussion-1 (1).pdf", []int{1}},
		{"Piccolo.pdf", []int{}},
		{"_A_wczora_z_wieczora_-Clarinetto_in_Bb_3.pdf", []int{3}},
	}

	for _, tc := range testCases {
		got := ParseVoice(tc.rawName)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("ParseVoice(%q) = %v, want %v", tc.rawName, got, tc.expected)
		}
	}
}
