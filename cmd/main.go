package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/koheiyamayama/google-cloud-go/models"
	"github.com/rs/zerolog/log"
)

func main() {
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
	mysqlClient := models.NewMySQLClient(dbx)
	InsertSeedData(ctx, mysqlClient)
}

func InsertSeedData(ctx context.Context, client *models.MySQLClient) error {
	log.Debug().Msg("start insert seed data")
	fp, err := os.ReadFile("./testdata/posts.json")
	if err != nil {
		return err
	}

	posts := []*models.Post{}
	err = json.Unmarshal(fp, &posts)
	if err != nil {
		return err
	}

	fu, err := os.ReadFile("./testdata/posts.json")
	if err != nil {
		return err
	}

	users := []*models.User{}
	err = json.Unmarshal(fu, &users)
	if err != nil {
		return err
	}

	var uErrors []error
	var pErrors []error
	// TODO: Batch Insertなメソッドを定義する
	tx := client.Dbx.MustBegin()
	for _, user := range users {
		user, uErr := client.InsertUser(ctx, user.Name)
		if uErr != nil {
			uErrors = append(uErrors, uErr)
		}

		for _, post := range posts {
			_, pErr := client.InsertPost(ctx, post.Title, post.Body, user.ID)
			if pErr != nil {
				pErrors = append(pErrors, pErr)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	if len(pErrors) != 0 {
		return fmt.Errorf("failed to insert %d posts", len((pErrors)))
	} else if len(uErrors) != 0 {
		return fmt.Errorf("failed to insert %d posts", len((uErrors)))
	}

	log.Debug().Msg("end insert seed data")

	return nil
}
