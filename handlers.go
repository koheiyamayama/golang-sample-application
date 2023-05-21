package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type (
	Handlers struct {
		db *MemDB
	}

	Post struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}
)

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

func FromStrToPost(str string) *Post {
	v := &Post{}
	if err := json.Unmarshal([]byte(str), v); err != nil {
		return &Post{}
	}

	return v
}

func NewHandlers(db *MemDB) *Handlers {
	return &Handlers{
		db: db,
	}
}

func (h *Handlers) Posts(w http.ResponseWriter, r *http.Request) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Debug().Msg("sample zerolog")

	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		items := h.db.List()
		var responses []*Post
		for _, item := range items {
			responses = append(responses, FromStrToPost(item))
		}
		b, _ := json.Marshal(responses)
		w.Write(b)
	case http.MethodPost:
		b := r.Body
		defer r.Body.Close()
		v := &Post{}
		if err := json.NewDecoder(b).Decode(v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v.ID = ulid.Make().String()
		if err := h.db.Create(v.ID, v.String()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(v.String()))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
