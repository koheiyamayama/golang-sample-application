package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avast/retry-go"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/koheiyamayama/ks-laboratory-backend/config"
	"github.com/koheiyamayama/ks-laboratory-backend/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Debug().Msg("start ks-laboratory-backend")

	ctx := context.Background()
	var dbx *sqlx.DB
	retry.Do(func() error {
		var err error
		dbx, err = sqlx.Open("mysql", config.ConnectDBInfo())
		log.Err(err).Msgf("failed to connect mysql: %s", config.ConnectDBInfo())
		return err
	}, retry.Attempts(5), retry.Delay(3*time.Second))

	mysqlClient := models.NewMySQLClient(dbx)

	h := NewHandlers(mysqlClient)

	r := chi.NewRouter()
	r.Use(Logging)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/healthd", func(r chi.Router) {
			r.Get("/", h.Health)
		})
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", h.ListPosts)
			r.Post("/", h.CreatePosts)
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", h.GetUserByID)
				r.Get("/posts", h.ListPostsByUser)
			})
		})
	})

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      r,
	}

	go func() {
		log.Debug().Msg("start ks-laboratory-backend process")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Msgf("exit ks-laboratory-backend: %s", err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Fatal().Msgf("Failed to gracefully shutdown: %s", err.Error())
	}
	log.Info().Msg("Server shutdown")
}
