package types

import "time"

type Book struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	AuthorId    int64     `json:"-"`
	Genre       int64     `json:"genre"`
	Description string    `json:"description"`
	Image       string    `json:"cover_image_name"`
	AccessRead  bool      `json:"access_read"`
	Active      bool      `json:"active"`
	Created     time.Time `json:"created"`
}

type User struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

type T struct {
	Token string `json:"token"`
}

type Chapter struct {
	ID      int64     `json:"id"`
	BookId  int64     `json:"book_id"`
	Number  int64     `json:"number"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

type Genre struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type BookTitle struct {
	Title string `json:"title"`
}

type BookId struct {
	Id int64 `json:"book_id"`
}

type ChapterId struct {
	Id int64 `json:"chapter_id"`
}

type AuthorName struct {
	Name string `json:"author"`
}

type AuthorId struct {
	Id         int64 `json:"author_id"`
	LastBookId int64 `json:"last_id"`
}

type GenreName struct {
	Name string `json:"genre_name"`
}

type GenreID struct {
	Id         int64 `json:"genre_id"`
	LastBookId int64 `json:"last_id"`
}

type ImageName struct {
	Name string `json:"image_name"`
}

type RatingID struct {
	Id int64 `json:"like_id"`
}

type Config struct {
	UserName   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Database   string `json:"database"`
	ImagesPath string `json:"images_path"`
}
