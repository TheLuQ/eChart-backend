package firestore

import (
	"fmt"
	"testing"

	"github.com/TheLuQ/eChart-backend/instrument"
)

func TestSaveSheet(t *testing.T) {
	var sdb2 InterfaceFirestoreDB = &mock{}
	sheet := Sheet{Instrument: instrument.Instrument{Name: "Piano"}}
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
		if sheet.Instrument.Name != "guitar" {
			t.Errorf("Expected instrument 'guitar' but got %s", sheet.Instrument)
		}
		if sheet.Title != "my song" {
			t.Errorf("Expected title 'my song' but got %s", sheet.Title)
		}
		if sheet.ParentPath != "music/sheets/my-song" {
			t.Errorf("Expected parent path 'music/sheets/my-song' but got %s", sheet.ParentPath)
		}
		if sheet.FullPath != "music/sheets/my-song/guitar.pdf" {
			t.Errorf("Expected full path 'music/sheets/my-song/guitar.pdf' but got %s", sheet.FullPath)
		}
	})

	t.Run("with valid path without instrument", func(t *testing.T) {
		fullPath := "music/sheets/another-song/clarinet.pdf"
		sheet, err := ToSheet(fullPath)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if sheet.Instrument.Name != "clarinet" {
			t.Errorf("Expected instrument 'clarinet' but got %s", sheet.Instrument)
		}
		if sheet.Title != "another song" {
			t.Errorf("Expected title 'another song' but got %s", sheet.Title)
		}
		if sheet.ParentPath != "music/sheets/another-song" {
			t.Errorf("Expected parent path 'music/sheets/another-song' but got %s", sheet.ParentPath)
		}
		if sheet.FullPath != "music/sheets/another-song/clarinet.pdf" {
			t.Errorf("Expected full path 'music/sheets/another-song/clarinet.pdf' but got %s", sheet.FullPath)
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

type mock struct{}

func (m *mock) SaveSheet(sheet Sheet) error {
	fmt.Printf("Saving sheet for instrument %s\n", sheet.Instrument)
	return nil
}

func (m *mock) SearchGroupById(id string) (*SheetGroup, error) {
	sheet := Sheet{Instrument: instrument.Instrument{Name: "Piano"}}
	group := &SheetGroup{
		Title:       "Sample Group",
		LastUpdated: "2024-01-01",
		Sheets: []Sheet{
			sheet,
		},
	}
	fmt.Printf("Found group with ID %s\n", id)
	return group, nil
}

func (m *mock) AddSheetToGroup(id string, sheet Sheet, title string) error {
	fmt.Printf("Updating group %s with adding sheet for instrument %s\n", title, sheet.Instrument)
	return nil
}

func (m *mock) RemoveSheetFromGroup(id string, sheet Sheet, title string) error {
	fmt.Printf("Updating group %s with removing sheet for instrument %s\n", title, sheet.Instrument)
	return nil
}
