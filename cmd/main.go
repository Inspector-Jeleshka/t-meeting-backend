package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"t-meeting-backend/internal/route"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	r := chi.NewRouter()
	connString := "postgres://postgres:coolpassword@localhost:5433/tmeeting?sslmode=disable"
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	timeout := 5 * time.Second

	route.Setup(timeout, db, r)

	port := 33
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
