package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errResp := ErrorResponse{
				Msg: err.Error(),
			}
			b, _ := json.Marshal(errResp)
			w.Write(b)
		}

		log.Info().Dict("request_body", zerolog.Dict().Bytes("request_body", body)).Send()
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		h.ServeHTTP(w, r)
	})
}
