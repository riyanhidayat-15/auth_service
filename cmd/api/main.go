package main

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/services"
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main () {
	_ = godotenv.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Gagal load config", "error", err)
		os.Exit(1)
	}

	pool, err := database.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("Koneksi gagal", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("Terhubung ke database")

	userRepo := repository.NewPostgresUserRepository(pool)
	authSvc := services.NewAuthService(userRepo)
	authhandler := handler.NewAuthHandler(authSvc)

	router := handler.NewRouter(authhandler)

	logger.Info("Server up and running", "port", cfg.Port)
	
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.Error("Server berhenti", "error", err)
		os.Exit(1)
	}
}
