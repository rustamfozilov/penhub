package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	errors "github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/types"
	"golang.org/x/crypto/bcrypt"
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
		return nil, errors.Wrap(err, "pgx fail connect")
	}
	return &DB{Pool: pool}, nil
}

func (d *DB) CreateBook(ctx context.Context, book *types.Book) error {
	_, err := d.Pool.Exec(ctx, `
	insert into books (title, author_id, description, cover_image, access_read, genre, active, created)
	 values ($1, $2, $3, $4, $5,$6, default, default)   
`, book.Title, book.ID, book.Description, book.Image, book.AccessRead, book.Genre)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (d *DB) RegistrationUser(ctx context.Context, user *types.User, hash []byte) error {

	_, err := d.Pool.Exec(ctx, `
		insert into users (name, login, password, active, created)
		values ($1, $2, $3, default, default)
`, user.Name, user.Login, hash)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) IsLoginUsed(ctx context.Context, login string) bool {

	rows, err := d.Pool.Query(ctx, `
		select login from users 
`)
	if errors.Is(err, pgx.ErrNoRows) {
		return false
	}

	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var loginDb string
		err := rows.Scan(&loginDb)
		if err != nil {
			return false
		}
		if loginDb == login {
			return true
		}
	}
	err = rows.Err()
	if err != nil {
		return false
	}

	return false
}

func (d *DB) ValidateLoginAndPassword(ctx context.Context, login, password string) (bool, int64, error) {
	var id int64
	var hash string
	err := d.Pool.QueryRow(ctx, `
		select password, id from users where login = $1
`, login).Scan(&hash, &id)
	if err != nil {
		return true, 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println("invalid password")
		return false, 0, err
	}
	return true, id, nil
}

func (d *DB) PutNewToken(ctx context.Context, token string, id int64) error {
	_, err := d.Pool.Exec(ctx, `
			insert into users_tokens (user_id, token, expire, created) 
			values ($1, $2, default, default)
`, id, token)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) IdByToken(cxt context.Context, token string) (id int64, expire time.Time, err error) {
	err = d.Pool.QueryRow(cxt, `
select user_id, expire from users_tokens where token = $1
`, token).Scan(&id, &expire)
	if err != nil {
		return 0, expire, err
	}
	return id, expire, nil
}

func (d *DB) GetBookId(ctx context.Context, title string) (id int64, err error) {

	err = d.Pool.QueryRow(ctx, `
		select id from books where title = $1
`, title).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("book not exist")
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *DB) BookAccess(ctx context.Context, userId, bookId int64) (bool, error) {
	var authorId int64
	err := d.Pool.QueryRow(ctx, `
		select author_id from books where id = $1
`, bookId).Scan(&authorId)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("book not exist")
	}
	if err != nil {
		return false, err
	}
	if authorId == userId {
		return true, nil
	}
	return false, nil
}

func (d *DB) WriteChapter(ctx context.Context, chapter *types.Chapter) error {
	_, err := d.Pool.Exec(ctx, `
		insert into chapters (book_id, number, name, content,active, created) 
		values ($1, $2, $3, $4, default, default)
`, chapter.BookId, chapter.Number, chapter.Name, chapter.Content)

	if err != nil {
		return err
	}
	return nil
}

func (d *DB) GetBooksById(ctx context.Context, id int64) ([]*types.Book, error) {

	books := make([]*types.Book, 0)

	rows, err := d.Pool.Query(ctx, `
		select id, title, genre, author_id, description, cover_image from books where author_id = $1
`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book types.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Genre, &book.AuthorId, &book.Description, &book.Image)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (d *DB) GetChaptersByBookId(ctx context.Context, id int64) ([]*types.Chapter, error) {
	rows, err := d.Pool.Query(ctx, `
	select id, book_id, number, name from chapters where book_id = $1
`, id)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	chapters := make([]*types.Chapter, 0)
	for rows.Next() {
		var chapter types.Chapter
		err := rows.Scan(&chapter.ID, &chapter.BookId, &chapter.Number, &chapter.Name)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		chapters = append(chapters, &chapter)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return chapters, nil
}

func (d *DB) ReadChapter(ctx context.Context, id int64) (*types.Chapter, error) {
	var chapter types.Chapter
	err := d.Pool.QueryRow(ctx, `
	select id, book_id, number, name, content, active, created from chapters where id = $1
`, id).Scan(&chapter.ID, &chapter.BookId, &chapter.Number,&chapter.Name, &chapter.Content, &chapter.Active, &chapter.Created)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return &chapter, nil
}

func (d *DB) EditTitle(ctx context.Context,id int64, title string) error {
	_, err := d.Pool.Exec(ctx, `
		update books set title = $1 where id = $2
`, title, id)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}