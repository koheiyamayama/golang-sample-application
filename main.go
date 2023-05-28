package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
	for {
		time.Sleep(5 * time.Second)
		dbx, err = sqlx.Open("mysql", ConnectDBInfo())
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
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      r,
	}

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

func ConnectDBInfo() string {
	type opt struct {
		Key   string
		Value string
	}
	type opts []*opt

	options := opts{
		{Key: "parseTime", Value: "true"},
	}
	info := strings.Builder{}
	base := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", GetDBUserName(), GetDBPassword(), GetDBHostName(), GetDBPort(), GetDatabaseName())
	info.WriteString(base)

	for _, option := range options {
		info.WriteString(fmt.Sprintf("?%s=%s", option.Key, option.Value))
	}

	return info.String()
}

func GetDBUserName() string {
	defaultName := "root"
	if userName := os.Getenv("DATABASE_USER_NAME"); userName != "" {
		return userName
	} else {
		return defaultName
	}
}

func GetDBPassword() string {
	defaultPassword := "root"
	if password := os.Getenv("DATABASE_PASSWORD"); password != "" {
		return password
	} else {
		return defaultPassword
	}
}

func GetDBHostName() string {
	defaultHostname := "rds"
	if hostname := os.Getenv("DATABASE_HOSTNAME"); hostname != "" {
		return hostname
	} else {
		return defaultHostname
	}
}

func GetDBPort() string {
	defaultPort := "3306"
	if port := os.Getenv("DATABASE_PORT"); port != "" {
		return port
	} else {
		return defaultPort
	}
}

func GetDatabaseName() string {
	defaultDatabaseName := "google-cloud-go"
	if databaseName := os.Getenv("DATABASE_NAME"); databaseName != "" {
		return databaseName
	} else {
		return defaultDatabaseName
	}
}
