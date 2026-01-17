package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"t-meeting-backend/internal/adapters/postgres"
	"t-meeting-backend/internal/config"
	"t-meeting-backend/internal/repository"
	"t-meeting-backend/internal/route"
	"t-meeting-backend/internal/service"
)

type App struct {
	cfg    *config.Config
	db     *postgres.DB
	router *chi.Mux
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	db, err := postgres.NewDB(ctx, cfg.DBDSN())
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	eventRepo := repository.NewPgxEventRepository(db.Pool())
	eventSvc := service.NewEventService(eventRepo)

	route.Setup(r, eventSvc)

	return &App{
		cfg:    cfg,
		db:     db,
		router: r,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.cfg.HTTPPort)
	fmt.Printf("Listening on %s...\n", addr)
	return http.ListenAndServe(addr, a.router)
}

func (a *App) Close() {
	a.db.Close()
}
