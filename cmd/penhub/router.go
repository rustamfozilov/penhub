package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rustamfozilov/penhub/internal/handlers"
)

func NewRouter(h *handlers.Handler) *chi.Mux {

	unAuthMux := chi.NewMux()
	unAuthMux.Route("/user", func(r chi.Router) {
		r.Post("/registration", h.RegistrationUser)
		r.Get("/token", h.GetTokenForUser)
	})

	authMux := chi.NewMux()
	authMux.Use(handlers.Authentication(h.Service.IdByToken))
	authMux.Post("/api/books", h.CreateBook)
	authMux.Post("/api/write", h.WriteBook)

	mux := chi.NewMux()
	mux.Mount(`/api/unauth`, unAuthMux)
	mux.Mount(`/api/`, authMux)

	return mux
}
