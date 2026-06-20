package entity

import "time"

// RecommendationSession represents one Q&A flow toward a Gemini-generated
// movie recommendation, based on the user's watchlist + their answers.
type RecommendationSession struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"` // "in_progress" | "completed"
	Step      int       `json:"step"`   // index of the last answered question (0-based)
	Answers   string    `json:"-"`      // raw JSON array string, stored in DB
	Result    string    `json:"result,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
