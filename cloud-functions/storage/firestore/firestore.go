package firestore

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
)

type Sheet struct {
	Instrument string `firestore:"instrument"`
	Title      string `firestore:"title"`
}

func NewSheet(fullPath string) (Sheet, error) {
	dirs := strings.Split(fullPath, "/")
	if len(dirs) > 1 {
		instrument := dirs[len(dirs)-1]
		title := dirs[len(dirs)-2]
		sheet := Sheet{
			Instrument: instrument,
			Title:      title,
		}
		return sheet, nil
	}
	return Sheet{}, fmt.Errorf("invalid path: %s", fullPath)
}

type InterfaceFirestoreDB interface {
	SaveSheet(sheet Sheet) error
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
