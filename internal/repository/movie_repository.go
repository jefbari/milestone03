package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"letter-square-api/internal/entity"
)

type MovieRepository interface {
	Create(m *entity.Movie) error
	FindByID(id int64) (*entity.Movie, error)
	FindAll(search, genre string, limit, offset int) ([]*entity.Movie, error)
	Update(m *entity.Movie) error
	Delete(id int64) error
	FindTopRated(limit int) ([]*entity.Movie, error)
}

type movieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) Create(m *entity.Movie) error {
	q := `INSERT INTO movies (title, director, genre, year, synopsis, poster_url) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.Exec(q, m.Title, m.Director, m.Genre, m.Year, m.Synopsis, m.PosterURL)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	m.ID = id
	return nil
}

func (r *movieRepository) FindByID(id int64) (*entity.Movie, error) {
	m := &entity.Movie{}
	q := `SELECT id, title, director, genre, year, synopsis, poster_url, created_at, updated_at FROM movies WHERE id = ?`
	err := r.db.QueryRow(q, id).Scan(&m.ID, &m.Title, &m.Director, &m.Genre, &m.Year, &m.Synopsis, &m.PosterURL, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return m, err
}

func (r *movieRepository) FindAll(search, genre string, limit, offset int) ([]*entity.Movie, error) {
	q := `SELECT id, title, director, genre, year, synopsis, poster_url, created_at, updated_at FROM movies WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		q += ` AND (title LIKE ? OR director LIKE ?)`
		like := fmt.Sprintf("%%%s%%", search)
		args = append(args, like, like)
	}
	if genre != "" {
		q += ` AND genre = ?`
		args = append(args, genre)
	}
	q += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*entity.Movie
	for rows.Next() {
		m := &entity.Movie{}
		if err := rows.Scan(&m.ID, &m.Title, &m.Director, &m.Genre, &m.Year, &m.Synopsis, &m.PosterURL, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *movieRepository) Update(m *entity.Movie) error {
	q := `UPDATE movies SET title=?, director=?, genre=?, year=?, synopsis=?, poster_url=?, updated_at=NOW() WHERE id=?`
	_, err := r.db.Exec(q, m.Title, m.Director, m.Genre, m.Year, m.Synopsis, m.PosterURL, m.ID)
	return err
}

func (r *movieRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM movies WHERE id = ?`, id)
	return err
}

func (r *movieRepository) FindTopRated(limit int) ([]*entity.Movie, error) {
	q := `
		SELECT m.id, m.title, m.director, m.genre, m.year, m.synopsis, m.poster_url, m.created_at, m.updated_at
		FROM movies m
		JOIN reviews r ON r.movie_id = m.id
		GROUP BY m.id
		ORDER BY AVG(r.rating) DESC
		LIMIT ?`
	rows, err := r.db.Query(q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*entity.Movie
	for rows.Next() {
		m := &entity.Movie{}
		if err := rows.Scan(&m.ID, &m.Title, &m.Director, &m.Genre, &m.Year, &m.Synopsis, &m.PosterURL, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}
