package handlers

import (
	"encoding/json"
	"github.com/rustamfozilov/penhub/internal/services"
	"github.com/rustamfozilov/penhub/internal/types"
	"log"
	"net/http"
)

type Handler struct {
	Service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {

	id, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, err)
		return
	}

	var b types.Book
	b.ID = id
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		badRequest(w, err)
		return
	}
	err = h.Service.CreateBook(r.Context(), &b)
	if err != nil {
		InternalServerError(w, err)
		return
	}

}

func (h *Handler) WriteBook(w http.ResponseWriter, r *http.Request) {
	chapter := &types.Chapter{}
	bookName := &types.BookTitle{}
	err := json.NewDecoder(r.Body).Decode(chapter)
	if err != nil {
		badRequest(w, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(bookName)
	if err != nil {
		badRequest(w, err)
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, err)
		return
	}
	bookId, err := h.Service.GetBookId(r.Context(), bookName)
	if err != nil {
		badRequest(w, err)
		return
	}
	access, err := h.Service.BookAccess(r.Context(), userId, bookId)
	if !access {
		log.Println("no access")
		badRequest(w, err)
		return
	}
	chapter.BookId = bookId
	err = h.Service.WriteChapter(r.Context(), chapter)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}
