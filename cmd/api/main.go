package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"letter-square-api/config"
	"letter-square-api/internal/handler"
	"letter-square-api/internal/middleware"
	"letter-square-api/internal/repository"
	"letter-square-api/internal/service"
	"letter-square-api/internal/thirdparty/gemini"
)

func main() {
	cfg := config.Load()

	// Database
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	log.Println("database connected")

	// Repositories
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	watchlistRepo := repository.NewWatchlistRepository(db)
	recSessionRepo := repository.NewRecommendationSessionRepository(db)

	// Third-party
	geminiClient := gemini.NewClient(cfg.GeminiAPIKey, cfg.GeminiModel)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHour)
	movieSvc := service.NewMovieService(movieRepo)
	reviewSvc := service.NewReviewService(reviewRepo, movieRepo)
	wlSvc := service.NewWatchlistService(watchlistRepo, movieRepo)
	recSvc := service.NewRecommendationService(recSessionRepo, watchlistRepo, geminiClient)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	movieH := handler.NewMovieHandler(movieSvc)
	reviewH := handler.NewReviewHandler(reviewSvc)
	wlH := handler.NewWatchlistHandler(wlSvc)
	recH := handler.NewRecommendationHandler(recSvc)

	// Router
	mux := http.NewServeMux()
	auth := middleware.Auth(cfg.JWTSecret)

	// Auth
	mux.HandleFunc("POST /api/auth/register", authH.Register)
	mux.HandleFunc("POST /api/auth/login", authH.Login)

	// Movies (public read, auth write)
	mux.HandleFunc("GET /api/movies", movieH.GetAll)
	mux.HandleFunc("GET /api/movies/{id}", movieH.GetByID)
	mux.Handle("POST /api/movies", auth(http.HandlerFunc(movieH.Create)))
	mux.Handle("PUT /api/movies/{id}", auth(http.HandlerFunc(movieH.Update)))
	mux.Handle("DELETE /api/movies/{id}", auth(http.HandlerFunc(movieH.Delete)))

	// Reviews
	mux.HandleFunc("GET /api/movies/{id}/reviews", reviewH.GetByMovie)
	mux.Handle("POST /api/movies/{id}/reviews", auth(http.HandlerFunc(reviewH.Create)))
	mux.Handle("GET /api/users/me/reviews", auth(http.HandlerFunc(reviewH.GetMyReviews)))
	mux.Handle("PUT /api/reviews/{reviewId}", auth(http.HandlerFunc(reviewH.Update)))
	mux.Handle("DELETE /api/reviews/{reviewId}", auth(http.HandlerFunc(reviewH.Delete)))

	// Watchlist
	mux.Handle("GET /api/watchlist", auth(http.HandlerFunc(wlH.GetMyWatchlist)))
	mux.Handle("POST /api/watchlist/{movieId}", auth(http.HandlerFunc(wlH.Add)))
	mux.Handle("DELETE /api/watchlist/{movieId}", auth(http.HandlerFunc(wlH.Remove)))

	// Gemini AI Recommendation (session-based Q&A flow)
	mux.Handle("POST /api/recommendations/start", auth(http.HandlerFunc(recH.Start)))
	mux.Handle("POST /api/recommendations/sessions/{id}/answer", auth(http.HandlerFunc(recH.Answer)))
	mux.Handle("GET /api/recommendations/sessions/{id}", auth(http.HandlerFunc(recH.GetSession)))

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("server running on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
