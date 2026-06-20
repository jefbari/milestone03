package repository

import (
	"database/sql"
	"errors"
	"letter-square-api/internal/entity"
)

type ReviewRepository interface {
	Create(r *entity.Review) error
	FindByID(id int64) (*entity.Review, error)
	FindByMovieID(movieID int64, limit, offset int) ([]*entity.Review, error)
	FindByUserID(userID int64, limit, offset int) ([]*entity.Review, error)
	Update(r *entity.Review) error
	Delete(id int64) error
	ExistsByUserAndMovie(userID, movieID int64) (bool, error)
}

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(rv *entity.Review) error {
	q := `INSERT INTO reviews (user_id, movie_id, rating, body) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(q, rv.UserID, rv.MovieID, rv.Rating, rv.Body)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	rv.ID = id
	return nil
}

func (r *reviewRepository) FindByID(id int64) (*entity.Review, error) {
	rv := &entity.Review{}
	q := `
		SELECT rv.id, rv.user_id, rv.movie_id, rv.rating, rv.body, rv.created_at, rv.updated_at,
		       u.username, m.title
		FROM reviews rv
		JOIN users u ON u.id = rv.user_id
		JOIN movies m ON m.id = rv.movie_id
		WHERE rv.id = ?`
	err := r.db.QueryRow(q, id).Scan(
		&rv.ID, &rv.UserID, &rv.MovieID, &rv.Rating, &rv.Body, &rv.CreatedAt, &rv.UpdatedAt,
		&rv.Username, &rv.MovieTitle,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return rv, err
}

func (r *reviewRepository) FindByMovieID(movieID int64, limit, offset int) ([]*entity.Review, error) {
	q := `
		SELECT rv.id, rv.user_id, rv.movie_id, rv.rating, rv.body, rv.created_at, rv.updated_at,
		       u.username, m.title
		FROM reviews rv
		JOIN users u ON u.id = rv.user_id
		JOIN movies m ON m.id = rv.movie_id
		WHERE rv.movie_id = ?
		ORDER BY rv.created_at DESC
		LIMIT ? OFFSET ?`
	return r.scanRows(r.db.Query(q, movieID, limit, offset))
}

func (r *reviewRepository) FindByUserID(userID int64, limit, offset int) ([]*entity.Review, error) {
	q := `
		SELECT rv.id, rv.user_id, rv.movie_id, rv.rating, rv.body, rv.created_at, rv.updated_at,
		       u.username, m.title
		FROM reviews rv
		JOIN users u ON u.id = rv.user_id
		JOIN movies m ON m.id = rv.movie_id
		WHERE rv.user_id = ?
		ORDER BY rv.created_at DESC
		LIMIT ? OFFSET ?`
	return r.scanRows(r.db.Query(q, userID, limit, offset))
}

func (r *reviewRepository) scanRows(rows *sql.Rows, err error) ([]*entity.Review, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []*entity.Review
	for rows.Next() {
		rv := &entity.Review{}
		if err := rows.Scan(&rv.ID, &rv.UserID, &rv.MovieID, &rv.Rating, &rv.Body,
			&rv.CreatedAt, &rv.UpdatedAt, &rv.Username, &rv.MovieTitle); err != nil {
			return nil, err
		}
		reviews = append(reviews, rv)
	}
	return reviews, nil
}

func (r *reviewRepository) Update(rv *entity.Review) error {
	q := `UPDATE reviews SET rating=?, body=?, updated_at=NOW() WHERE id=?`
	_, err := r.db.Exec(q, rv.Rating, rv.Body, rv.ID)
	return err
}

func (r *reviewRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM reviews WHERE id = ?`, id)
	return err
}

func (r *reviewRepository) ExistsByUserAndMovie(userID, movieID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM reviews WHERE user_id=? AND movie_id=?`, userID, movieID).Scan(&count)
	return count > 0, err
}
