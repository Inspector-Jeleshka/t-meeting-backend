package app

import (
	"context"
	"fmt"
	"net/http"
	"t-meeting-backend/internal/api/controller"
	"t-meeting-backend/internal/api/route"
	"time"

	"t-meeting-backend/internal/adapters/postgres"
	"t-meeting-backend/internal/config"
	"t-meeting-backend/internal/repository"
	"t-meeting-backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type App struct {
	cfg    *config.Config
	db     *postgres.DB
	router *chi.Mux
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	// db setup
	db, err := postgres.NewDB(ctx, cfg.DBDSN())
	if err != nil {
		return nil, err
	}

	// events setup
	eventRepo := repository.NewPgxEventRepository(db.Pool())
	eventSvc := service.NewEventService(eventRepo)
	eventController := controller.EventController{Svc: eventSvc}

	jwtSvc := service.NewJWTService(
		"dev-secret-change-me",
		15*time.Minute,
		7*24*time.Hour,
	)

	// users/auth setup
	userRepo, err := repository.NewUserRepository(db.Pool())
	if err != nil {
		return nil, fmt.Errorf("new user repository: %w", err)
	}
	userSvc := service.NewUserService(userRepo)
	authController := controller.NewAuthController(userSvc, jwtSvc)

	// router setup
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	route.NewEventRouter(r, eventController, jwtSvc)
	route.SetupAuthRoutes(r, authController, jwtSvc)

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
