package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
	var dbx *sqlx.DB
	var err error
	count := 0
	for {
		time.Sleep(5 * time.Second)
		dbx, err = sqlx.Open("mysql", "root:root@tcp(rds:3306)/google-cloud-go")
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
	mysqlClient := NewMySQLClient(dbx)

	res := os.Getenv("INSERT_SEED_DATA")
	if res == "" {
		err = InsertSeedData(ctx, mysqlClient)
		if err != nil {
			log.Warn().Msg(err.Error())
		} else {
			os.Setenv("INSERT_SEED_DATA", "1")
		}
	}

	db := NewMemDB()
	h := NewHandlers(db, mysqlClient)

	http.HandleFunc("/posts", h.Posts)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Msgf("exit google-cloud-go: %s", err.Error())
		os.Exit(1)
	}

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
