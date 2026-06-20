package service

import (
	"letter-square-api/internal/apperror"
	"letter-square-api/internal/dto"
	"letter-square-api/internal/entity"
	"letter-square-api/internal/repository"
)

type MovieService interface {
	Create(req *dto.CreateMovieRequest) (*entity.Movie, error)
	GetByID(id int64) (*entity.Movie, error)
	GetAll(search, genre string, page, limit int) ([]*entity.Movie, error)
	Update(id int64, req *dto.UpdateMovieRequest) (*entity.Movie, error)
	Delete(id int64) error
}

type movieService struct {
	movieRepo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) MovieService {
	return &movieService{movieRepo: repo}
}

func (s *movieService) Create(req *dto.CreateMovieRequest) (*entity.Movie, error) {
	m := &entity.Movie{
		Title:     req.Title,
		Director:  req.Director,
		Genre:     req.Genre,
		Year:      req.Year,
		Synopsis:  req.Synopsis,
		PosterURL: req.PosterURL,
	}
	if err := s.movieRepo.Create(m); err != nil {
		return nil, apperror.ErrInternal
	}
	return m, nil
}

func (s *movieService) GetByID(id int64) (*entity.Movie, error) {
	m, err := s.movieRepo.FindByID(id)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if m == nil {
		return nil, apperror.ErrMovieNotFound
	}
	return m, nil
}

func (s *movieService) GetAll(search, genre string, page, limit int) ([]*entity.Movie, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	movies, err := s.movieRepo.FindAll(search, genre, limit, offset)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return movies, nil
}

func (s *movieService) Update(id int64, req *dto.UpdateMovieRequest) (*entity.Movie, error) {
	m, err := s.movieRepo.FindByID(id)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if m == nil {
		return nil, apperror.ErrMovieNotFound
	}

	if req.Title != nil {
		m.Title = *req.Title
	}
	if req.Director != nil {
		m.Director = *req.Director
	}
	if req.Genre != nil {
		m.Genre = *req.Genre
	}
	if req.Year != nil {
		m.Year = *req.Year
	}
	if req.Synopsis != nil {
		m.Synopsis = *req.Synopsis
	}
	if req.PosterURL != nil {
		m.PosterURL = *req.PosterURL
	}

	if err := s.movieRepo.Update(m); err != nil {
		return nil, apperror.ErrInternal
	}
	return m, nil
}

func (s *movieService) Delete(id int64) error {
	m, err := s.movieRepo.FindByID(id)
	if err != nil {
		return apperror.ErrInternal
	}
	if m == nil {
		return apperror.ErrMovieNotFound
	}
	if err := s.movieRepo.Delete(id); err != nil {
		return apperror.ErrInternal
	}
	return nil
}
