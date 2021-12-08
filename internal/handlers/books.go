package handlers

import (
	"encoding/json"
	errors "github.com/pkg/errors"
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
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}
	var b types.Book

	data := r.FormValue("data")
	//log.Println("data", data)
	err = json.Unmarshal([]byte(data), &b)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	b.AuthorId = id

	file, header, err := r.FormFile("image")
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	filename := header.Filename
	log.Println(header)

	book, err := h.Service.SaveImage(file, filename, &b)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	err = h.Service.CreateBook(r.Context(), book)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) WriteBook(w http.ResponseWriter, r *http.Request) {
	chapter := &types.Chapter{}

	err := json.NewDecoder(r.Body).Decode(chapter)
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

	access, err := h.Service.BookAccess(r.Context(), userId, chapter.BookId)
	if !access {
		log.Println("no access")
		badRequest(w, err)
		return
	}
	if err != nil {
		InternalServerError(w, err)
		return
	}
	err = h.Service.WriteChapter(r.Context(), chapter)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}

func (h *Handler) GetBooksByUserId(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromContext(r.Context())
	if err != nil {
		err := errors.WithStack(err)
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

func (h *Handler) GetChaptersByBookId(w http.ResponseWriter, r *http.Request) {
	var BookIdReq types.BookId
	err := json.NewDecoder(r.Body).Decode(&BookIdReq)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	chapters, err := h.Service.GetChaptersByBookId(r.Context(), &BookIdReq)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data, err := json.Marshal(chapters)
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

func (h *Handler) ReadChapter(w http.ResponseWriter, r *http.Request) {
	var chapterId types.ChapterId
	err := json.NewDecoder(r.Body).Decode(&chapterId)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	chapter, err := h.Service.ReadChapter(r.Context(), &chapterId)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data, err := json.Marshal(chapter)
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

func (h *Handler) GetBooksByAuthorId(w http.ResponseWriter, r *http.Request) {
	var authorId types.AuthorId
	err := json.NewDecoder(r.Body).Decode(&authorId)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}
	books, err := h.Service.GetBooksById(r.Context(), authorId.Id)
	if err != nil {
		if err != nil {
			InternalServerError(w, err)
			return
		}
	}
	data, err := json.Marshal(books)
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

func (h *Handler) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.Service.GetAllGenres(r.Context())
	if err != nil {
		InternalServerError(w, err)
	}
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

func (h *Handler) GetBooksByGenreId(w http.ResponseWriter, r *http.Request) {
	var genreId types.GenreID
	err := json.NewDecoder(r.Body).Decode(&genreId)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}

	books, err := h.Service.GetBooksByGenreId(r.Context(), genreId)
	if err != nil {
		InternalServerError(w, err)
	}
	data, err := json.Marshal(books)
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

func (h *Handler) GetGenreById(w http.ResponseWriter, r *http.Request) {
	var genreId types.GenreID
	err := json.NewDecoder(r.Body).Decode(&genreId)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}

	genre, err := h.Service.GetGenreById(r.Context(), genreId)
	if err != nil {
		InternalServerError(w, err)
	}
	data, err := json.Marshal(genre)
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

func (h *Handler) GetImageByName(w http.ResponseWriter, r *http.Request) {
	var n types.ImageName
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		err := errors.WithStack(err)
		badRequest(w, err)
		return
	}

	file, err := h.Service.GetImageByName(n.Name)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "image/png") //TODO

	_, err = w.Write(file)
	if err != nil {
		err := errors.WithStack(err)
		InternalServerError(w, err)
		return
	}

}
