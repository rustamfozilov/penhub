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



func (s *Service) BookAccess(ctx context.Context, userId, bookId int64 ) (bool, error) {
	return  s.db.BookAccess(ctx, userId, bookId)
}

func (s *Service) WriteChapter(ctx context.Context, chapter *types.Chapter) error {
		return s.db.WriteChapter(ctx, chapter)
}

func (s *Service) GetBooksById(ctx context.Context, id int64) ([]*types.Book, error) {
	return s.db.GetBooksById(ctx, id)
}

func (s *Service) GetChaptersByBookId(ctx context.Context, bookId *types.BookId) ([]*types.Chapter, error){
		return s.db.GetChaptersByBookId(ctx, bookId.Id)
}

func (s *Service) ReadChapter(ctx context.Context, chaptId *types.ChapterId) (*types.Chapter,error) {
	return s.db.ReadChapter(ctx,chaptId.Id)
}

func (s *Service) EditTitle(ctx context.Context, edit *types.Book) error {

	return s.db.EditTitle(ctx, edit.ID, edit.Title)
}