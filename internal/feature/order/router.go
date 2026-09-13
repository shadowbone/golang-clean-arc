package order

import (
	"golang-rest-api/internal/provider"

	"github.com/gofiber/fiber/v3"
)

type Router struct{}

func (Router) Name() string { return "order" }

func (Router) Register(router fiber.Router, prov provider.Deps) {
	handler := NewOrderHandler(
		NewOrderUseCase(
			NewOrderRepository(prov.DB),
			NewProductAdapter(prov.DB),
			prov.Tx,
		),
	)

	api := router.Group("/orders")
	api.Use(provider.RequireAuth(prov.Tokens))
	api.Get("/", handler.FetchAll)
	api.Get("/:id", handler.FetchById)
	api.Post("/", handler.Store)
	api.Patch("/:id/cancel", handler.Cancel)
}
