package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type (
	User struct {
		ID   ulid.ULID `json:"id"`
		Name string    `json:"name"`
	}

	Post struct {
		ID     ulid.ULID `json:"id"`
		Title  string    `json:"title"`
		Body   string    `json:"body"`
		UserID ulid.ULID `json:"user_id"`
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
)

func NewUser(id ulid.ULID, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
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

func (mp *MySQLPost) ToPost() *Post {
	return NewPost(ulid.MustParse(mp.ID), mp.Title, mp.Body, ulid.MustParse(mp.UserID))
}

func (mysql *MySQLClient) InsertUser(ctx context.Context, name string) (*User, error) {
	query := `
		INSERT INTO users (id, name) VALUES (?, ?); 
	`
	id := ulid.Make()
	_, err := mysql.Dbx.ExecContext(ctx, query, id.String(), name)
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

func (mysql *MySQLClient) SelectPostsByUserID(ctx context.Context, userID ulid.ULID) ([]*Post, error) {
	query := `
		SELECT * FROM posts WHERE posts.user_id = ?; 
	`
	posts := []*Post{}
	err := mysql.Dbx.SelectContext(ctx, &posts, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("models.SelectPostsByUserID: failed to %s: %w", query, err)
	}

	return posts, nil
}

func (mysql *MySQLClient) ListPosts(ctx context.Context, limit *int) ([]*Post, error) {
	query := strings.Builder{}
	if limit != nil || *limit == 0 {
		limit = ToPtr(10)
	}
	query.WriteString("select * from posts")
	query.WriteString(" order by posts.id desc")
	query.WriteString(fmt.Sprintf(" limit %d", *limit))

	log.Debug().Msgf("query: %s", query.String())
	mPosts := []*MySQLPost{}
	err := mysql.Dbx.SelectContext(ctx, &mPosts, query.String())
	if err != nil {
		return nil, fmt.Errorf("models.ListPosts: failed to %s: %w", query, err)
	}

	posts := make([]*Post, len(mPosts))
	for i, mpost := range mPosts {
		posts[i] = mpost.ToPost()
	}
	return posts, nil
}

func ToPtr[T any](v T) *T {
	return &v
}
