package firestore

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/TheLuQ/eChart-backend/instrument"
)

type Sheet struct {
	instrument.Instrument
	Id string `firestore:"id"`
}

type SheetInfo struct {
	instrument.Instrument
	ParentPath string
	FullPath   string
	Title      string
}

func ToSheet(rawPath string) (*SheetInfo, error) {
	fileName := path.Base(rawPath)
	fileName = strings.TrimSuffix(fileName, path.Ext(fileName))

	parentName := path.Base(path.Dir(rawPath))
	cleanParentName := strings.ReplaceAll(parentName, "-", " ")
	if fileName == "" || cleanParentName == "." {
		return nil, fmt.Errorf("Invalid path format: %s", rawPath)
	}
	instrument, err := instrument.ParseFileName(fileName)
	if err != nil {
		return nil, fmt.Errorf("Invalid instrument name format: %s", fileName)
	}
	return &SheetInfo{Instrument: instrument, ParentPath: path.Dir(rawPath), FullPath: rawPath, Title: cleanParentName}, nil
}

type SheetGroup struct {
	Title       string  `firestore:"title,omitempty"`
	LastUpdated string  `firestore:"last_updated"`
	Sheets      []Sheet `firestore:"sheets"`
}

func NewSheetGroup(title string, sheet Sheet) *SheetGroup {
	return &SheetGroup{
		Title:       title,
		LastUpdated: time.Now().Local().String(),
		Sheets:      []Sheet{sheet},
	}
}

type InterfaceFirestoreDB interface {
	SaveSheet(sheet Sheet) error
	SearchGroupById(id string) (*SheetGroup, error)
	AddSheetToGroup(id string, sheet Sheet, title string) error
	RemoveSheetFromGroup(id string, sheet Sheet, title string) error
}

type SheetDb struct {
	client              *firestore.Client
	sheetCollectionName string
}

func (s *SheetDb) SaveSheet(sheet Sheet) error {
	result, err := s.client.Collection(s.sheetCollectionName).NewDoc().Create(context.Background(), sheet)
	if err != nil {
		return err
	}
	fmt.Printf("Sheet saved with result: %v\n", result)
	return nil
}

func (s *SheetDb) SearchGroupById(id string) (*SheetGroup, error) {
	doc, err := s.client.Collection(s.sheetCollectionName).Doc(id).Get(context.Background())
	if err != nil {
		return nil, err
	}
	var group SheetGroup
	err = doc.DataTo(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *SheetDb) AddSheetToGroup(id string, sheet Sheet, title string) error {
	return s.UpdateSheetGroup(id, sheet, title, firestore.ArrayUnion(sheet))
}

func (s *SheetDb) RemoveSheetFromGroup(id string, sheet Sheet, title string) error {
	return s.UpdateSheetGroup(id, sheet, title, firestore.ArrayRemove(sheet))
}

func (s *SheetDb) UpdateSheetGroup(id string, sheet Sheet, title string, arrayResult interface{}) error {
	_, err := s.client.Collection(s.sheetCollectionName).Doc(id).
		Set(context.Background(), map[string]interface{}{"instruments": arrayResult, "title": title}, firestore.MergeAll)
	return err
}

func New(dbName string, sheetCollectionName string) (*SheetDb, error) {
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		return nil, fmt.Errorf("PROJECT_ID environment variable is not set")
	}

	client, err := firestore.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, fmt.Errorf("Error creating Firestore client: %v", err)
	}
	return &SheetDb{
		client:              client,
		sheetCollectionName: sheetCollectionName,
	}, nil
}
