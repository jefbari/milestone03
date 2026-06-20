package entity

import "time"

type Movie struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Genre       string    `json:"genre"`
	Year        int       `json:"year"`
	Synopsis    string    `json:"synopsis"`
	PosterURL   string    `json:"poster_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
