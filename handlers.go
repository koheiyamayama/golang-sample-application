package main

import (
	"encoding/json"
	"net/http"

	"github.com/koheiyamayama/google-cloud-go/models"
)

type (
	Handlers struct {
		db          *MemDB
		mysqlClient *models.MySQLClient
	}

	InternalErrorResponse struct {
		Msg string `json:"msg"`
	}
)

func NewHandlers(db *MemDB, mysqlClient *models.MySQLClient) *Handlers {
	return &Handlers{
		db:          db,
		mysqlClient: mysqlClient,
	}
}

func (h *Handlers) ListPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	posts, err := h.mysqlClient.ListPosts(ctx, models.ToPtr(10))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errRes := &InternalErrorResponse{
			Msg: err.Error(),
		}
		b, _ := json.Marshal(errRes)
		w.Write(b)
		return
	}

	type postsResponse struct {
		Posts []*models.Post `json:"posts"`
	}
	postsRes := &postsResponse{
		Posts: posts,
	}
	b, err := json.Marshal(postsRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errRes := &InternalErrorResponse{
			Msg: err.Error(),
		}
		b, _ := json.Marshal(errRes)
		w.Write(b)
		return
	}
	w.Write(b)
}

func (h *Handlers) CreatePosts(w http.ResponseWriter, r *http.Request) {
	b := r.Body
	defer r.Body.Close()
	v := &models.Post{}
	if err := json.NewDecoder(b).Decode(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	retVal, err := h.mysqlClient.InsertPost(ctx, v.Title, v.Body, v.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errRes := &InternalErrorResponse{
			Msg: err.Error(),
		}
		b, _ := json.Marshal(errRes)
		w.Write(b)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(retVal.String()))
}
