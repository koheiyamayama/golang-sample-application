package main

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type (
	User struct {
		ID   ulid.ULID `json:"id"`
		Name string    `json:"name"`
	}

	Post struct {
		ID     ulid.ULID `json:"id" db:"id"`
		Title  string    `json:"title" db:"title"`
		Body   string    `json:"body" db:"body"`
		UserID ulid.ULID `json:"user_id" db:"user_id"`
	}

	MySQLClient struct {
		dbx *sqlx.DB
	}
)

func NewUser(id ulid.ULID, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
}

func NewPost(id ulid.ULID, title string, body string) *Post {
	return &Post{
		ID:    id,
		Title: title,
		Body:  body,
	}
}

func NewMySQLClient(dbx *sqlx.DB) *MySQLClient {
	return &MySQLClient{
		dbx: dbx,
	}
}

func (mysql *MySQLClient) InsertUser(ctx context.Context, name string) (*User, error) {
	query := `
		INSERT INTO users (id, name) VALUES (?, ?); 
	`
	id := ulid.Make()
	_, err := mysql.dbx.ExecContext(ctx, query, id.String(), name)
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
	_, err := mysql.dbx.ExecContext(ctx, query, id.String(), title, body, userID.String())
	if err != nil {
		return nil, fmt.Errorf("models.InsertPost: failed to %s: %w", query, err)
	}

	return &Post{ID: id, Title: title, Body: body}, nil
}

func (mysql *MySQLClient) SelectPostsByUserID(ctx context.Context, userID ulid.ULID) ([]*Post, error) {
	query := `
		SELECT * FROM posts WHERE posts.user_id = ?; 
	`
	posts := []*Post{}
	err := mysql.dbx.SelectContext(ctx, &posts, query, userID)
	if err != nil {
		return nil, fmt.Errorf("models.SelectPostsByUserID: failed to %s: %w", query, err)
	}

	return posts, nil
}
