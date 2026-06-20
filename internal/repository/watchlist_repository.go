package repository

import (
	"database/sql"
	"letter-square-api/internal/entity"
)

type WatchlistRepository interface {
	Add(userID, movieID int64) error
	Remove(userID, movieID int64) error
	FindByUserID(userID int64, limit, offset int) ([]*entity.Watchlist, error)
	Exists(userID, movieID int64) (bool, error)
}

type watchlistRepository struct {
	db *sql.DB
}

func NewWatchlistRepository(db *sql.DB) WatchlistRepository {
	return &watchlistRepository{db: db}
}

func (r *watchlistRepository) Add(userID, movieID int64) error {
	_, err := r.db.Exec(`INSERT INTO watchlists (user_id, movie_id) VALUES (?, ?)`, userID, movieID)
	return err
}

func (r *watchlistRepository) Remove(userID, movieID int64) error {
	_, err := r.db.Exec(`DELETE FROM watchlists WHERE user_id=? AND movie_id=?`, userID, movieID)
	return err
}

func (r *watchlistRepository) FindByUserID(userID int64, limit, offset int) ([]*entity.Watchlist, error) {
	q := `
		SELECT w.id, w.user_id, w.movie_id, w.created_at,
		       m.id, m.title, m.director, m.genre, m.year, m.synopsis, m.poster_url, m.created_at, m.updated_at
		FROM watchlists w
		JOIN movies m ON m.id = w.movie_id
		WHERE w.user_id = ?
		ORDER BY w.created_at DESC
		LIMIT ? OFFSET ?`
	rows, err := r.db.Query(q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.Watchlist
	for rows.Next() {
		w := &entity.Watchlist{Movie: &entity.Movie{}}
		if err := rows.Scan(
			&w.ID, &w.UserID, &w.MovieID, &w.CreatedAt,
			&w.Movie.ID, &w.Movie.Title, &w.Movie.Director, &w.Movie.Genre,
			&w.Movie.Year, &w.Movie.Synopsis, &w.Movie.PosterURL,
			&w.Movie.CreatedAt, &w.Movie.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, w)
	}
	return items, nil
}

func (r *watchlistRepository) Exists(userID, movieID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM watchlists WHERE user_id=? AND movie_id=?`, userID, movieID).Scan(&count)
	return count > 0, err
}
