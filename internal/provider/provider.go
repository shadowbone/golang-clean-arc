package provider

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// disini untuk mendaftaran depedensi yg nantinya akan dibutuhkan
type Deps struct {
	DB *pgxpool.Pool
}

type AppServiceProvider interface {
	Name() string
	Register(router fiber.Router, d Deps)
}
