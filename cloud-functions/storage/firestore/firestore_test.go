package firestore

import (
	"fmt"
	"testing"
)

func TestSaveSheet(t *testing.T) {
	var sdb2 InterfaceFirestoreDB = &mock{}
	sheet := Sheet{
		Instrument: "Piano",
	}
	err := sdb2.SaveSheet(sheet)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestNewSheet(t *testing.T) {
	t.Run("with valid path", func(t *testing.T) {
		fullPath := "music/sheets/my-song/guitar.pdf"
		sheet, err := ToSheet(fullPath)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if sheet.Instrument != "guitar" {
			t.Errorf("Expected instrument 'guitar' but got %s", sheet.Instrument)
		}
		if sheet.Title != "my song" {
			t.Errorf("Expected title 'my song' but got %s", sheet.Title)
		}
		if sheet.Id != "music/sheets/my-song" {
			t.Errorf("Expected id 'music/sheets/my-song' but got %s", sheet.Id)
		}
	})

	t.Run("with valid path without instrument", func(t *testing.T) {
		fullPath := "music/sheets/another-song/clarinet.pdf"
		sheet, err := ToSheet(fullPath)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if sheet.Instrument != "clarinet" {
			t.Errorf("Expected instrument 'clarinet' but got %s", sheet.Instrument)
		}
		if sheet.Title != "another song" {
			t.Errorf("Expected title 'another song' but got %s", sheet.Title)
		}
		if sheet.Id != "music/sheets/another-song" {
			t.Errorf("Expected id 'music/sheets/another-song' but got %s", sheet.Id)
		}
	})
}

func TestNewSheetInvalidPath(t *testing.T) {
	fullPath := "invalidpath"
	_, err := ToSheet(fullPath)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestAddSheet(t *testing.T) {
	sheet1 := Sheet{Instrument: "Piano"}
	sheet2 := Sheet{Instrument: "Guitar"}

	t.Run("adds sheet when absent", func(t *testing.T) {
		group := &SheetGroup{"My Group", "2026-01-15", []Sheet{}}
		group.AddSheet(sheet1)

		if len(group.Sheets) != 1 || group.Sheets[0] != sheet1 {
			t.Fatalf("Expected receiver to contain exactly %#v; receiver=%#v", sheet1, group.Sheets)
		}
	})

	t.Run("does not add duplicate", func(t *testing.T) {
		singleGroup := &SheetGroup{"My Group", "2026-01-15", []Sheet{sheet1}}

		singleGroup.AddSheet(sheet1)

		if len(singleGroup.Sheets) != 1 || singleGroup.Sheets[0] != sheet1 {
			t.Fatalf("Expected receiver to remain unchanged after duplicate add; receiver=%#v", singleGroup.Sheets)
		}
	})

	t.Run("adds different sheet", func(t *testing.T) {
		singleGroup := &SheetGroup{"My Group", "2026-01-15", []Sheet{sheet1}}
		singleGroup.AddSheet(sheet2)

		if len(singleGroup.Sheets) != 2 || singleGroup.Sheets[1] != sheet2 {
			t.Fatalf("Expected receiver to end with %#v appended; receiver=%#v", sheet2, singleGroup.Sheets)
		}
	})
}

func TestUpdateTitle(t *testing.T) {
	t.Run("does nothing when title is empty", func(t *testing.T) {
		group := &SheetGroup{"Original", "old", []Sheet{{Instrument: "Piano"}}}

		group.UpdateTitle("")

		if group.Title != "Original" {
			t.Fatalf("Expected title to remain %q; got %q", "Original", group.Title)
		}
		if group.LastUpdated != "old" {
			t.Fatalf("Expected LastUpdated to remain %q; got %q", "old", group.LastUpdated)
		}
	})

	t.Run("does nothing when title is unchanged", func(t *testing.T) {
		group := &SheetGroup{"Same", "old", []Sheet{{Instrument: "Piano"}}}

		group.UpdateTitle("Same")

		if group.Title != "Same" {
			t.Fatalf("Expected title to remain %q; got %q", "Same", group.Title)
		}
		if group.LastUpdated != "old" {
			t.Fatalf("Expected LastUpdated to remain %q; got %q", "old", group.LastUpdated)
		}
	})
}

type mock struct{}

func (m *mock) SaveSheet(sheet Sheet) error {
	fmt.Printf("Saving sheet for instrument %s\n", sheet.Instrument)
	return nil
}

func (m *mock) SearchGroupById(id string) (*SheetGroup, error) {
	group := &SheetGroup{
		Title:       "Sample Group",
		LastUpdated: "2024-01-01",
		Sheets: []Sheet{
			{Instrument: "Piano"},
			{Instrument: "Guitar"},
		},
	}
	fmt.Printf("Found group with ID %s\n", id)
	return group, nil
}

func (m *mock) UpdateSheetGroup(id string, sheet Sheet, title string) error {
	fmt.Printf("Updating group %s with sheet for instrument %s\n", title, sheet.Instrument)
	return nil
}
