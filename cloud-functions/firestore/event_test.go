package firestore

import "testing"

func TestNormalizePaginationLimit(t *testing.T) {
	cases := []struct {
		name  string
		input int
		want  int
	}{
		{name: "uses default for values below min", input: 0, want: defaultPaginationLimit},
		{name: "uses min accepted value", input: minPaginationLimit, want: minPaginationLimit},
		{name: "keeps in-range limit", input: 42, want: 42},
		{name: "caps values above max", input: maxPaginationLimit + 1, want: maxPaginationLimit},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizePaginationLimit(tc.input)
			if got != tc.want {
				t.Fatalf("expected %d but got %d", tc.want, got)
			}
		})
	}
}
