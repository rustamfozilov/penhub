package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rustamfozilov/penhub/internal/types"
	"log"
	"time"
)

type DB struct {
	Pool *pgxpool.Pool
}


func NewDB() (*DB, error) {
	dsn := "postgres://app:pass@localhost:5432/penhub_db"

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
			log.Println(err)
			return nil, err
	}
	return &DB{Pool: pool}, nil
}

func (d *DB) CreateBook(ctx context.Context, book *types.Book) error {
	_, err := d.Pool.Exec(ctx, `
	insert into books (title, author_id, description, cover_image, access_read, active, created)
	 values ( $1, 1, $2, $3, $4, default, default)   
`, book.Title, book.Description, book.Image, book.AccessRead) //TODO author_id взять из id пользвателя
	if err != nil {
		return err
	}
	return nil
}
