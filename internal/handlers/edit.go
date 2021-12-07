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
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	access, err := h.Service.BookAccess(r.Context(), userId, edit.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if !access {
		badRequest(w, errors.New("no access"))
		return
	}

	err = h.Service.EditTitle(r.Context(), &edit)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) EditContent(w http.ResponseWriter, r *http.Request) {
	var editChapter types.Chapter
	err := json.NewDecoder(r.Body).Decode(&editChapter)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	access, err := h.Service.BookAccess(r.Context(), userId, editChapter.BookId)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	if !access {
		badRequest(w, errors.New("no access"))
		return
	}
	err = h.Service.EditContent(r.Context(), &editChapter)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) EditAccess(w http.ResponseWriter, r *http.Request) {
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
		InternalServerError(w, err)
		return
	}
	if !access {
		log.Println("no access")
		badRequest(w, err)
		return
	}
	err = h.Service.EditAccess(r.Context(), &edit)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) EditChapterName(w http.ResponseWriter, r *http.Request)  {
	var editChapter types.Chapter
	err := json.NewDecoder(r.Body).Decode(&editChapter)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	access, err := h.Service.BookAccess(r.Context(), userId, editChapter.BookId)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	if !access {
		badRequest(w, errors.New("no access"))
		return
	}
	err = h.Service.EditChapterName(r.Context(), &editChapter)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) EditImage(w http.ResponseWriter, r *http.Request) {
	var b types.Book
	data := r.FormValue("data")

	err := json.Unmarshal([]byte(data), &b)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	filename := header.Filename

	book, err := h.Service.SaveImage(file, filename, &b)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	err = h.Service.EditImage(r.Context(), book)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

//TODO edit chapter number, edit genre, edit description