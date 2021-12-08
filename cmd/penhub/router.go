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

	authMux.Route("/books", func(r chi.Router) {

	})
	authMux.Post("/books/create", h.CreateBook)
	authMux.Get("/books/genres", h.GetAllGenres)
	authMux.Get("/books/genres/genre", h.GetGenreById)
	authMux.Post("/books/write", h.WriteBook)
	authMux.Get("/books", h.GetBooksByUserId) // my books
	authMux.Get("/chapters/list", h.GetChaptersByBookId)
	authMux.Get("/books/read", h.ReadChapter)
	authMux.Put("/books/title/edit", h.EditTitle)
	authMux.Put("/books/access/edit", h.EditAccess)
	authMux.Put("/chapters/content/edit", h.EditContent)
	authMux.Put("/chapters/name/edit", h.EditChapterName)
	authMux.Get("/search/title", h.SearchByTitle)
	authMux.Get("/search/author", h.SearchAuthor)
	authMux.Get("/books/author", h.GetBooksByAuthorId)
	authMux.Get("/search/genre", h.SearchGenre)
	authMux.Get("/books/genre", h.GetBooksByGenreId)
	authMux.Put("/books/image/edit", h.EditImage)
	authMux.Get("/books/image", h.GetImageByName)

	mux := chi.NewMux()
	mux.Mount(`/api/unauth`, unAuthMux)
	mux.Mount(`/api`, authMux)

	authMux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("he he"))
	}) // TODO delete
	return mux
}
