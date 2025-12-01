package route

import (
	"t-meeting-backend/internal/controller"
	"t-meeting-backend/internal/repository"
	"t-meeting-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewEventRouter(router chi.Router) {
	er := repository.NewEventRepository()
	ec := &controller.EventController{
		Svc: service.NewEventService(er),
	}
	router.Post("/event", ec.Create)
	router.Get("/events", ec.GetAll)
	router.Route("/event/{eventID}", func(router chi.Router) {
		router.Get("/", ec.GetEventById)
		router.Put("/", ec.Update)
		router.Delete("/", ec.Delete)
	})
}
