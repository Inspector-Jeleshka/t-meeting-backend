package route

import (
	"t-meeting-backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup(r *chi.Mux, svc service.EventService, jwtSvc *service.JWTService) {
	r.Use(middleware.StripSlashes)
	NewEventRouter(r, svc, jwtSvc)
}
