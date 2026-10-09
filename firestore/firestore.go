package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
)

type FireDb struct {
	client *firestore.Client
}

func CountQuery() AggregationQuery {
	return func(cr *firestore.CollectionRef) *firestore.AggregationQuery {
		return cr.NewAggregationQuery().WithCount("count")
	}
}

func MaxQuery(field string) DocQuery {
	return func(collection *firestore.CollectionRef) firestore.Query {
		return collection.OrderBy(field, firestore.Desc).Limit(1)
	}
}

func GroupKeyQuery(groupKey string) DocQuery {
	return func(collection *firestore.CollectionRef) firestore.Query {
		return collection.Where("group_key", "==", groupKey)
	}
}

func GroupKeysQuery(groupKeys []string) DocQuery {
	return func(collection *firestore.CollectionRef) firestore.Query {
		return collection.Where("group_key", "in", groupKeys)
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

func LastUpdatedQuery() DocQuery {
	return func(collection *firestore.CollectionRef) firestore.Query {
		return MaxQuery("last_updated")(collection)
	}
}

func PaginationQuery(startAfterID string, limit int) DocQuery {
	validatedLimit := normalizePaginationLimit(limit)
	return func(cr *firestore.CollectionRef) firestore.Query {
		q := cr.OrderBy(firestore.DocumentID, firestore.Asc).Limit(validatedLimit)
		if startAfterID != "" {
			q = q.StartAfter(startAfterID)
		}
		return q
	}
}

func normalizePaginationLimit(limit int) int {
	if limit < minPaginationLimit {
		return defaultPaginationLimit
	}
	if limit > maxPaginationLimit {
		return maxPaginationLimit
	}
	return limit
}

type DocUpdateFn func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error

type CreateDocFn func(ref *firestore.DocumentRef) (*firestore.WriteResult, error)

type DocQuery func(*firestore.CollectionRef) firestore.Query

type AggregationQuery func(*firestore.CollectionRef) *firestore.AggregationQuery

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

func (s *FireDb) SearchWithAggregate(collectionName string, queryFn AggregationQuery) (firestore.AggregationResult, error) {
	collection := s.client.Collection(collectionName)
	query := queryFn(collection)
	result, err := query.Get(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *FireDb) CountDocuments(collectionName string) (int64, error) {
	result, err := s.SearchWithAggregate(collectionName, CountQuery())
	if err != nil {
		return 0, err
	}
	count, ok := result["count"].(*firestorepb.Value)
	if !ok {
		return 0, fmt.Errorf("failed to parse count from aggregation result")
	}
	return count.GetIntegerValue(), nil
}

func (s *FireDb) SearchSingleDocument(collectionName string, queryFn DocQuery) (*firestore.DocumentSnapshot, error) {
	docs, err := s.SearchByQuery(collectionName, queryFn)
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("no documents found")
	}
	return docs[0], nil
}

func ParseDocument[T any](doc *firestore.DocumentSnapshot) (*T, error) {
	var result T
	if err := doc.DataTo(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func ParseDocuments[T any](docs []*firestore.DocumentSnapshot) ([]T, error) {
	result := make([]T, 0, len(docs))
	for _, doc := range docs {
		parsed, err := ParseDocument[T](doc)
		if err != nil {
			return nil, err
		}
		result = append(result, *parsed)
	}
	return result, nil
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
