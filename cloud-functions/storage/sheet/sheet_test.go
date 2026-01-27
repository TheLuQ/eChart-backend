package sheet

import "testing"

func TestParsingFileName(t *testing.T) {
	t.Run("Parse file name", func(t *testing.T) {
		output, err := ParseFileName("Clarinet 2 in B")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.Name != "Clarinet" || output.Number != "2" || output.Key != "B" {
			t.Errorf("Expected Clarinet 2 in B, got %+v", output)
		}
	})

	t.Run("Parse file name (additional cases)", func(t *testing.T) {
		cases := []struct {
			name     string
			input    string
			expected Instrument
			wantErr  bool
		}{
			{
				name:  "two-part instrument name",
				input: "Bass Clarinet 2 in Bb",
				expected: Instrument{
					Name:   "Bass Clarinet",
					Number: "2",
					Key:    "Bb",
				},
			},
			{
				name:  "instrument only",
				input: "Flute",
				expected: Instrument{
					Name:   "Flute",
					Number: "",
					Key:    "",
				},
			},
			{
				name:  "instrument and number",
				input: "Oboe 1",
				expected: Instrument{
					Name:   "Oboe",
					Number: "1",
					Key:    "",
				},
			},
			{
				name:  "case-insensitive extension and key",
				input: "clarinet 2 in b",
				expected: Instrument{
					Name:   "clarinet",
					Number: "2",
					Key:    "b",
				},
			},
			{
				name:  "extra whitespace",
				input: "  Trumpet   10   in   Eb  ",
				expected: Instrument{
					Name:   "Trumpet",
					Number: "10",
					Key:    "Eb",
				},
			},
			{
				name:  "no space before number",
				input: "Clarinet2 in B",
				expected: Instrument{
					Name:   "Clarinet",
					Number: "2",
					Key:    "B",
				},
			},
			{
				name:    "reject empty string",
				input:   "",
				wantErr: true,
			},
			{
				name:    "reject wrong extension",
				input:   "Clarinet 2 in B.txt",
				wantErr: true,
			},
			{
				name:    "reject missing 'in' keyword",
				input:   "Clarinet 2 B",
				wantErr: true,
			},
			{
				name:    "reject number 0",
				input:   "Clarinet 0 in B",
				wantErr: true,
			},
			{
				name:    "reject non-anchored suffix",
				input:   "Clarinet 2 in B.pdfx",
				wantErr: true,
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				got, err := ParseFileName(tc.input)
				if tc.wantErr {
					if err == nil {
						t.Fatalf("expected error, got nil; output=%+v", got)
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tc.expected {
					t.Fatalf("expected %+v, got %+v", tc.expected, got)
				}
			})
		}
	})
}
