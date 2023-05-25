package main

import (
	"net/http"
	"time"

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

	db := NewMemDB()
	h := NewHandlers(db)

	http.HandleFunc("/posts", h.Posts)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Msgf("exit app-engine-go: %s", err.Error())
	}
}
