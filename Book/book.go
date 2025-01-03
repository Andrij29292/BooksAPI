package book

import "time"

type Book struct {
	ID           int64     `json:"id,omitempty"`
	Name         string    `json:"name"`
	Author       string    `json:"author"`
	PagesCount   int       `json:"pagesCount"`
	RegisteredAt time.Time `json:"registeredAt,omitempty"`
}
