package firestore

import (
	"encoding/json"
	"os"
	"testing"
)

func TestToSheet(t *testing.T) {
	t.Run("with complete metadata", func(t *testing.T) {
		rawPath := "my-band/sheets/my-song"
		metadata := map[string]string{
			"instrument_name_en":  "guitar",
			"instrument_name_pol": "gitara",
			"voice":               "1",
			"key":                 "C",
		}

		sheet, err := ToSheet(rawPath, metadata["instrument_name_en"], metadata["instrument_name_pol"], metadata["voice"], metadata["key"])
		if err != nil {
			t.Fatalf("Expected no error but got %v", err)
		}

		expectedID := "my-band/sheets/my-song"
		if sheet.Id != expectedID {
			t.Errorf("Expected sheet id %q but got %q", expectedID, sheet.Id)
		}
		if sheet.FileName != "my-song" {
			t.Errorf("Expected file name 'my-song' but got %s", sheet.FileName)
		}
		if sheet.Title != "sheets" {
			t.Errorf("Expected title 'sheets' but got %s", sheet.Title)
		}
		if sheet.Band != "my-band" {
			t.Errorf("Expected band 'my-band' but got %s", sheet.Band)
		}
		if sheet.Instrument.Name != "guitar" {
			t.Errorf("Expected instrument name 'guitar' but got %s", sheet.Instrument.Name)
		}
		if sheet.Instrument.NamePol != "gitara" {
			t.Errorf("Expected instrument polish name 'gitara' but got %s", sheet.Instrument.NamePol)
		}
		if sheet.Instrument.Voice != "1" {
			t.Errorf("Expected voice '1' but got %s", sheet.Instrument.Voice)
		}
		if sheet.Instrument.Key != "C" {
			t.Errorf("Expected key 'C' but got %s", sheet.Instrument.Key)
		}
	})
}

func TestSheetMarshalJSON(t *testing.T) {
	sheet := Sheet{
		Instrument: Instrument{
			Name:    "guitar",
			NamePol: "gitara",
			Voice:   "1",
			Key:     "C",
		},
		Id:       "sheet-123",
		Title:    "my song",
		FileName: "asd.pdf",
		Url:      "https://storage.googleapis.com/my-bucket/music/sheets/my-song",
	}

	data, err := json.Marshal(sheet)
	if err != nil {
		t.Fatalf("Expected no error marshaling Sheet but got %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Expected valid JSON but got error: %v", err)
	}

	checks := map[string]string{
		"instrument_name_en":  "guitar",
		"instrument_name_pol": "gitara",
		"key":                 "C",
		"file_id":             "sheet-123",
		"title":               "my song",
		"file_name":           "asd.pdf",
		"voice":               "1",
		"url":                 "https://storage.googleapis.com/my-bucket/music/sheets/my-song",
	}
	for key, want := range checks {
		got, ok := result[key]
		if !ok {
			t.Errorf("Expected JSON key %q to be present", key)
			continue
		}
		if got != want {
			t.Errorf("Expected JSON[%q] = %q but got %q", key, want, got)
		}
	}
}

func TestGetSheets(t *testing.T) {
	t.Setenv(sheetStorageBaseURLEnv, "https://storage.googleapis.com")
	t.Setenv(sheetStorageBucketEnv, "my-bucket")

	sg := &SheetGroup{
		Title: "group title",
		Band:  "my-band",
		Sheets: []Sheet{
			{Id: "my-band/sheets/my-song"},
			{},
		},
	}

	result := sg.GetSheets()
	if len(result) != 2 {
		t.Fatalf("Expected 2 sheets but got %d", len(result))
	}
	if result[0].Title != "group title" {
		t.Fatalf("Expected title to be propagated but got %q", result[0].Title)
	}
	if result[0].Band != "my-band" {
		t.Fatalf("Expected band to be propagated but got %q", result[0].Band)
	}
	if result[0].Url != "https://storage.googleapis.com/storage/v1/b/my-bucket/o/my-band%2Fsheets%2Fmy-song?alt=media" {
		t.Fatalf("Expected generated url but got %q", result[0].Url)
	}
	if result[1].Url != "" {
		t.Fatalf("Expected empty url for empty id but got %q", result[1].Url)
	}
}

func TestGetSheetsWithoutBaseURL(t *testing.T) {
	if err := os.Unsetenv(sheetStorageBaseURLEnv); err != nil {
		t.Fatalf("failed to unset env: %v", err)
	}
	t.Setenv(sheetStorageBucketEnv, "my-bucket")

	sg := &SheetGroup{
		Title: "group title",
		Sheets: []Sheet{
			{Id: "my-band/sheets/my-song"},
		},
	}

	result := sg.GetSheets()
	if len(result) != 1 {
		t.Fatalf("Expected 1 sheet but got %d", len(result))
	}
	if result[0].Url != "https://storage.googleapis.com/storage/v1/b/my-bucket/o/my-band%2Fsheets%2Fmy-song?alt=media" {
		t.Fatalf("Expected empty url when %s is not set but got %q", sheetStorageBaseURLEnv, result[0].Url)
	}
}

func TestGroupKey(t *testing.T) {
	cases := []struct {
		name  string
		band  string
		title string
		want  string
	}{
		{name: "normal values", band: "Frank Sinatra Orchestra", title: "Fly Me to the Moon", want: "frank-sinatra-orchestra__fly-me-to-the-moon"},
		{name: "collapses spaces and symbols", band: "  Big   Band ", title: "Fly  me --- to   the moon", want: "big-band__fly-me-to-the-moon"},
		{name: "empty title", band: "Only Band", title: "", want: "only-band"},
		{name: "empty band", band: "", title: "Only Title", want: "only-title"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := GroupKey(tc.band, tc.title)
			if got != tc.want {
				t.Fatalf("Expected group key %q but got %q", tc.want, got)
			}
		})
	}
}

func TestSheetGroupAddSheetIfMissing(t *testing.T) {
	t.Run("adds when id is new", func(t *testing.T) {
		sg := &SheetGroup{Sheets: []Sheet{{Id: "a"}}}
		sg.AddSheetIfMissing(&Sheet{Id: "b"})
		if len(sg.Sheets) != 2 {
			t.Fatalf("Expected 2 sheets but got %d", len(sg.Sheets))
		}
		if sg.Sheets[1].Id != "b" {
			t.Fatalf("Expected appended sheet id %q but got %q", "b", sg.Sheets[1].Id)
		}
	})

	t.Run("does not add when id already exists", func(t *testing.T) {
		sg := &SheetGroup{Sheets: []Sheet{{Id: "dup", FileName: "old.pdf"}}}
		sg.AddSheetIfMissing(&Sheet{Id: "dup", FileName: "new.pdf"})
		if len(sg.Sheets) != 1 {
			t.Fatalf("Expected 1 sheet but got %d", len(sg.Sheets))
		}
		if sg.Sheets[0].FileName != "new.pdf" {
			t.Fatalf("Expected existing sheet to be replaced")
		}
	})

	t.Run("works when existing slice is nil", func(t *testing.T) {
		sg := &SheetGroup{}
		sg.AddSheetIfMissing(&Sheet{Id: "first"})
		if len(sg.Sheets) != 1 {
			t.Fatalf("Expected 1 sheet but got %d", len(sg.Sheets))
		}
	})
}
