package repository

import (
	"database/sql"
	"errors"
	"letter-square-api/internal/entity"
)

type RecommendationSessionRepository interface {
	Create(s *entity.RecommendationSession) error
	FindByID(id int64) (*entity.RecommendationSession, error)
	Update(s *entity.RecommendationSession) error
}

type recommendationSessionRepository struct {
	db *sql.DB
}

func NewRecommendationSessionRepository(db *sql.DB) RecommendationSessionRepository {
	return &recommendationSessionRepository{db: db}
}

func (r *recommendationSessionRepository) Create(s *entity.RecommendationSession) error {
	if s.Answers == "" {
		s.Answers = "[]"
	}
	q := `INSERT INTO recommendation_sessions (user_id, status, step, answers, result) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.Exec(q, s.UserID, s.Status, s.Step, s.Answers, s.Result)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	s.ID = id
	return nil
}

func (r *recommendationSessionRepository) FindByID(id int64) (*entity.RecommendationSession, error) {
	s := &entity.RecommendationSession{}
	q := `SELECT id, user_id, status, step, answers, result, created_at, updated_at
	      FROM recommendation_sessions WHERE id = ?`
	err := r.db.QueryRow(q, id).Scan(
		&s.ID, &s.UserID, &s.Status, &s.Step, &s.Answers, &s.Result, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return s, err
}

func (r *recommendationSessionRepository) Update(s *entity.RecommendationSession) error {
	q := `UPDATE recommendation_sessions SET status=?, step=?, answers=?, result=?, updated_at=NOW() WHERE id=?`
	_, err := r.db.Exec(q, s.Status, s.Step, s.Answers, s.Result, s.ID)
	return err
}
