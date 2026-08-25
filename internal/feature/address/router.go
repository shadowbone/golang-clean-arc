package address

import (
	"golang-rest-api/internal/provider"

	"github.com/gofiber/fiber/v3"
)

type Router struct{}

func (Router) Name() string { return "addresses" }

func (Router) Register(router fiber.Router, prov provider.Deps) {
	handler := NewHandler(
		NewUseCase(
			NewAddressRepository(prov.DB),
			prov.Tx,
		),
	)
	api := router.Group("/users")
	api.Get("/:userId/addresses", handler.FindById)
	api.Post("/:userId/addresses", handler.Store)
	api.Put("/:userId/addresses/:id", handler.Update)
	api.Delete("/:userId/addresses/:id", handler.Destroy)
	api.Patch("/:userId/addresses/:id/default", handler.SetDefault)
}
