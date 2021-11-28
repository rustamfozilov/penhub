package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/types"
	"log"
	"net/http"
)

func (h *Handler) EditTitle(w http.ResponseWriter, r *http.Request) {
	var edit types.Book

	err := json.NewDecoder(r.Body).Decode(&edit)
	if err != nil {
		err := errors.Wrap(err, "")
		badRequest(w, err)
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		err := errors.Wrap(err, "")
		InternalServerError(w, err)
		return
	}
	access, err := h.Service.BookAccess(r.Context(), userId, edit.ID)
	if err != nil {
		err := errors.Wrap(err, "")
		InternalServerError(w, err)
		return
	}
	if !access {
		log.Println("no access")
		badRequest(w, err)
		return
	}


	err = h.Service.EditTitle(r.Context(), &edit)
	if err != nil {
		err := errors.Wrap(err, "")
			InternalServerError(w, err)
		return
	}
}


//func (h *Handler) EditContent(w http.ResponseWriter, r *http.Request)  {
//var n
//		json.NewDecoder(r.Body).Decode(&)
//
//
//}