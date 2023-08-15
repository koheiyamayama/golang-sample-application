package main

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/avast/retry-go"
	"github.com/jmoiron/sqlx"
	"github.com/koheiyamayama/ks-laboratory-backend/config"
	"github.com/koheiyamayama/ks-laboratory-backend/models"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	var dbx *sqlx.DB
	retry.Do(func() error {
		var err error
		dbx, err = sqlx.Open("mysql", config.ConnectDBInfo())
		log.Err(err).Msgf("failed to connect mysql: %s", config.ConnectDBInfo())
		return err
	}, retry.Attempts(5), retry.Delay(3*time.Second))

	mysqlClient := models.NewMySQLClient(dbx)
	errors := []error{}
	user1, e := mysqlClient.InsertUser(ctx, "kohei")
	if e != nil {
		errors = append(errors, e)
	}
	_, e = mysqlClient.InsertUser(ctx, "yamayama")
	if e != nil {
		errors = append(errors, e)
	}

	_, e = mysqlClient.InsertPost(ctx, "title1", "body1", user1.ID)
	if e != nil {
		errors = append(errors, e)
	}
	_, e = mysqlClient.InsertPost(ctx, "title2", "body2", user1.ID)
	if e != nil {
		errors = append(errors, e)
	}

	if len(errors) > 0 {
		for _, e := range errors {
			log.Err(e).Msgf("failed to insert: %s", e.Error())
		}
	}
}
