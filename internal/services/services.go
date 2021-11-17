package services

import (
	"context"
	"errors"
	"github.com/rustamfozilov/penhub/internal/db"
	"github.com/rustamfozilov/penhub/internal/types"
	"log"
)

type Service struct {
	db *db.DB
}

func NewService(db *db.DB) *Service {
	return &Service{db: db}
}

 var ErrInternal = errors.New("internal error")

func (s *Service) CreateBook(ctx context.Context, book *types.Book) error {

	err := s.db.CreateBook(ctx, book)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}

	return nil
}
