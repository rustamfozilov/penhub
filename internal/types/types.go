package types

import "time"

type Book struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	AuthorId    int64     `json:"author_id"`
	Description string    `json:"description"`
	Image       string    `json:"cover_image"`
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
	ID int64 `json:"id"`
	BookId int64 `json:"book_id"`
	Number int64 `json:"number"`
	Name string `json:"name"`
	Content string `json:"content"`
	Active   bool      `json:"active"`
	Created time.Time `json:"created"`
}

type BookTitle struct {
	Title string `json:"title"`
}