package routes

import (
	"golang-rest-api/internal/feature/address"
	"golang-rest-api/internal/feature/auth"
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
	address.Router{},
	auth.Router{},
}

func SetUpRouter(router fiber.Router, prov provider.Deps) {
	for _, r := range listRouter {
		r.Register(router, prov)
	}
}
