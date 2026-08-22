package main

import (
	"context"
	"errors"
	"golang-rest-api/internal/database"
	"golang-rest-api/internal/provider"
	"golang-rest-api/internal/response"
	"golang-rest-api/internal/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
	log.Println("aplikasi berhenti dengan bersih")
}

func run() error {
	_ = godotenv.Load()

	ctx := context.Background()

	pool, err := database.NewPool(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	log.Println("database terhubung")

	app := fiber.New(fiber.Config{
		ErrorHandler: response.ErrorHandler,
	})

	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(logger.New())
	// nanti bisa setting disini
	app.Use(cors.New(cors.Config{}))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Setyabudi Dwisandi Arifin")
	})

	routes.SetUpRouter(
		app.Group("/api"),
		provider.Deps{
			DB: pool,
			Tx: database.NewTransactor(pool),
		},
	)

	// Server jalan di goroutine terpisah
	serverErr := make(chan error, 1)
	go func() {
		if err := app.Listen(":3000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	log.Println("server jalan di :3000")

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

	log.Println("server ditutup, menutup koneksi database")
	return nil
}
