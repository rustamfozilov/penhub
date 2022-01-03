package services

import (
	"context"
	"github.com/google/uuid"
	errors "github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/db"
	"github.com/rustamfozilov/penhub/internal/types"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	db            *db.DB
	imagesDirPath string
}

func NewService(db *db.DB, imagesDirPath string) *Service {
	return &Service{db: db, imagesDirPath: imagesDirPath}
}



var ErrExpired = errors.New("token expired")
var ErrInternal = errors.New("internal error")
var ErrLoginUsed = errors.New("login already registered")
var ErrNoSuchUser = errors.New("no such user")
var ErrInvalidPassword = errors.New("invalid password")
var ErrNotFound = errors.New("not found")
var ErrInvalidData = errors.New("invalid data")

func (s *Service) CreateBook(ctx context.Context, book *types.Book) error {
	return s.db.CreateBook(ctx, book)
}

func (s *Service) RegistrationUser(ctx context.Context, user *types.User) error {
	if s.db.IsLoginUsed(ctx, user.Login) {
		return ErrLoginUsed
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}
	err = s.db.RegistrationUser(ctx, user, hash)
	if err != nil {
		err := errors.WithStack(err)
		return err
	}
	return nil
}

func (s *Service) GetBookId(ctx context.Context, bookName *types.BookTitle) (id int64, err error) {
	return s.db.GetBookId(ctx, bookName.Title)
}

func (s *Service) HaveAccessToEditBook(ctx context.Context, userId, bookId int64) (bool, error) {

	return s.db.BookAccess(ctx, userId, bookId)
}

func (s *Service) WriteChapter(ctx context.Context, chapter *types.Chapter) error {
	return s.db.WriteChapter(ctx, chapter)
}

func (s *Service) GetBooksById(ctx context.Context, id *types.AuthorId) ([]*types.Book, error) {
	return s.db.GetBooksById(ctx, id)
}

func (s *Service) GetChaptersByBookId(ctx context.Context, bookId *types.BookId) ([]*types.Chapter, error) {
	return s.db.GetChaptersByBookId(ctx, bookId.Id)
}

func (s *Service) ReadChapter(ctx context.Context, chapterId *types.ChapterId) (*types.Chapter, error) {
	return s.db.ReadChapter(ctx, chapterId.Id)

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

func (s *Service) EditChapterNumber(ctx context.Context, edit *types.Chapter) error {
	return s.db.EditChapterNumber(ctx, edit)
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

func (s *Service) GetBooksByGenreId(ctx context.Context, genreId *types.GenreID) ([]*types.Book, error) {
	return s.db.GetBooksByGenreId(ctx, genreId)
}

func (s *Service) GetGenreById(ctx context.Context, genreId types.GenreID) (*types.Genre, error) {
	return s.db.GetGenreById(ctx, genreId)
}

func (s *Service) SaveImage(file io.Reader, fileName string, book *types.Book) (*types.Book, error) {

	extension := fileName[len(fileName)-4:]
	imageName := uuid.New().String()
	path := filepath.Join(s.imagesDirPath, imageName+extension)
	imageFile, err := os.Create(path)
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	defer imageFile.Close()

	_, err = io.Copy(imageFile, file)
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	book.Image = imageName + extension
	return book, nil
}

func (s *Service) EditImage(ctx context.Context, book *types.Book) error {
	return s.db.EditImageName(ctx, book)
}

func (s *Service) GetImageByName(name string) ([]byte, error) {
	path := filepath.Join(s.imagesDirPath, name)
	file, err := os.ReadFile(path)
	if err != nil {
		err := errors.WithStack(err)
		return nil, err
	}
	return file, nil
}

func (s *Service) EditGenre(ctx context.Context, book *types.Book) error {
	return s.db.EditGenre(ctx, book)
}

func (s *Service) EditDescription(ctx context.Context, book *types.Book) error {
	return s.db.EditDescription(ctx, book)
}

func (s *Service) DeleteBook(ctx context.Context, book *types.Book) error {
	return s.db.DeleteBook(ctx, book)
}
func (s *Service) DeleteChapter(ctx context.Context, chapter *types.Chapter) error {
	return s.db.DeleteChapter(ctx, chapter)
}

func (s *Service) AddLike(ctx context.Context, userId *int64, id *types.BookId) error {
	return s.db.AddLike(ctx, userId, id)
}

func (s *Service) DeleteLike(ctx context.Context, id *types.RatingID) error {
	return s.db.DeleteLike(ctx, id)
}

func (s *Service) GetLikeId(ctx context.Context, userId *int64, id *types.BookId) (*types.RatingID, error) {
	return s.db.GetLikeId(ctx, userId, id)
}

func (s *Service) BookLikes(ctx context.Context, id *types.BookId) (int64, error) {
	return s.db.BookLikes(ctx, id)
}

func (s *Service) ValidateUser(user *types.User) error {
	if len(user.Name) > 20 || len(user.Name) < 3 {
		return ErrInvalidData
	}
	if len(user.Login) > 20 || len(user.Login) < 3 {
		return ErrInvalidData
	}
	if len(user.Password) > 20 || len(user.Password) < 6 {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, "_") || strings.Contains(user.Password, "-") {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, "@") || strings.Contains(user.Password, "#") {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, "$") || strings.Contains(user.Password, "%") {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, "&") || strings.Contains(user.Password, "*") {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, "(") || strings.Contains(user.Password, ")") {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, ":") || strings.Contains(user.Password, ".") {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, "/") || strings.Contains(user.Password, `\`) {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, ",") || strings.Contains(user.Password, ";") {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, "?") || strings.Contains(user.Password, `"`) {
		return ErrInvalidData
	}
	if strings.Contains(user.Password, "!") || strings.Contains(user.Password, "~") {
		return ErrInvalidData
	}
	return nil
}

func (s *Service) ValidateBook(book *types.Book) error {
	if len(book.Title) > 20 || len(book.Title) < 1 {
		return ErrInvalidData
	}
	if len(book.Description) > 250 || len(book.Description) < 5 {
		return ErrInvalidData
	}
	return nil
}

func (s *Service) ValidateImage(size int64) error {
	if size > 5_000_000_000 {
		return ErrInvalidData
	}
	return nil
}

func (s *Service) ValidateChapter(chapter *types.Chapter) error {
	if len(chapter.Name) > 20 || len(chapter.Name) < 1 {
		return ErrInvalidData
	}
	return nil
}
func (s *Service) ValidateGenreId(genreId *types.GenreID) error {
	if genreId.Id > 200 {
		return ErrInvalidData
	}
	return nil
}

func (s *Service) ValidateTitle(title string) error {
	if len(title) > 20 || len(title) < 1 {
		return ErrInvalidData
	}
	return nil
}

func (s *Service) ValidateDescription(description string) error {
	if len(description) > 250 || len(description) < 5 {
		return ErrInvalidData
	}
	return nil
}

func (s *Service) ValidateImageName(name string) error {
	if len(name) > 50 || len(name) < 5 {
		return ErrInvalidData
	}
	return nil
}
