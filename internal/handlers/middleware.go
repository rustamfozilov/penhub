package handlers

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rustamfozilov/penhub/internal/services"
	"log"
	"net/http"
)

type IDFunc func(ctx context.Context, token string) (id int64, err error)

type contextKey struct {
	key string
}

var AuthenticateContextKey = &contextKey{key: "authentication key"}

func Authentication(idFunc IDFunc) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			id, err := idFunc(r.Context(), token)
			if errors.Is(err, services.ErrExpired) || errors.Is(err, services.ErrNoAuthorization) {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

func NotFound(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		tctx := chi.NewRouteContext()
		if !rctx.Routes.Match(tctx, r.Method, r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
