package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/TheLuQ/eChart-backend/firestore"
)

type mockSheetStore struct{}

func (m *mockSheetStore) GetAllTitles() ([]firestore.SheetGroup, error) { return nil, nil }
func (m *mockSheetStore) SearchByGroupKeys(groupKeys []string) ([]firestore.Sheet, error) {
	return nil, nil
}
func (m *mockSheetStore) AddSheet(sheet *firestore.Sheet) error    { return nil }
func (m *mockSheetStore) UpsertSheet(sheet *firestore.Sheet) error { return nil }
func (m *mockSheetStore) RemoveSheet(sheet *firestore.Sheet) error { return nil }
func (m *mockSheetStore) GetEtag() (string, error)                 { return "", nil }

type mockEventStore struct {
	detailsCalls    int
	paginationCalls int
	gotIDs          []string
	gotStartAfterID string
	gotLimit        int
	detailsResult   []firestore.Event
	pageResult      []firestore.Event
}

func (m *mockEventStore) AddEvent(event *firestore.Event) error { return nil }

func (m *mockEventStore) GetShortEventsWithPagination(startAfterID string, limit int) ([]firestore.Event, error) {
	m.paginationCalls++
	m.gotStartAfterID = startAfterID
	m.gotLimit = limit
	return m.pageResult, nil
}

func (m *mockEventStore) GetEventDetails(ids []string) ([]firestore.Event, error) {
	m.detailsCalls++
	m.gotIDs = append([]string{}, ids...)
	return m.detailsResult, nil
}

func TestEventsGetUsesIDFilterWhenIDsProvided(t *testing.T) {
	eventsStore := &mockEventStore{detailsResult: []firestore.Event{{ID: "xyz1"}, {ID: "xyz2"}}}
	handler := NewRouter(&mockSheetStore{}, eventsStore)

	req := httptest.NewRequest(http.MethodGet, "/api/events/?id=xyz1&id=xyz2&startAfterId=ignored&limit=2", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %d", res.Code)
	}
	if eventsStore.detailsCalls != 1 {
		t.Fatalf("expected details call once but got %d", eventsStore.detailsCalls)
	}
	if eventsStore.paginationCalls != 0 {
		t.Fatalf("expected no pagination calls but got %d", eventsStore.paginationCalls)
	}
	if !reflect.DeepEqual(eventsStore.gotIDs, []string{"xyz1", "xyz2"}) {
		t.Fatalf("expected ids [xyz1 xyz2] but got %v", eventsStore.gotIDs)
	}
}

func TestEventsGetUsesPaginationWhenNoIDsProvided(t *testing.T) {
	eventsStore := &mockEventStore{pageResult: []firestore.Event{{ID: "abc"}}}
	handler := NewRouter(&mockSheetStore{}, eventsStore)

	req := httptest.NewRequest(http.MethodGet, "/api/events/?startAfterId=cursor-1&limit=50", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %d", res.Code)
	}
	if eventsStore.paginationCalls != 1 {
		t.Fatalf("expected pagination call once but got %d", eventsStore.paginationCalls)
	}
	if eventsStore.detailsCalls != 0 {
		t.Fatalf("expected no details calls but got %d", eventsStore.detailsCalls)
	}
	if eventsStore.gotStartAfterID != "cursor-1" {
		t.Fatalf("expected cursor cursor-1 but got %q", eventsStore.gotStartAfterID)
	}
	if eventsStore.gotLimit != 50 {
		t.Fatalf("expected limit 50 but got %d", eventsStore.gotLimit)
	}

	var body []firestore.Event
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid JSON body but got %v", err)
	}
	if len(body) != 1 || body[0].ID != "abc" {
		t.Fatalf("expected paginated event response but got %v", body)
	}
}

func TestEventsGetRejectsInvalidLimit(t *testing.T) {
	eventsStore := &mockEventStore{}
	handler := NewRouter(&mockSheetStore{}, eventsStore)

	req := httptest.NewRequest(http.MethodGet, "/api/events/?limit=0", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 but got %d", res.Code)
	}
	if eventsStore.paginationCalls != 0 || eventsStore.detailsCalls != 0 {
		t.Fatalf("expected no event store calls but got pagination=%d details=%d", eventsStore.paginationCalls, eventsStore.detailsCalls)
	}
}

func TestFirestoreSheetGroupToJson(t *testing.T) {
	t.Run("marshals to expected JSON", func(t *testing.T) {
		sg := firestore.SheetGroup{
			Title:       "title",
			Band:        "band",
			GroupKey:    "group_key",
			LastUpdated: time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC),
		}
		tr := TitleResponse{
			Title:       sg.Title,
			Band:        sg.Band,
			GroupKey:    sg.GroupKey,
			LastUpdated: sg.LastUpdated,
		}
		want := fmt.Sprintf(`{"title":"%s","band":"%s","group_key":"%s","last_updated":"%s"}`,
			sg.Title, sg.Band, sg.GroupKey, sg.LastUpdated.Format(time.RFC3339))
		got, err := json.Marshal(tr)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(got) != want {
			t.Fatalf("MarshalJSON() = %q, want %q", got, want)
		}
	})
}
