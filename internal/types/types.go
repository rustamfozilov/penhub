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
