package route

import (
	"t-meeting-backend/internal/controller"
	"t-meeting-backend/internal/middleware"
	"t-meeting-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewEventRouter(router chi.Router, svc service.EventService, jwtSvc *service.JWTService) {
	ec := &controller.EventController{
		Svc: svc,
	}
	router.Route("", func(r chi.Router) {
		r.Use(middleware.JWTVerifier(jwtSvc))
		r.Post("/event", ec.Create)
		r.Get("/events", ec.GetAll)
		r.Route("/event/{eventID}", func(r chi.Router) {
			r.Get("/", ec.GetEventById)
			r.Put("/", ec.Update)
			r.Delete("/", ec.Delete)
		})
	})
}
