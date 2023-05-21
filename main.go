package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
)

func main() {
	println("start app-engine-go")
	srv := &http.Server{
		Addr:         "0.0.0.0:8888",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	ctx := context.Background()

	projectID := GetGCPProjectID()

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	logName := "app-engine-go"
	logger := client.Logger(logName).StandardLogger(logging.Info)

	db := NewMemDB(logger)
	h := NewHandlers(db)

	http.HandleFunc("/posts", h.Posts)

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalf("exit app-engine-go: %s", err.Error())
	}
}
