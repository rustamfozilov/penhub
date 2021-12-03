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

func Authentication(idFunc IDFunc) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			//log.Println("token:", token)
			id, err := idFunc(r.Context(), token)
			if errors.Is(err, services.ErrExpired) || errors.Is(err, pgx.ErrNoRows) {
				badRequest(w, err)
				return
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
