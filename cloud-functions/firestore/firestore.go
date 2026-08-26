package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
)

type FireDb struct {
	client *firestore.Client
}

func GroupKeyQuery(groupKey string) DocQuery {
	return func(collection *firestore.CollectionRef) firestore.Query {
		return collection.Where("group_key", "==", groupKey)
	}
}

func IdQuery(ids []string) DocQuery {
	return func(cr *firestore.CollectionRef) firestore.Query {
		var refs []*firestore.DocumentRef
		for _, id := range ids {
			refs = append(refs, cr.Doc(id))
		}
		return cr.Where(firestore.DocumentID, "in", refs)
	}
}

type DocUpdateFn func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error

type CreateDocFn func(ref *firestore.DocumentRef) (*firestore.WriteResult, error)

type DocQuery func(*firestore.CollectionRef) firestore.Query

func (s *FireDb) UpdateDocumentWithGroupKey(collectionName string, groupKey string, updateFn DocUpdateFn) error {
	return s.UpdateDocument(collectionName, updateFn, GroupKeyQuery(groupKey))
}

func (s *FireDb) UpdateDocument(collectionName string, updateFn DocUpdateFn, query DocQuery) error {
	docQuery := query(s.client.Collection(collectionName))
	return s.client.RunTransaction(context.Background(), func(ctx context.Context, t *firestore.Transaction) error {
		documents := t.Documents(docQuery)
		docs, err := documents.GetAll()

		if err != nil {
			return err
		}

		if len(docs) == 0 {
			newRef := s.client.Collection(collectionName).NewDoc()
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

func (s *FireDb) UpdateDocumentWithId(collectionName string, id string, updateFn CreateDocFn) error {
	_, error := updateFn(s.client.Collection(collectionName).Doc(id))
	if error != nil {
		return error
	}
	println("Document update result: " + id)
	return nil
}

func (s *FireDb) CreateDocument(collectionName string, createFn CreateDocFn) error {
	newRef := s.client.Collection(collectionName).NewDoc()
	if _, err := createFn(newRef); err != nil {
		return err
	}
	println("Document created succesfully with ID: " + newRef.ID)
	return nil
}

func (s *FireDb) SearchByQuery(collectionName string, query DocQuery) ([]*firestore.DocumentSnapshot, error) {
	documents := query(s.client.Collection(collectionName))
	docs, err := documents.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func New(dbName string) (*FireDb, error) {
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		return nil, fmt.Errorf("PROJECT_ID environment variable is not set")
	}

	client, err := firestore.NewClientWithDatabase(context.Background(), projectId, dbName)
	if err != nil {
		return nil, fmt.Errorf("Error creating Firestore client: %v", err)
	}
	return &FireDb{
		client: client,
	}, nil
}
