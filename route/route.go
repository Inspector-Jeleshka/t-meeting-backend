package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup(r *chi.Mux) {
	r.Use(middleware.StripSlashes)
	NewEventRouter(r)
}
