package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/services"
	"github.com/rustamfozilov/penhub/internal/types"
	"net/http"
)

func (h *Handler) SearchByTitle(w http.ResponseWriter, r *http.Request) {
	var SearchTitle types.BookTitle

	err := json.NewDecoder(r.Body).Decode(&SearchTitle)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	items, err := h.Service.SearchByTitle(r.Context(), &SearchTitle)
	if errors.Is(err, services.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return //TODO что отправить тогда? и вобще надо ли если нет такой ошибки
	}
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) SearchAuthor(w http.ResponseWriter, r *http.Request) {
	var SearchAuthor types.AuthorName

	err := json.NewDecoder(r.Body).Decode(&SearchAuthor)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	items, err := h.Service.SearchByAuthor(r.Context(), &SearchAuthor)
	if errors.Is(err, services.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return //TODO что отправить тогда? и вобще надо ли если нет такой ошибки
	}
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) SearchGenre(w http.ResponseWriter, r *http.Request) {

	var genreName types.GenreName

	err := json.NewDecoder(r.Body).Decode(&genreName)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}

	genres, err := h.Service.SearchGenre(r.Context(), genreName)
	data, err := json.Marshal(genres)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}

}
