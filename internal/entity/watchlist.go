package entity

import "time"

type Watchlist struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	MovieID   int64     `json:"movie_id"`
	CreatedAt time.Time `json:"created_at"`

	// Joins
	Movie *Movie `json:"movie,omitempty"`
}
