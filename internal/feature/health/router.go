package health

import (
	"context"
	"golang-rest-api/internal/provider"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct{}

type handler struct {
	pool *pgxpool.Pool
}

func (Router) Name() string { return "health" }

func (Router) Register(router fiber.Router, prov provider.Deps) {
	h := &handler{pool: prov.DB}
	router.Get("/health", h.check)
	router.Get("/ready", h.ready)
}

func (h *handler) check(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *handler) ready(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if err := h.pool.Ping(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":   "unavailable",
			"database": "down",
		})
	}

	stat := h.pool.Stat()

	return c.JSON(fiber.Map{
		"status":   "ok",
		"database": "up",
		"pool": fiber.Map{
			"total":    stat.TotalConns(),
			"idle":     stat.IdleConns(),
			"acquired": stat.AcquiredConns(),
		},
	})
}
