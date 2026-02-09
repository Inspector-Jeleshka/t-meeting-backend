package controller

import (
	"encoding/json"
	"net/http"
	"t-meeting-backend/internal/domain"
	"t-meeting-backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type EventController struct {
	Svc service.EventService
}

func (ec *EventController) Create(w http.ResponseWriter, r *http.Request) {
	var event domain.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "decode body: "+err.Error(), http.StatusBadRequest)
		//http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = ec.Svc.Create(r.Context(), &event)
	if err != nil {
		http.Error(w, "create event: "+err.Error(), http.StatusBadRequest)
		//http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ec *EventController) GetAll(w http.ResponseWriter, r *http.Request) {
	events, err := ec.Svc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if events == nil {
		events = make([]*domain.Event, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *EventController) GetEventById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "eventID")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	e, err := c.Svc.GetByID(ctx, eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if e == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ec *EventController) Update(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(r.PathValue("eventID"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var event domain.Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ec.Svc.Update(r.Context(), eventID, &event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ec *EventController) Delete(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(r.PathValue("eventID"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = ec.Svc.Delete(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
