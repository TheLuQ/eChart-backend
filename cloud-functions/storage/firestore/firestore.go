package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
)

type FireDb struct {
	client              *firestore.Client
	sheetCollectionName string
}

func ParentPathQuery(title string) DocQuery {
	return func(collection *firestore.CollectionRef) firestore.Query {
		return collection.Where("parent_path", "==", title)
	}
}

type DocUpdateFn func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error

type DocQuery func(*firestore.CollectionRef) firestore.Query

func (s *FireDb) UpdateDocumentWithParentPath(path string, updateFn DocUpdateFn) error {
	return s.UpdateDocument(updateFn, ParentPathQuery(path))
}

func (s *FireDb) UpdateDocument(updateFn DocUpdateFn, query DocQuery) error {
	docQuery := query(s.client.Collection(s.sheetCollectionName))
	return s.client.RunTransaction(context.Background(), func(ctx context.Context, t *firestore.Transaction) error {
		documents := t.Documents(docQuery)
		docs, err := documents.GetAll()

		if err != nil {
			return err
		}

		if len(docs) == 0 {
			newRef := s.client.Collection(s.sheetCollectionName).NewDoc()
			return updateFn(t, newRef)
		}

		for _, d := range docs {
			err := updateFn(t, d.Ref)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *FireDb) SearchGroupById(id string) (*SheetGroup, error) {
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

func New(dbName string, sheetCollectionName string) (*FireDb, error) {
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		return nil, fmt.Errorf("PROJECT_ID environment variable is not set")
	}

	client, err := firestore.NewClientWithDatabase(context.Background(), projectId, dbName)
	if err != nil {
		return nil, fmt.Errorf("Error creating Firestore client: %v", err)
	}
	return &FireDb{
		client:              client,
		sheetCollectionName: sheetCollectionName,
	}, nil
}
