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

	err := json.NewDecoder(r.Body).Decode(chapter)
	if err != nil {
		log.Println("1")
		badRequest(w, err)
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, err)
		return
	}

	access, err := h.Service.BookAccess(r.Context(), userId, chapter.BookId)
	if !access {
		log.Println("no access")
		badRequest(w, err)
		return
	}

	err = h.Service.WriteChapter(r.Context(), chapter)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) GetBooksById(w http.ResponseWriter, r *http.Request)  {
	id, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, err)
		return
	}

	books, err := h.Service.GetBooksById(r.Context(), id)
	if err != nil {
		if err != nil {
			InternalServerError(w, err)
			return
		}
	}
	data, err := json.Marshal(books)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}