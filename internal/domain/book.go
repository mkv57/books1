package domain

import "time"

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int       `json:"user_id"`
}
