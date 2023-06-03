package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackskj/carta"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type (
	User struct {
		ID         ulid.ULID `json:"id"`
		Name       string    `json:"name"`
		ProfileURL string    `json:"profile_url"`
	}

	Post struct {
		ID     ulid.ULID `json:"id"`
		Title  string    `json:"title"`
		Body   string    `json:"body"`
		UserID ulid.ULID `json:"user_id"`
	}

	PostWithUser struct {
		ID    ulid.ULID `json:"id"`
		Title string    `json:"title"`
		Body  string    `json:"body"`
		User  User      `json:"user"`
	}

	MySQLClient struct {
		Dbx *sqlx.DB
	}

	MySQLPost struct {
		ID     string `db:"id"`
		Title  string `db:"title"`
		Body   string `db:"body"`
		UserID string `db:"user_id"`
	}

	MySQLPostWithUser struct {
		ID    string `db:"post_id"`
		Title string `db:"post_title"`
		Body  string `db:"post_body"`
		User  struct {
			ID         string `db:"user_id"`
			Name       string `db:"user_name"`
			ProfileURL string `db:"user_profile_url"`
		}
	}
)

func NewUser(id ulid.ULID, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
}

func (u *User) ProfileImageURL() string {
	url := strings.Builder{}
	url.WriteString("https://storage.cloud.google.com/be-yamaguchi-kohei.appspot.com/profile_images/")
	url.WriteString(u.ProfileURL)
	return url.String()
}

func NewPost(id ulid.ULID, title string, body string, userID ulid.ULID) *Post {
	return &Post{
		ID:     id,
		Title:  title,
		Body:   body,
		UserID: userID,
	}
}

func NewMySQLClient(dbx *sqlx.DB) *MySQLClient {
	return &MySQLClient{
		Dbx: dbx,
	}
}

func (p *Post) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	return string(b)
}

func (p *Post) Key() string {
	return fmt.Sprintf("posts:%s", p.ID)
}

func (mp *MySQLPost) ToModel() *Post {
	return NewPost(ulid.MustParse(mp.ID), mp.Title, mp.Body, ulid.MustParse(mp.UserID))
}

func (mp *MySQLPostWithUser) ToModel() *PostWithUser {
	return &PostWithUser{
		ID:    ulid.MustParse(mp.ID),
		Title: mp.Title,
		Body:  mp.Body,
		User: User{
			ID:         ulid.MustParse(mp.User.ID),
			Name:       mp.User.Name,
			ProfileURL: mp.User.ProfileURL,
		},
	}
}

func (mysql *MySQLClient) InsertUser(ctx context.Context, name string, profileURL string) (*User, error) {
	query := `
		INSERT INTO users (id, name, profile_url) VALUES (?, ?, ?);
	`
	id := ulid.Make()
	_, err := mysql.Dbx.ExecContext(ctx, query, id.String(), name, profileURL)
	if err != nil {
		return nil, fmt.Errorf("models.InsertUser: failed to %s: %w", query, err)
	}

	return &User{ID: id, Name: name}, nil
}

func (mysql *MySQLClient) InsertPost(ctx context.Context, title string, body string, userID ulid.ULID) (*Post, error) {
	query := `
		INSERT INTO posts (id, title, body, user_id) VALUES (?, ?, ?, ?); 
	`
	id := ulid.Make()
	_, err := mysql.Dbx.ExecContext(ctx, query, id.String(), title, body, userID.String())
	if err != nil {
		return nil, fmt.Errorf("models.InsertPost: failed to %s: %w", query, err)
	}

	return &Post{ID: id, Title: title, Body: body, UserID: userID}, nil
}

func (mysql *MySQLClient) SelectPostsByUserID(ctx context.Context, userID ulid.ULID, limit *int) ([]*Post, error) {
	query := strings.Builder{}
	query.WriteString(`
	  select posts.id, posts.title, posts.body, posts.user_id 
		from posts 
		where posts.user_id = ?
	`)
	query.WriteString(" order by posts.id desc")
	if limit == nil || *limit == 0 {
		limit = ToPtr(10)
	}
	query.WriteString(fmt.Sprintf(" limit %d", *limit))

	log.Debug().Msgf(query.String())
	mPosts := []*MySQLPost{}
	err := mysql.Dbx.SelectContext(ctx, &mPosts, query.String(), userID.String())
	if err != nil {
		return nil, fmt.Errorf("models.SelectPostsByUserID: failed to %s: %w", query.String(), err)
	}
	posts := make([]*Post, len(mPosts))
	for i, mpost := range mPosts {
		posts[i] = mpost.ToModel()
	}

	return posts, nil
}

func (mysql *MySQLClient) ListPosts(ctx context.Context, limit *int) ([]*PostWithUser, error) {
	query := strings.Builder{}
	if limit != nil || *limit == 0 {
		limit = ToPtr(10)
	}
	query.WriteString(`
	  select P.id as post_id, P.title as post_title, P.body as post_body, U.id as user_id, U.name as user_name, U.profile_url as user_profile_url
		from posts as P
		join users as U
		on U.id = P.user_id
	`)
	query.WriteString(" order by P.id desc")
	query.WriteString(fmt.Sprintf(" limit %d", *limit))

	mPosts := []*MySQLPostWithUser{}
	rows, err := mysql.Dbx.QueryContext(ctx, query.String())
	if err != nil {
		err = fmt.Errorf("models.ListPosts: failed to %s: %w", query.String(), err)
		log.Debug().Msgf(err.Error())
		return nil, fmt.Errorf(err.Error())
	}

	err = carta.Map(rows, &mPosts)
	if err != nil {
		err = fmt.Errorf("models.ListPosts: failed to carta.Map: %w", err)
		log.Debug().Msgf(err.Error())
		return nil, fmt.Errorf(err.Error())
	}
	posts := make([]*PostWithUser, len(mPosts))
	for i, p := range mPosts {
		posts[i] = p.ToModel()
	}

	return posts, nil
}

func ToPtr[T any](v T) *T {
	return &v
}
