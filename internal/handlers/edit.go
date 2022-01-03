package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/types"
	"net/http"
)

func (h *Handler) EditChapter(w http.ResponseWriter, r *http.Request) {
	var editChapter types.Chapter
	err := json.NewDecoder(r.Body).Decode(&editChapter)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, errors.WithStack(err))
		return
	}
	access, err := h.Service.HaveAccessToEditBook(r.Context(), userId, editChapter.BookId)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if !access {
		Forbidden(w, errors.New("no access"))
		return
	}

	if editChapter.Content != "" {
		err = h.Service.EditContent(r.Context(), &editChapter)
		if err != nil {
			InternalServerError(w, err)
			return
		}
	}

	if editChapter.Name != "" {
		err := h.Service.ValidateChapter(&editChapter)
		if err != nil {
			badRequest(w, errors.WithStack(err))
			return
		}

		err = h.Service.EditChapterName(r.Context(), &editChapter)
		if err != nil {
			InternalServerError(w, err)
			return
		}
	}
	if editChapter.Number != 0 {
		err = h.Service.EditChapterNumber(r.Context(), &editChapter)
		if err != nil {
			InternalServerError(w, err)
			return
		}
	}

}

func (h *Handler) EditImage(w http.ResponseWriter, r *http.Request) {
	var b types.Book
	data := r.FormValue("data")
	err := json.Unmarshal([]byte(data), &b)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, errors.WithStack(err))
		return
	}
	access, err := h.Service.HaveAccessToEditBook(r.Context(), userId, b.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if !access {
		Forbidden(w, errors.New("no access"))
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	filename := header.Filename
	err = h.Service.ValidateImage(header.Size)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
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

func (h *Handler) EditBook(w http.ResponseWriter, r *http.Request) {
	var edit types.Book
	err := json.NewDecoder(r.Body).Decode(&edit)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, errors.WithStack(err))
		return
	}
	access, err := h.Service.HaveAccessToEditBook(r.Context(), userId, edit.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if !access {
		Forbidden(w, errors.New("no access"))
		return
	}
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	err = h.Service.EditAccess(r.Context(), &edit) //permanent change active
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if edit.Title != "" {
		err = h.Service.ValidateTitle(edit.Title)
		if err != nil {
			badRequest(w, errors.WithStack(err))
			return
		}
		err = h.Service.EditTitle(r.Context(), &edit)
		if err != nil {
			InternalServerError(w, err)
			return
		}
	}
	if edit.Genre != 0 {
		err := h.Service.EditGenre(r.Context(), &edit)
		if err != nil {
			InternalServerError(w, err)
			return
		}
	}
	if edit.Description != "" {
		err = h.Service.ValidateDescription(edit.Description)
		if err != nil {
			badRequest(w, errors.WithStack(err))
			return
		}
		err = h.Service.EditDescription(r.Context(), &edit)
		if err != nil {
			InternalServerError(w, err)
			return
		}
	}
}

func (h Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	var edit types.Book
	err := json.NewDecoder(r.Body).Decode(&edit)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, errors.WithStack(err))
		return
	}
	access, err := h.Service.HaveAccessToEditBook(r.Context(), userId, edit.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if !access {
		Forbidden(w, errors.New("no access"))
		return
	}
	err = h.Service.DeleteBook(r.Context(), &edit)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	var editChapter types.Chapter
	err := json.NewDecoder(r.Body).Decode(&editChapter)
	if err != nil {
		badRequest(w, errors.WithStack(err))
		return
	}
	userId, err := GetIdFromContext(r.Context())
	if err != nil {
		InternalServerError(w, errors.WithStack(err))
		return
	}
	access, err := h.Service.HaveAccessToEditBook(r.Context(), userId, editChapter.BookId)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if !access {
		Forbidden(w, errors.New("no access"))
		return
	}
	err = h.Service.DeleteChapter(r.Context(), &editChapter)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}
