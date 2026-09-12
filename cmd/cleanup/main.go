package main

import (
	"context"
	"golang-rest-api/internal/config"
	"golang-rest-api/internal/database"
	"golang-rest-api/internal/feature/auth"
	"golang-rest-api/internal/logger"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("cleanup gagal", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	applog := logger.New(cfg.Log)
	slog.SetDefault(applog)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	dbCfg := cfg.DB
	dbCfg.MaxConns = 1
	dbCfg.MinConns = 1
	dbCfg.AppName = "cleanup"

	pool, err := database.NewPool(ctx, dbCfg)
	if err != nil {
		return err
	}

	defer pool.Close()

	n, err := auth.NewAuthRepository(pool).DeleteExpired(ctx)
	if err != nil {
		return err
	}

	applog.Info("refresh token dibersihkan", slog.Int64("jumlah", n))
	return nil
}
