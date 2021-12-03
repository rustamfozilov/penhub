package services

import (
	"context"
	errors "github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/db"
	"github.com/rustamfozilov/penhub/internal/types"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Service struct {
	db *db.DB
}

func NewService(db *db.DB) *Service {
	return &Service{db: db}
}

var ErrExpired = errors.New("token expired")
var ErrInternal = errors.New("internal error")
var ErrLoginUsed = errors.New("login already registered")
var ErrNoSuchUser = errors.New("no such user")
var ErrInvalidPassword = errors.New("invalid password")
var ErrNotFound = errors.New("not found")

func (s *Service) CreateBook(ctx context.Context, book *types.Book) error {
	err := s.db.CreateBook(ctx, book)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}

	return nil
}

func (s *Service) RegistrationUser(ctx context.Context, user *types.User) error {
	if s.db.IsLoginUsed(ctx, user.Login) {
		return ErrLoginUsed
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}

	err = s.db.RegistrationUser(ctx, user, hash)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}

	return nil
}

func (s *Service) GetBookId(ctx context.Context, bookName *types.BookTitle) (id int64, err error) {
	return s.db.GetBookId(ctx, bookName.Title)
}

func (s *Service) BookAccess(ctx context.Context, userId, bookId int64) (bool, error) {
	return s.db.BookAccess(ctx, userId, bookId)
}

func (s *Service) WriteChapter(ctx context.Context, chapter *types.Chapter) error {
	return s.db.WriteChapter(ctx, chapter)
}

func (s *Service) GetBooksById(ctx context.Context, id int64) ([]*types.Book, error) {
	return s.db.GetBooksById(ctx, id)
}

func (s *Service) GetChaptersByBookId(ctx context.Context, bookId *types.BookId) ([]*types.Chapter, error) {
	return s.db.GetChaptersByBookId(ctx, bookId.Id)
}

func (s *Service) ReadChapter(ctx context.Context, chapterId *types.ChapterId) (*types.Chapter, error) {
	chapter, err := s.db.ReadChapter(ctx, chapterId.Id)
	if err != nil {
		return nil, err
	}
	if chapter.Active == true {
		return chapter, nil
	}
	return nil, errors.New("haven't access to read")
}

func (s *Service) EditTitle(ctx context.Context, edit *types.Book) error {
	return s.db.EditTitle(ctx, edit.ID, edit.Title)
}

func (s *Service) EditContent(ctx context.Context, edit *types.Chapter) error {

	return s.db.EditContent(ctx, edit)

}

func (s *Service) EditAccess(ctx context.Context, edit *types.Book) error {
	return s.db.EditAccess(ctx, edit)
}

func (s *Service) EditChapterName(ctx context.Context, edit *types.Chapter) error {
	return s.db.EditChapterName(ctx, edit)
}

func (s *Service) SearchByTitle(ctx context.Context, title *types.BookTitle) ([]*types.Book, error) {
	books, err := s.db.SearchByTitle(ctx, title)
	if errors.Is(err, db.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *Service) SearchByAuthor(ctx context.Context, author *types.AuthorName) ([]*types.User, error) {
	return s.db.SearchByAuthor(ctx, author)
}

func (s *Service) GetAllGenres(ctx context.Context) ([]*types.Genre, error) {
	return s.db.GetAllGenres(ctx)
}

func (s *Service) SearchGenre(ctx context.Context, genreName types.GenreName) ([]*types.Genre, error) {
	return s.db.SearchGenre(ctx, genreName)
}

func (s *Service) GetBooksByGenreId(ctx context.Context, genreId types.GenreID) ([]*types.Book, error) {
	return s.db.GetBooksByGenreId(ctx, genreId)
}

func (s *Service) GetGenreById(ctx context.Context, genreId types.GenreID) (*types.Genre, error) {
	return s.db.GetGenreById(ctx, genreId)
}
