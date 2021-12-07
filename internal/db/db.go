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

var ErrNotFound = errors.New("not found")

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
	insert into books (title, author_id, description, cover_image, access_read, genre_id, active, created)
	 values ($1, $2, $3, $4, $5,$6, default, default)   
`, book.Title, book.AuthorId, book.Description, book.Image, book.AccessRead, book.Genre)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *DB) RegistrationUser(ctx context.Context, user *types.User, hash []byte) error {

	_, err := d.Pool.Exec(ctx, `
		insert into users (name, login, password, active, created)
		values ($1, $2, $3, default, default)
`, user.Name, user.Login, hash)
	if err != nil {
		err:= errors.WithStack(err)

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
		err:= errors.WithStack(err)
		return true, 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println("invalid password")
		err := errors.WithStack(err)
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
		err := errors.WithStack(err)
		return err
	}
	return nil
}

func (d *DB) IdByToken(cxt context.Context, token string) (id int64, expire time.Time, err error) {
	err = d.Pool.QueryRow(cxt, `
select user_id, expire from users_tokens where token = $1
`, token).Scan(&id, &expire)
	if err != nil {
		err := errors.WithStack(err)
		return 0, expire, err
	}
	return id, expire, nil
}

func (d *DB) GetBookId(ctx context.Context, title string) (id int64, err error) {

	err = d.Pool.QueryRow(ctx, `
		select id from books where title = $1
`, title).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		err := errors.WithStack(err)
		log.Println("book not exist")
		return 0, err
	}
	if err != nil {
		err := errors.WithStack(err)
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
		err := errors.WithStack(err)
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
		err:= errors.WithStack(err)
		return err
	}
	return nil
}

func (d *DB) GetBooksById(ctx context.Context, id int64) ([]*types.Book, error) {

	books := make([]*types.Book, 0)

	rows, err := d.Pool.Query(ctx, `
		select id, title, genre_id, author_id, description, cover_image, active, created from books where author_id = $1
`, id)
	if err != nil {
		err:= errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book types.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Genre, &book.AuthorId, &book.Description, &book.Image, &book.Active, &book.Created)
		if err != nil {
			err:= errors.WithStack(err)
			return nil, err
		}
			if book.Active{
				books = append(books, &book)
			}
	}
	err = rows.Err()
	if err != nil {
		err:= errors.WithStack(err)
		return nil, err
	}
	return books, nil
}

func (d *DB) GetChaptersByBookId(ctx context.Context, id int64) ([]*types.Chapter, error) {
	rows, err := d.Pool.Query(ctx, `
	select id, book_id, number, name, active, created from chapters where book_id = $1
		order by number 
`, id)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	chapters := make([]*types.Chapter, 0)
	for rows.Next() {
		var chapter types.Chapter
		err := rows.Scan(&chapter.ID, &chapter.BookId, &chapter.Number, &chapter.Name, &chapter.Active, &chapter.Created)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		if chapter.Active {
			chapters = append(chapters, &chapter)
		}
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
`, id).Scan(&chapter.ID, &chapter.BookId, &chapter.Number, &chapter.Name, &chapter.Content, &chapter.Active, &chapter.Created)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &chapter, nil
}

func (d *DB) EditTitle(ctx context.Context, id int64, title string) error {
	_, err := d.Pool.Exec(ctx, `
		update books set title = $1 where id = $2
`, title, id)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (d *DB) EditContent(ctx context.Context, edit *types.Chapter) error {

	_, err := d.Pool.Exec(ctx, `
			update chapters set content = $1 where id = $2
`, edit.Content, edit.ID)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (d *DB) EditAccess(ctx context.Context, edit *types.Book) error {

	_, err := d.Pool.Exec(ctx, `
				update books set  access_read = $1 where id = $2
`, edit.AccessRead, edit.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *DB) EditChapterName(ctx context.Context, edit *types.Chapter) error   {

	_, err := d.Pool.Exec(ctx, `
			update chapters set name = $1 where id = $2
`, edit.Name, edit.ID)
	if err != nil {
		return errors.WithStack(err)
	}
return nil
}



// ????? возвращает пустой слайс вместо NoRows, при неправильном поисковом запросе?????
func (d *DB) SearchByTitle(ctx context.Context, title *types.BookTitle) ([]*types.Book, error)  {
		books := make([]*types.Book, 0)
	rows, err := d.Pool.Query(ctx, `
			select *from books where "like"(title, $1) 
`,title.Title + "%")
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book types.Book

		err := rows.Scan(&book.ID, &book.Title, &book.AuthorId, &book.Genre, &book.Description, &book.Image, &book.AccessRead, &book.Active, &book.Created)
		if err != nil {
			err := errors.WithStack(err)
				return nil, err
		}
		books = append(books, &book)
	}
	err = rows.Err()
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	return books, nil
}

func (d *DB) SearchByAuthor(ctx context.Context, author *types.AuthorName) ([]*types.User, error )  {
	authors := make([]*types.User, 0)
	rows, err := d.Pool.Query(ctx, `
			select id, name, active, created from users where "like"(name, $1) 
`,author.Name + "%")
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user types.User

		err := rows.Scan(&user.ID, &user.Name, &user.Active, &user.Created)
		if err != nil {
			err := errors.WithStack(err)
			return nil, err
		}
		if user.Active{
			authors = append(authors, &user)
		}
	}
	err = rows.Err()
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	return authors, nil
}

func (d *DB) GetAllGenres(ctx context.Context) ([]*types.Genre, error) {

	genres := make([]*types.Genre, 0)

	rows, err := d.Pool.Query(ctx, `
	select *from genres
`)
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre types.Genre
		err := rows.Scan(&genre.Id, &genre.Name, &genre.Active)
		if err != nil {
			err := errors.WithStack(err)
			return nil, err
		}

		if genre.Active {
			genres = append(genres, &genre)
		}
	}
return genres, nil
}

func (d *DB) SearchGenre(ctx context.Context, genreName types.GenreName) ([]*types.Genre, error)  {

	genres := make([]*types.Genre, 0)

	rows, err := d.Pool.Query(ctx, `
	select *from genres where "like"(name, $1)
`, genreName.Name+"%")
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre types.Genre
		err := rows.Scan(&genre.Id, &genre.Name, &genre.Active)
		if err != nil {
			err := errors.WithStack(err)
			return nil, err
		}

		if genre.Active {
			genres = append(genres, &genre)
		}
	}
	return genres, nil
}

func (d *DB) GetBooksByGenreId(ctx context.Context, genreId types.GenreID) ([]*types.Book,error){

	books := make([]*types.Book, 0)

	rows, err := d.Pool.Query(ctx, `
		select id, title, genre_id, author_id, description, cover_image, active, created from books where genre_id = $1
`, genreId.Id)
	if err != nil {
		err:= errors.WithStack(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book types.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Genre, &book.AuthorId, &book.Description, &book.Image, &book.Active, &book.Created)
		if err != nil {
			err:= errors.WithStack(err)
			return nil, err
		}
		if book.Active{
			books = append(books, &book)
		}
	}
	err = rows.Err()
	if err != nil {
		err:= errors.WithStack(err)
		return nil, err
	}
	return books, nil
}

func (d *DB) GetGenreById(ctx context.Context, genreId types.GenreID) (*types.Genre, error)  {
	var genre types.Genre

	err := d.Pool.QueryRow(ctx, `
		select *from genres where id = $1
`, genreId.Id).Scan(&genre.Id, &genre.Name, &genre.Active)
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
		return &genre, nil
}

func (d *DB) EditImage(ctx context.Context, book *types.Book) error  {
	_, err := d.Pool.Exec(ctx, `
				update books set cover_image = $1 where id = $2
`, book.Image, book.ID)
	if err != nil {
		err := errors.WithStack(err)
		return err
	}
return nil
}