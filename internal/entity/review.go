package entity

import "time"

type Review struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	MovieID   int64     `json:"movie_id"`
	Rating    float64   `json:"rating"` // 1.0 - 5.0
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Joins
	Username  string  `json:"username,omitempty"`
	MovieTitle string `json:"movie_title,omitempty"`
}
