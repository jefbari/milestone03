package repository

import (
	"database/sql"
	"errors"
	"letter-square-api/internal/entity"
)

type UserRepository interface {
	Create(u *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id int64) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u *entity.User) error {
	q := `INSERT INTO users (username, email, password, bio) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(q, u.Username, u.Email, u.Password, u.Bio)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	u := &entity.User{}
	q := `SELECT id, username, email, password, bio, created_at, updated_at FROM users WHERE email = ? LIMIT 1`
	err := r.db.QueryRow(q, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Bio, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *userRepository) FindByID(id int64) (*entity.User, error) {
	u := &entity.User{}
	q := `SELECT id, username, email, password, bio, created_at, updated_at FROM users WHERE id = ? LIMIT 1`
	err := r.db.QueryRow(q, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Bio, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	u := &entity.User{}
	q := `SELECT id, username, email, password, bio, created_at, updated_at FROM users WHERE username = ? LIMIT 1`
	err := r.db.QueryRow(q, username).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Bio, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}
