package main

import (
	"fmt"
	"log"
	"net/http"
	"t-meeting-backend/internal/repository"
	"t-meeting-backend/internal/service"

	"t-meeting-backend/internal/config"
	"t-meeting-backend/internal/route"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()
	r := chi.NewRouter()
	repo := repository.NewEventRepository(cfg.DBDSN())
	svc := service.NewEventService(repo)
	route.Setup(r, svc)

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	fmt.Printf("Listening on port %s...\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
