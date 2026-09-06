package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

const (
	defaultPaginationLimit = 20
	minPaginationLimit     = 1
	maxPaginationLimit     = 100
)

type AgendaItem struct {
	Time        string `json:"time"        firestore:"time"`
	Title       string `json:"title"       firestore:"title"`
	Description string `json:"description,omitempty" firestore:"description,omitempty"`
}

type Event struct {
	ID          string       `json:"id"          firestore:"-"`
	Title       string       `json:"title"       firestore:"title"`
	Date        string       `json:"date"        firestore:"date"`
	Time        string       `json:"time"        firestore:"time,omitempty"`
	Description string       `json:"description" firestore:"description"`
	Location    string       `json:"location"    firestore:"location"`
	Agenda      []AgendaItem `json:"agenda"      firestore:"agenda,omitempty"`
	GroupKeys   []string     `json:"group_keys"  firestore:"group_keys,omitempty"`
	Tags        []string     `json:"tags"        firestore:"tags,omitempty"`
}

func CreateEvent(event *Event) CreateDocFn {
	return func(ref *firestore.DocumentRef) (*firestore.WriteResult, error) {
		return ref.Set(context.Background(), event)
	}
}

func ShortEventQuery() DocQuery {
	return func(cr *firestore.CollectionRef) firestore.Query {
		return cr.SelectPaths(
			firestore.FieldPath{"title"},
			firestore.FieldPath{"date"},
			firestore.FieldPath{"location"},
			firestore.FieldPath{"description"},
			firestore.FieldPath{"group_keys"},
			firestore.FieldPath{"id"},
		)
	}
}

type EventConnector struct {
	Db             *FireDb
	CollectionName string
}

func (c *EventConnector) AddEvent(event *Event) error {
	return c.Db.CreateDocument(c.CollectionName, CreateEvent(event))
}

func (c *EventConnector) GetShortEvents() ([]Event, error) {
	var events []Event
	docs, err := c.Db.SearchByQuery(c.CollectionName, ShortEventQuery())
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		var event Event
		if err := doc.DataTo(&event); err != nil {
			return nil, err
		}
		event.ID = doc.Ref.ID
		events = append(events, event)
	}
	return events, nil
}

func (c *EventConnector) GetShortEventsWithPagination(startAfterID string, limit int) ([]Event, error) {
	var events []Event
	docs, err := c.Db.SearchByQuery(c.CollectionName, PaginationQuery(startAfterID, limit))
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		var event Event
		if err := doc.DataTo(&event); err != nil {
			return nil, err
		}
		event.ID = doc.Ref.ID
		events = append(events, event)
	}
	return events, nil
}

func (c EventConnector) GetEventDetails(id []string) ([]Event, error) {
	docs, err := c.Db.SearchByQuery(c.CollectionName, IdQuery(id))
	if err != nil {
		return nil, err
	}
	var events []Event
	for _, doc := range docs {
		var event Event
		if err := doc.DataTo(&event); err != nil {
			return nil, err
		}
		event.ID = doc.Ref.ID
		events = append(events, event)
	}
	return events, nil
}
