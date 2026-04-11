package route

import (
	"t-meeting-backend/internal/api/controller"
	"t-meeting-backend/internal/api/middleware"
	"t-meeting-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewEventRouter(router chi.Router, eventController controller.EventController, jwtSvc *service.JWTService) {
	router.Route("", func(r chi.Router) {
		r.Use(middleware.JWTVerifier(jwtSvc))
		r.Post("/event", eventController.Create)
		r.Get("/events", eventController.GetAll)
		r.Route("/event/{eventID}", func(r chi.Router) {
			r.Get("/", eventController.GetEventById)
			r.Put("/", eventController.Update)
			r.Delete("/", eventController.Delete)
		})
	})
}
