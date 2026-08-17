package routes

import (
	"golang-rest-api/internal/feature/health"
	"golang-rest-api/internal/feature/product"
	"golang-rest-api/internal/feature/user"
	"golang-rest-api/internal/provider"

	"github.com/gofiber/fiber/v3"
)

var listRouter = []provider.AppServiceProvider{
	health.Router{},
	product.Router{},
	user.Router{},
}

func SetUpRouter(router fiber.Router, prov provider.Deps) {
	for _, r := range listRouter {
		r.Register(router, prov)
	}
}
