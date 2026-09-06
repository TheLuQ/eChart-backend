package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheLuQ/eChart-backend/firestore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type SheetQueryRequest struct {
	GroupKeys []string `json:"group_keys"`
}

type EventQueryRequest struct {
	IDs []string `json:"ids"`
}

func NewRouter(sheets firestore.SheetStore, events firestore.EventStore) http.Handler {
	r := chi.NewRouter()

	r.Use(corsMiddleware)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am working!"))
	})

	r.Get("/titles", func(w http.ResponseWriter, r *http.Request) {
		titles, err := sheets.GetAllTitles()
		if err != nil {
			http.Error(w, "Failed to retrieve titles: "+err.Error(), http.StatusInternalServerError)
			return
		}
		titleJson, err := json.Marshal(titles)
		if err != nil {
			http.Error(w, "Failed to marshal titles: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(titleJson)
	})

	r.Route("/api/sheets", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			options := ParseSheetOptions(r.URL.Query())
			result, err := sheets.SearchByGroupKeys(options.GroupKeys)
			if err != nil {
				http.Error(w, "Failed to retrieve sheets: "+err.Error(), http.StatusInternalServerError)
				return
			}
			sheetJson, err := json.Marshal(result)
			if err != nil {
				http.Error(w, "Failed to marshal sheets: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(sheetJson)
		})
		r.Get("/{groupKey}", func(w http.ResponseWriter, r *http.Request) {
			groupKey := chi.URLParam(r, "groupKey")
			result, err := sheets.SearchByGroupKeys([]string{groupKey})
			if err != nil {
				http.Error(w, "Failed to retrieve sheet: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if len(result) == 0 {
				http.Error(w, "Sheet not found", http.StatusNotFound)
				return
			}
			sheetJson, err := json.Marshal(result)
			if err != nil {
				http.Error(w, "Failed to marshal sheet: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(sheetJson)
		})
		r.Post("/query", func(w http.ResponseWriter, r *http.Request) {
			var query SheetQueryRequest
			if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
				http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
				return
			}
			result, err := sheets.SearchByGroupKeys(query.GroupKeys)
			if err != nil {
				http.Error(w, "Failed to search sheets: "+err.Error(), http.StatusInternalServerError)
				return
			}
			sheetJson, err := json.Marshal(result)
			if err != nil {
				http.Error(w, "Failed to marshal sheets: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(sheetJson)
		})
	})

	r.Route("/api/events", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			options, err := ParseEventOptions(r.URL.Query())
			if err != nil {
				http.Error(w, "Invalid query params: "+err.Error(), http.StatusBadRequest)
				return
			}

			var result []firestore.Event
			if len(options.IDs) > 0 {
				result, err = events.GetEventDetails(options.IDs)
			} else {
				result, err = events.GetShortEventsWithPagination(options.StartAfterID, options.Limit)
			}
			if err != nil {
				http.Error(w, "Failed to retrieve events: "+err.Error(), http.StatusInternalServerError)
				return
			}
			eventJson, err := json.Marshal(result)
			if err != nil {
				http.Error(w, "Failed to marshal events: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(eventJson)
		})
		r.Post("/query", func(w http.ResponseWriter, r *http.Request) {
			var query EventQueryRequest
			if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
				http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
				return
			}
			result, err := events.GetEventDetails(query.IDs)
			if err != nil {
				http.Error(w, "Failed to retrieve events: "+err.Error(), http.StatusInternalServerError)
				return
			}
			eventJson, err := json.Marshal(result)
			if err != nil {
				http.Error(w, "Failed to marshal events: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(eventJson)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var event firestore.Event
			if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
				http.Error(w, "Invalid event data: "+err.Error(), http.StatusBadRequest)
				return
			}
			if err := events.AddEvent(&event); err != nil {
				http.Error(w, "Failed to add event: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			events, err := events.GetEventDetails([]string{id})
			if err != nil {
				http.Error(w, "Failed to retrieve event details: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if events == nil {
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}
			eventJson, err := json.Marshal(events[0])
			if err != nil {
				http.Error(w, "Failed to marshal event details: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(eventJson)
		})
	})

	return r
}
