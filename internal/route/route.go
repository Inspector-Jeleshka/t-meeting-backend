package route

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(timeout time.Duration, db *pgxpool.Pool, r *chi.Mux) {
	r.Use(middleware.StripSlashes)
	NewEventRouter(timeout, db, r)
}
