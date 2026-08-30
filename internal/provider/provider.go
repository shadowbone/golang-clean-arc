package provider

import (
	"golang-rest-api/internal/repository"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// disini untuk mendaftaran depedensi yg nantinya akan dibutuhkan
type Deps struct {
	DB     *pgxpool.Pool
	Tx     repository.Transactor
	Log    *slog.Logger
	Tokens TokenVerifier
}

type AppServiceProvider interface {
	Name() string
	Register(router fiber.Router, d Deps)
}
