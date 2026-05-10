package route

import (
	"t-meeting-backend/internal/api/controller"
	"t-meeting-backend/internal/api/middleware"
	"t-meeting-backend/internal/jwt"

	"github.com/go-chi/chi/v5"
)

func NewEventRouter(router chi.Router, eventController controller.EventController, jwtManager *jwt.JWTManager) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.JWTVerifier(jwtManager))
		r.Post("/event", eventController.Create)
		r.Get("/events", eventController.GetAll)
		r.Route("/event/{eventID}", func(r chi.Router) {
			r.Get("/", eventController.GetEventById)
			r.Put("/", eventController.Update)
			r.Delete("/", eventController.Delete)
		})
	})

	router.Get("/published-event/{eventID}", eventController.GetPublishedEventByID)
}
