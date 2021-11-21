package handlers

import (
	"encoding/json"
	"github.com/rustamfozilov/penhub/internal/services"
	"github.com/rustamfozilov/penhub/internal/types"
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
