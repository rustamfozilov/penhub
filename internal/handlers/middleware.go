package handlers

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/rustamfozilov/penhub/internal/services"
	"net/http"
)

type IDFunc func(ctx context.Context, token string) (id int64, err error)

type contextKey struct {
	key string
}

var AuthenticateContextKey = &contextKey{key: "authentication key"}

//var ErrNoAuthentication = errors.New("no authentication")

func Identification(idFunc IDFunc) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			id, err := idFunc(r.Context(), token)
			if errors.Is(err, services.ErrExpired) || errors.Is(err, pgx.ErrNoRows) {
				badRequest(w, err)
			}
			if err != nil {
				InternalServerError(w, err)
				return
			}
			ctx := context.WithValue(r.Context(), AuthenticateContextKey, id)
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
		})
	}
}
