package provider

import (
	"golang-rest-api/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// disini untuk mendaftaran depedensi yg nantinya akan dibutuhkan
type Deps struct {
	DB *pgxpool.Pool
	Tx repository.Transactor
}

type AppServiceProvider interface {
	Name() string
	Register(router fiber.Router, d Deps)
}
