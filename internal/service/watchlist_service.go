package service

import (
	"letter-square-api/internal/apperror"
	"letter-square-api/internal/entity"
	"letter-square-api/internal/repository"
)

type WatchlistService interface {
	Add(userID, movieID int64) error
	Remove(userID, movieID int64) error
	GetByUser(userID int64, page, limit int) ([]*entity.Watchlist, error)
}

type watchlistService struct {
	watchlistRepo repository.WatchlistRepository
	movieRepo     repository.MovieRepository
}

func NewWatchlistService(wlRepo repository.WatchlistRepository, movieRepo repository.MovieRepository) WatchlistService {
	return &watchlistService{watchlistRepo: wlRepo, movieRepo: movieRepo}
}

func (s *watchlistService) Add(userID, movieID int64) error {
	m, err := s.movieRepo.FindByID(movieID)
	if err != nil {
		return apperror.ErrInternal
	}
	if m == nil {
		return apperror.ErrMovieNotFound
	}

	exists, err := s.watchlistRepo.Exists(userID, movieID)
	if err != nil {
		return apperror.ErrInternal
	}
	if exists {
		return apperror.ErrDuplicateWatchlist
	}

	if err := s.watchlistRepo.Add(userID, movieID); err != nil {
		return apperror.ErrInternal
	}
	return nil
}

func (s *watchlistService) Remove(userID, movieID int64) error {
	exists, err := s.watchlistRepo.Exists(userID, movieID)
	if err != nil {
		return apperror.ErrInternal
	}
	if !exists {
		return apperror.ErrNotFound
	}
	if err := s.watchlistRepo.Remove(userID, movieID); err != nil {
		return apperror.ErrInternal
	}
	return nil
}

func (s *watchlistService) GetByUser(userID int64, page, limit int) ([]*entity.Watchlist, error) {
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	items, err := s.watchlistRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return items, nil
}
