package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Event struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Location    string `json:"location"`
}

func main() {
	events := []Event{}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	r.Post("/event", func(w http.ResponseWriter, r *http.Request) {
		var event Event

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&event); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		events = append(events, event)
		res, _ := json.Marshal(event)
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	r.Route("/event/{eventId}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			eventId := r.PathValue("eventId")
			i := slices.IndexFunc(events, func(e Event) bool {
				return e.Id == eventId
			})
			if i < 0 {
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}

			res, err := json.Marshal(events[i])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		})
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			eventId := r.PathValue("eventId")
			i := slices.IndexFunc(events, func(e Event) bool {
				return e.Id == eventId
			})

			var event Event
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&event); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if i < 0 {
				events = append(events, event)
			} else {
				events[i] = event
			}

			res, _ := json.Marshal(event)
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			eventId := r.PathValue("eventId")
			i := slices.IndexFunc(events, func(e Event) bool {
				return e.Id == eventId
			})
			if i < 0 {
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}

			events = slices.Delete(events, i, i+1)
		})
	})

	fmt.Printf("Listening on port %d...\n", 33)
	http.ListenAndServe(":33", r)
}
