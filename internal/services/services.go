package services

import (
	"context"
	"errors"
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

 var ErrInternal = errors.New("internal error")
var ErrLoginUsed = errors.New("login already registered")


func (s *Service) CreateBook(ctx context.Context, book *types.Book) error {

	err := s.db.CreateBook(ctx, book)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}

	return nil
}

func (s *Service) RegistrationUser(ctx context.Context, user *types.User) error {
	if s.db.IsLoginUsed(ctx, user.Login){
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
