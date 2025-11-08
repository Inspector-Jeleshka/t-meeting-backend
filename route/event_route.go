package route

import (
	"t-meeting-backend/controller"
	"t-meeting-backend/repository"
	"t-meeting-backend/usecase"

	"github.com/go-chi/chi/v5"
)

func NewEventRouter(router chi.Router) {
	er := repository.NewEventRepository()
	ec := &controller.EventController{
		EventUsecase: usecase.NewEventUsecase(er),
	}
	router.Post("/event", ec.Create)
	router.Get("/events", ec.GetAll)
	router.Route("/event/{eventID}", func(router chi.Router) {
		router.Get("/", ec.GetByID)
		router.Put("/", ec.Update)
		router.Delete("/", ec.Delete)
	})
}
