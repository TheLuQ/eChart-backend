package firestore

import (
	"context"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/TheLuQ/eChart-backend/sheet"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Sheet struct {
	sheet.Instrument
	Id string `firestore:"id"`
}

type SheetInfo struct {
	sheet.Instrument
	Id    string
	Title string
}

func ToSheet(rawPath string) (*SheetInfo, error) {
	fileName := path.Base(rawPath)
	fileName = strings.TrimSuffix(fileName, path.Ext(fileName))

	parentName := path.Base(path.Dir(rawPath))
	cleanParentName := strings.ReplaceAll(parentName, "-", " ")
	if fileName == "" || cleanParentName == "." {
		return nil, fmt.Errorf("Invalid path format: %s", rawPath)
	}
	instrument, err := sheet.ParseFileName(fileName)
	if err != nil {
		return nil, fmt.Errorf("Invalid instrument name format: %s", fileName)
	}
	return &SheetInfo{Instrument: instrument, Id: path.Dir(rawPath), Title: cleanParentName}, nil
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

func (sh *SheetGroup) AddSheet(sheet Sheet) {
	if slices.Contains(sh.Sheets, sheet) {
		return
	}
	sh.Sheets = append(sh.Sheets, sheet)
	sh.LastUpdated = time.Now().Local().String()
}

func (sh *SheetGroup) UpdateTitle(title string) {
	if title == "" || sh.Title == title {
		return
	}
	sh.Title = title
	sh.LastUpdated = time.Now().Local().String()
}

type InterfaceFirestoreDB interface {
	SaveSheet(sheet Sheet) error
	SearchGroupById(id string) (*SheetGroup, error)
	UpdateSheetGroup(id string, sheet Sheet, title string) error
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

func (s *SheetDb) UpdateSheetGroup(id string, sheet Sheet, title string) error {
	docRef := s.client.Collection(s.sheetCollectionName).Doc(id)
	return s.client.RunTransaction(context.Background(), func(ctx context.Context, t *firestore.Transaction) error {
		doc, err := t.Get(docRef)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return t.Set(docRef, NewSheetGroup(title, sheet))
			}
			return err
		}

		var group SheetGroup
		if err = doc.DataTo(&group); err != nil {
			return err
		}
		group.AddSheet(sheet)
		group.UpdateTitle(title)
		return t.Set(docRef, group)
	})
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
