package main

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Debug().Msg("start google-cloud-go")

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	ctx := context.Background()

	projectID := GetGCPProjectID()

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal().Msgf("failed to create cloud logging client")
	}
	defer client.Close()

	db := NewMemDB()
	h := NewHandlers(db)

	http.HandleFunc("/posts", h.Posts)

	if err := srv.ListenAndServe(); err != nil {
		// logger.Fatalf("exit app-engine-go: %s", err.Error())
		log.Fatal().Msgf("exit app-engine-go: %w", err)
	}
}
