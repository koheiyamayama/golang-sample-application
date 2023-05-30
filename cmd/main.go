package main

import (
	"context"
	"fmt"
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
	// CreateSeedSQL(ctx)
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
	profileImages := []string{
		"profile_1.jpg",
		"profile_2.jpg",
		"profile_3.jpg",
		"profile_4.jpg",
		"profile_5.jpg",
		"profile_6.jpg",
		"profile_7.jpg",
		"profile_8.jpg",
		"profile_9.jpg",
		"profile_10.jpg",
		"profile_11.jpg",
		"profile_12.jpg",
	}
	tx := client.Dbx.MustBegin()
	for _, user := range users {
		i := rand.Intn(11)
		user, uErr := client.InsertUser(ctx, user.Name, profileImages[i])
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

func CreateSeedSQL(ctx context.Context) error {
	log.Info().Msg("start insert seed data")

	profileImages := []string{
		"profile_1.jpg",
		"profile_2.jpg",
		"profile_3.jpg",
		"profile_4.jpg",
		"profile_5.jpg",
		"profile_6.jpg",
		"profile_7.jpg",
		"profile_8.jpg",
		"profile_9.jpg",
		"profile_10.jpg",
		"profile_11.jpg",
		"profile_12.jpg",
	}

	faker := faker.New()
	p := faker.Person()
	users := make([]*models.User, 100)
	for i, _ := range users {
		n := rand.Intn(11)
		users[i] = &models.User{
			ID:         ulid.Make(),
			Name:       p.Name(),
			ProfileURL: profileImages[n],
		}
	}

	lorem := faker.Lorem()

	for _, user := range users {
		fmt.Printf("INSERT INTO users (`id`, `name`, `profile_url`) VALUES ('%s', '%s', '%s');\n", user.ID.String(), user.Name, user.ProfileURL)
		posts := make([]*models.Post, 100)
		for i, _ := range posts {
			posts[i] = &models.Post{
				ID:     ulid.Make(),
				Title:  lorem.Text(rand.Intn(30)),
				Body:   lorem.Text(rand.Intn(280)),
				UserID: user.ID,
			}
		}

		for _, post := range posts {
			fmt.Printf("INSERT INTO posts (`id`, `title`, `body`, `user_id`) VALUES ('%s', '%s', '%s', '%s');\n", post.ID.String(), post.Title, post.Body, post.UserID.String())
		}
	}

	log.Info().Msg("end insert seed data")

	return nil
}
