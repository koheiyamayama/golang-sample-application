package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/koheiyamayama/google-cloud-go/config"
	"github.com/koheiyamayama/google-cloud-go/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Debug().Msg("start google-cloud-go")

	ctx := context.Background()
	var dbx *sqlx.DB
	var err error
	count := 0
	log.Debug().Msg("connect mysql...")
	for {
		time.Sleep(5 * time.Second)
		dbx, err = sqlx.Open("mysql", config.ConnectDBInfo())
		if err != nil {
			log.Warn().Msgf("retry because of %s", err.Error())
		} else {
			break
		}

		if count >= 6 {
			log.Fatal().Msg("failed to open mysql")
		}

		count += 1
	}
	mysqlClient := models.NewMySQLClient(dbx)

	h := NewHandlers(mysqlClient)

	r := chi.NewRouter()
	r.Use(Logging)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", h.ListPosts)
			r.Post("/", h.CreatePosts)
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/posts", h.ListPostsByUser)
			})
		})
	})

	srv := &http.Server{
		Addr:         config.AppServerAddr(),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      r,
	}

	log.Debug().Msg("start app server")
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Msgf("exit google-cloud-go: %s", err.Error())
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
		log.Fatal().Msgf("Failed to gracefully shutdown:", err)
	}
	log.Info().Msg("Server shutdown")
}
