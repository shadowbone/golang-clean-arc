package product

import (
	"golang-rest-api/internal/provider"

	"github.com/gofiber/fiber/v3"
)

type Router struct{}

func (Router) Name() string { return "products" }

func (Router) Register(router fiber.Router, prov provider.Deps) {
	handler := NewProductHandler(
		NewProductUseCase(
			NewProductRepository(prov.DB),
		),
	)
	api := router.Group("/products")
	api.Use(provider.RequireAuth(prov.Tokens))
	api.Get("/", handler.FetchAll)
	api.Get("/:id", handler.FetchById)
	api.Post("/", handler.Store)
	api.Put("/:id", handler.Update)
	api.Delete("/:id", handler.Destroy)
}
