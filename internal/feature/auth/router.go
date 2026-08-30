package auth

import (
	"github.com/gofiber/fiber/v3"

	"golang-rest-api/internal/feature/user"
	"golang-rest-api/internal/provider"
)

type Router struct{}

func (Router) Name() string { return "auth" }

func (Router) Register(router fiber.Router, prov provider.Deps) {
	tm, ok := prov.Tokens.(*TokenManager)
	if !ok {
		panic("Tokens harus *auth.TokenManager")
	}

	handler := NewAuthHandler(
		NewAuthUseCase(
			NewAuthRepository(prov.DB),
			NewUserAdapter(user.NewUserRepository(prov.DB)),
			tm,
			NewHasher(12),
			prov.Tx,
		),
	)

	api := router.Group("/auth")
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)
	api.Post("/refresh", handler.Refresh)
	api.Post("/logout", handler.Logout)

	protected := api.Group("")
	protected.Use(provider.RequireAuth(prov.Tokens))
	protected.Get("/me", handler.Me)
	protected.Post("/logout-all", handler.LogoutAll)
}
