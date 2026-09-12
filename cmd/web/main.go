package main

import (
	"context"
	"errors"
	"golang-rest-api/internal/cache"
	"golang-rest-api/internal/config"
	"golang-rest-api/internal/database"
	"golang-rest-api/internal/feature/auth"
	"golang-rest-api/internal/provider"
	"golang-rest-api/internal/response"
	"golang-rest-api/internal/routes"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/joho/godotenv"

	applog "golang-rest-api/internal/logger"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
	log.Println("aplikasi berhenti dengan bersih")
}

func run() error {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logapps := applog.New(cfg.Log)
	slog.SetDefault(logapps)
	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DB)
	if err != nil {
		return err
	}
	defer pool.Close()

	logapps.Info("database terhubung",
		slog.String("host", cfg.DB.Host),
		slog.String("db", cfg.DB.Name),
	)

	app := fiber.New(fiber.Config{
		ErrorHandler: response.ErrorHandler,
	})

	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(applog.Middleware(
		logapps,
		cfg.Log.RequestBody,
	))
	// nanti bisa setting disini
	app.Use(cors.New(cors.Config{}))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Setyabudi Dwisandi Arifin")
	})

	tokens := auth.NewTokenManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessTokenTTL,
		cfg.Auth.RefreshTokenTTL,
	)

	var limiter provider.RateLimiter = provider.NoopLimiter{}
	if cfg.RateLimit.Enabled {
		redisClient, err := cache.NewClient(ctx, cfg.Redis)
		if err != nil {
			return err
		}

		defer func() {
			if err := redisClient.Close(); err != nil {
				logapps.Warn("gagal menutup redis", slog.String("error", err.Error()))
			}
		}()

		logapps.Info("redis terhubung", slog.String("addr", cfg.Redis.Addr))
		limiter = cache.NewRedisLimiter(redisClient)
	}

	routes.SetUpRouter(
		app.Group("/api"),
		provider.Deps{
			DB:          pool,
			Tx:          database.NewTransactor(pool),
			Log:         logapps,
			Tokens:      tokens,
			Limiter:     limiter,
			RateLimiter: cfg.RateLimit,
		},
	)

	addr := ":" + cfg.App.Port
	// Server jalan di goroutine terpisah
	serverErr := make(chan error, 1)
	go func() {
		if err := app.Listen(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	logapps.Info("server jalan", slog.String("addr", addr), slog.String("env", cfg.App.Env))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-stop:
		log.Printf("sinyal diterima: %s, memulai shutdown...", sig)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		return err
	}

	logapps.Info("server ditutup, menutup koneksi database")
	return nil
}
