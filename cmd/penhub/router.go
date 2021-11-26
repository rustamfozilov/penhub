package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rustamfozilov/penhub/internal/handlers"
	"net/http"
)

func NewRouter(h *handlers.Handler) *chi.Mux {

	unAuthMux := chi.NewMux()
	unAuthMux.Route("/user", func(r chi.Router) {
		r.Post("/registration", h.RegistrationUser)
		r.Get("/token", h.GetTokenForUser)
	})

	authMux := chi.NewMux()
	authMux.Use(handlers.Authentication(h.Service.IdByToken))
	authMux.Post("/book", h.CreateBook)
	authMux.Post("/write", h.WriteBook)
	authMux.Get("/books", h.GetBooksById)

	mux := chi.NewMux()
	mux.Mount(`/api/unauth`, unAuthMux)
	mux.Mount(`/api/`, authMux)

	authMux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("he he"))
	})
	return mux
}
