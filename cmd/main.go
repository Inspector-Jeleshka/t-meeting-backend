package main

import (
	"fmt"
	"log"
	"net/http"
	"t-meeting-backend/route"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	route.Setup(r)

	port := 33
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
