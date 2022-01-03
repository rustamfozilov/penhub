package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/types"
	"net/http"
)

func (h *Handler) AddLike(w http.ResponseWriter, r *http.Request) {
	userID, err := GetIdFromContext(r.Context())
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	var bookId types.BookId
	err = json.NewDecoder(r.Body).Decode(&bookId)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}

	err = h.Service.AddLike(r.Context(), &userID, &bookId)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	var ratingId types.RatingID
	err := json.NewDecoder(r.Body).Decode(&ratingId)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}

	err = h.Service.DeleteLike(r.Context(), &ratingId)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) GetLikeId(w http.ResponseWriter, r *http.Request) {
	userID, err := GetIdFromContext(r.Context())
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	var bookId types.BookId
	err = json.NewDecoder(r.Body).Decode(&bookId)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	likeId, err := h.Service.GetLikeId(r.Context(), &userID, &bookId)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	FormatAndSending(w, likeId)
}

func (h *Handler) BookLikes(w http.ResponseWriter, r *http.Request) {
	var bookId types.BookId
	err := json.NewDecoder(r.Body).Decode(&bookId)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	likes, err := h.Service.BookLikes(r.Context(), &bookId)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	FormatAndSending(w, likes)
}
