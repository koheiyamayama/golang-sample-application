package main

import (
	"context"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker"
	"github.com/jmoiron/sqlx"
	"github.com/koheiyamayama/google-cloud-go/config"
	"github.com/koheiyamayama/google-cloud-go/models"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	var dbx *sqlx.DB
	var err error
	count := 0
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
	InsertSeedData(ctx, mysqlClient)
}

func InsertSeedData(ctx context.Context, client *models.MySQLClient) error {
	log.Info().Msg("start insert seed data")

	faker := faker.New()
	p := faker.Person()
	users := make([]*models.User, 100)
	for i, _ := range users {
		users[i] = &models.User{
			ID:   ulid.Make(),
			Name: p.Name(),
		}
	}

	lorem := faker.Lorem()

	// TODO: Batch Insertなメソッドを定義する
	tx := client.Dbx.MustBegin()
	for _, user := range users {
		user, uErr := client.InsertUser(ctx, user.Name)
		if uErr != nil {
			log.Error().Stack().Err(uErr).Send()
			return uErr
		}

		posts := make([]*models.Post, 30)
		for i, _ := range posts {
			posts[i] = &models.Post{
				ID:     ulid.Make(),
				Title:  lorem.Text(rand.Intn(30)),
				Body:   lorem.Text(rand.Intn(280)),
				UserID: user.ID,
			}
		}

		for _, post := range posts {
			_, pErr := client.InsertPost(ctx, post.Title, post.Body, user.ID)
			if pErr != nil {
				log.Error().Stack().Err(pErr).Send()
				return pErr
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Info().Msg("end insert seed data")

	return nil
}
