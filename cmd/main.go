package main

import (
	"context"
	"log"

	"t-meeting-backend/internal/app"
	"t-meeting-backend/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application, err := app.New(context.Background(), cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	defer application.Close()

	if err = application.Run(); err != nil {
		log.Fatal(err)
	}
}

//
//func init() {
//	cfg := config.MustLoad()
//	_, err := app.New(context.Background(), cfg)
//	if err != nil {
//	}
//}
