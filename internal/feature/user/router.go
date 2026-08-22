package user

import (
	"golang-rest-api/internal/provider"

	"github.com/gofiber/fiber/v3"
)

type Router struct{}

func (Router) Name() string { return "users" }

func (Router) Register(router fiber.Router, prov provider.Deps) {
	handler := NewUserHandler(
		NewUserUseCase(NewUserRepository(prov.DB)),
	)
	api := router.Group("/users")
	api.Get("/", handler.FetchAll)
	api.Get("/:id", handler.FetchById)
	api.Post("/", handler.Store)
	api.Put("/:id", handler.Update)
	api.Delete("/:id", handler.Destroy)
}
