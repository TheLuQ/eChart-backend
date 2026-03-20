package firestore

import (
	"path"
	"time"

	"cloud.google.com/go/firestore"
)

type Instrument struct {
	Name    string `firestore:"instrument_name_en" json:"instrument_name_en"`
	NamePol string `firestore:"instrument_name_pol" json:"instrument_name_pol"`
	Voice   string `firestore:"voice,omitempty" json:"voice,omitempty"`
	Key     string `firestore:"key,omitempty" json:"key,omitempty"`
}

type Sheet struct {
	Instrument
	Id       string `firestore:"file_id" json:"file_id"`
	Title    string `firestore:"title,omitempty" json:"title"`
	FileName string `firestore:"file_name,omitempty" json:"file_name,omitempty"`
}

type SheetGroup struct {
	Title       string  `firestore:"title,omitempty"`
	LastUpdated string  `firestore:"last_updated"`
	ParentPath  string  `firestore:"parent_path,omitempty"`
	Sheets      []Sheet `firestore:"sheets"`
}

func (sg *SheetGroup) AddSheetIfMissing(sheet *Sheet) {
	if sg == nil {
		return
	}
	for i := range sg.Sheets {
		if sg.Sheets[i].Id == sheet.Id {
			sg.Sheets[i] = *sheet
			return
		}
	}
	sg.Sheets = append(sg.Sheets, *sheet)
}

func AddSheetFn(sheet *Sheet, parentPath string) DocUpdateFn {
	return func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error {
		err := transaction.Set(ref, map[string]interface{}{"sheets": firestore.ArrayUnion(sheet), "parent_path": parentPath,
			"last_updated": time.Now().UTC().String()}, firestore.MergeAll)
		return err
	}
}

func UpsertSheetFn(sheet *Sheet) DocUpdateFn {
	return func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error {
		var sg SheetGroup
		if doc, err := transaction.Get(ref); err == nil {
			if err := doc.DataTo(&sg); err != nil {
				return err
			}
		}
		sg.AddSheetIfMissing(sheet)
		sg.LastUpdated = time.Now().UTC().String()
		return transaction.Set(ref, sg)
	}
}

func RemoveSheetFn(sheet *Sheet) DocUpdateFn {
	return func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error {
		err := transaction.Set(ref, map[string]interface{}{"sheets": firestore.ArrayRemove(sheet),
			"last_updated": time.Now().UTC().String()}, firestore.MergeAll)
		return err
	}
}

func (sg *SheetGroup) GetSheets() []Sheet {
	for i := range sg.Sheets {
		sg.Sheets[i].Title = sg.Title
	}
	return sg.Sheets
}

func ToSheet(rawPath string, metadata map[string]string) (*Sheet, error) {
	fileName := path.Base(rawPath)
	instrumentName := metadata["instrument_name_en"]
	instrumentNamePol := metadata["instrument_name_pol"]
	voice := metadata["voice"]
	key := metadata["key"]

	instrument := Instrument{
		Name: instrumentName, NamePol: instrumentNamePol, Key: key, Voice: voice}
	return &Sheet{Instrument: instrument, Id: rawPath, FileName: fileName}, nil
}

func NewSheetGroup(title string, sheet Sheet) *SheetGroup {
	return &SheetGroup{
		Title:       title,
		LastUpdated: time.Now().UTC().String(),
		Sheets:      []Sheet{sheet},
	}
}
