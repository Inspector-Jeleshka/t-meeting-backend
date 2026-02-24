package route

import (
	"t-meeting-backend/internal/controller"
	"t-meeting-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewEventRouter(router chi.Router, svc service.EventService) {
	ec := &controller.EventController{
		Svc: svc,
	}
	router.Post("/event", ec.Create)
	router.Get("/events", ec.GetAll)
	router.Route("/event/{eventID}", func(r chi.Router) {
		r.Get("/", ec.GetEventById)
		r.Put("/", ec.Update)
		r.Delete("/", ec.Delete)
	})
}
