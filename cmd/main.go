package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"t-meeting-backend/internal/adapters/postgres"
	"t-meeting-backend/internal/config"
	"t-meeting-backend/internal/repository"
	"t-meeting-backend/internal/route"
	"t-meeting-backend/internal/service"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()

	db, err := postgres.NewDB(ctx, cfg.DBDSN())
	if err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	eventRepo := repository.NewPgxEventRepository(db.Pool())
	svc := service.NewEventService(eventRepo)

	route.Setup(r, svc)

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	fmt.Printf("Listening on port %s...\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
