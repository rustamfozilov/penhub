package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/jackc/pgx/v4"
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

func (s *Service) GetTokenForUser(ctx context.Context, user *types.User) (token string, err error) {
	ok,id, err := s.db.ValidateLoginAndPassword(ctx, user.Login, user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoSuchUser
	}
	if !ok {
		return "", ErrInvalidPassword
	}
	token, err = s.MakeToken(token)
	if err != nil {
		return "", err
	}
	err = s.db.PutNewToken(ctx, token, id)
	if err != nil {
		log.Println(err)
		return "", ErrInternal
	}
	return token, nil
}

func (s *Service) MakeToken(token string) (string, error) {
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) {
		return "", ErrInternal
	}
	if err != nil {
		log.Println(err)
		return "", ErrInternal
	}
	token = hex.EncodeToString(buffer)
	return token, nil
}
