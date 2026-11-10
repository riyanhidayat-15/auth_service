package main

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main () {
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

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	logger.Info("Server up and running", "port", cfg.Port)
	
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logger.Error("Server berhenti", "error", err)
		os.Exit(1)
	}
}
