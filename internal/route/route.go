package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"t-meeting-backend/internal/service"
)

func Setup(r *chi.Mux, svc service.EventService) {
	r.Use(middleware.StripSlashes)
	NewEventRouter(r, svc)
}
