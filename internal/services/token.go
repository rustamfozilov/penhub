package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/jackc/pgx/v4"
	errors "github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/types"
	"log"
	"time"
)

func (s *Service) GetTokenForUser(ctx context.Context, user *types.User) (token string, err error) {
	log.Println(user.Login, user.Password)
	ok, id, err := s.db.ValidateLoginAndPassword(ctx, user.Login, user.Password)
	log.Println(ok, id)
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

func (s *Service) IdByToken(cxt context.Context, token string) (id int64, err error) {
	id, expire, err := s.db.IdByToken(cxt, token)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.Wrap(err, "no authorization")
	}
	if err != nil {
		return 0, err
	}
	if time.Now().After(expire) { //TODO дает фолс когда должен тру - исправить
		return 0, ErrExpired
	}
	return id, nil
}
