package firestore

import (
	"fmt"
	"testing"
)

func TestSaveSheet(t *testing.T) {
	var sdb2 InterfaceFirestoreDB = &mock{}
	sheet := Sheet{
		Instrument: "Piano",
		Title:      "Another Song",
	}
	err := sdb2.SaveSheet(sheet)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestNewSheet(t *testing.T) {
	fullPath := "music/sheets/My Song/Guitar"
	sheet, err := NewSheet(fullPath)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if sheet.Instrument != "Guitar" {
		t.Errorf("Expected instrument 'Guitar' but got %s", sheet.Instrument)
	}
	if sheet.Title != "My Song" {
		t.Errorf("Expected title 'My Song' but got %s", sheet.Title)
	}
}

func TestNewSheetInvalidPath(t *testing.T) {
	fullPath := "invalidpath"
	_, err := NewSheet(fullPath)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

type mock struct{}

func (m *mock) SaveSheet(sheet Sheet) error {
	fmt.Printf("Saving sheet %s for instrument %s\n", sheet.Title, sheet.Instrument)
	return nil
}
