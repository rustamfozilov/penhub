package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rustamfozilov/penhub/internal/handlers"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(handlers.NotFound)

	unAuthMux := chi.NewMux()
	unAuthMux.Route("/user", func(r chi.Router) {
		r.Post("/registration", h.RegistrationUser)
		r.Get("/token", h.GetTokenForUser)
	})
	authMux := chi.NewMux()
	authMux.Use(handlers.Authentication(h.Service.IdByToken))
	authMux.Route("/books", func(r chi.Router) {
		r.Post("/create", h.CreateBook)
		r.Get("/genres", h.GetAllGenres)
		r.Get("/genres/id", h.GetGenreById)
		r.Get("/", h.GetBooksByUserId) // my books
		r.Put("/edit", h.EditBook)     // нужно всегда отправлять access true чтобы депжать активной
		r.Put("/image/edit", h.EditImage)
		r.Get("/image", h.GetImageByName)
		r.Delete("/delete", h.DeleteBook) // also, for recover
	})
	authMux.Route("/chapters", func(r chi.Router) {
		r.Post("/write", h.WriteChapter)
		r.Get("/list", h.GetChaptersByBookId)
		r.Get("/read", h.ReadChapter)
		r.Put("/edit", h.EditChapter)
		r.Delete("/delete", h.DeleteChapter) //also, for recover
	})
	authMux.Route("/search", func(r chi.Router) {
		r.Get("/title", h.SearchByTitle)
		r.Get("/author", h.SearchAuthor)
		r.Get("/author/books", h.GetBooksByAuthorId)
		r.Get("/genre", h.SearchGenre)
		r.Get("/genre/books", h.GetBooksByGenreId)
	})
	authMux.Route("/rating", func(r chi.Router) {
		r.Post("/like", h.AddLike)
		r.Get("/like", h.GetLikeId)
		r.Delete("/like", h.DeleteLike)
		r.Get("/book", h.BookLikes)
	})
	mux.Mount(`/api/unauth`, unAuthMux)
	mux.Mount(`/api`, authMux)

	return mux
}
