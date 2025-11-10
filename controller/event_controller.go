package controller

import (
	"encoding/json"
	"net/http"
	"t-meeting-backend/domain"
	"t-meeting-backend/usecase"

	"github.com/google/uuid"
)

type EventController struct {
	EventUsecase usecase.EventUsecase
}

func (ec *EventController) Create(w http.ResponseWriter, r *http.Request) {
	var event domain.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = ec.EventUsecase.Create(r.Context(), &event)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
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

func (ec *EventController) GetAll(w http.ResponseWriter, r *http.Request) {
	events, err := ec.EventUsecase.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(events)
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

func (ec *EventController) GetByID(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(r.PathValue("eventID"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	event, err := ec.EventUsecase.GetByID(r.Context(), eventID)
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

	err = ec.EventUsecase.Update(r.Context(), eventID, &event)
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

	err = ec.EventUsecase.Delete(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
