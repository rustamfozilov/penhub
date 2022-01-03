package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/types"
	"net/http"
)

func (h *Handler) SearchByTitle(w http.ResponseWriter, r *http.Request) {
	var searchTitle types.BookTitle
	err := json.NewDecoder(r.Body).Decode(&searchTitle)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	err = h.Service.ValidateTitle(searchTitle.Title)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	items, err := h.Service.SearchByTitle(r.Context(), &searchTitle)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	FormatAndSending(w, items)
}

func (h *Handler) SearchAuthor(w http.ResponseWriter, r *http.Request) {
	var searchAuthor types.AuthorName
	err := json.NewDecoder(r.Body).Decode(&searchAuthor)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	err = h.Service.ValidateTitle(searchAuthor.Name)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	items, err := h.Service.SearchByAuthor(r.Context(), &searchAuthor)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	FormatAndSending(w, items)
}

func (h *Handler) SearchGenre(w http.ResponseWriter, r *http.Request) {
	var genreName types.GenreName
	err := json.NewDecoder(r.Body).Decode(&genreName)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	err = h.Service.ValidateTitle(genreName.Name)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	genres, err := h.Service.SearchGenre(r.Context(), genreName)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	FormatAndSending(w, genres)

}
