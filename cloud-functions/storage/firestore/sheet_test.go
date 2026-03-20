package firestore

import (
	"encoding/json"
	"testing"
)

func TestToSheet(t *testing.T) {
	t.Run("with complete metadata", func(t *testing.T) {
		rawPath := "music/sheets/my-song"
		metadata := map[string]string{
			"instrument_name_en":  "guitar",
			"instrument_name_pol": "gitara",
			"voice":               "1",
			"key":                 "C",
		}

		sheet, err := ToSheet(rawPath, metadata)
		if err != nil {
			t.Fatalf("Expected no error but got %v", err)
		}

		if sheet.Id != rawPath {
			t.Errorf("Expected sheet id %q but got %q", rawPath, sheet.Id)
		}
		if sheet.FileName != "my-song" {
			t.Errorf("Expected file name 'my-song' but got %s", sheet.FileName)
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

	t.Run("with missing metadata keys", func(t *testing.T) {
		rawPath := "music/sheets/another-song"

		sheet, err := ToSheet(rawPath, map[string]string{})
		if err != nil {
			t.Fatalf("Expected no error but got %v", err)
		}
		if sheet.Instrument.Name != "" {
			t.Errorf("Expected empty instrument name but got %s", sheet.Instrument.Name)
		}
		if sheet.Instrument.NamePol != "" {
			t.Errorf("Expected empty instrument polish name but got %s", sheet.Instrument.NamePol)
		}
		if sheet.Instrument.Voice != "" {
			t.Errorf("Expected empty voice but got %s", sheet.Instrument.Voice)
		}
		if sheet.Instrument.Key != "" {
			t.Errorf("Expected empty key but got %s", sheet.Instrument.Key)
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
