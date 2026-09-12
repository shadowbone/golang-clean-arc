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

	loginLimit := provider.RateLimit(prov.Limiter, "login",
		prov.RateLimiter.LoginMax, prov.RateLimiter.LoginWindow)
	registerLimit := provider.RateLimit(prov.Limiter, "register",
		prov.RateLimiter.RegisterMax, prov.RateLimiter.RegisterWindow)

	login := api.Group("")
	login.Use(loginLimit)
	login.Post("/login", handler.Login)
	login.Post("/refresh", handler.Refresh)

	reg := api.Group("")
	reg.Use(registerLimit)
	reg.Post("/register", handler.Register)

	api.Post("/logout", handler.Logout)

	protected := api.Group("")
	protected.Use(provider.RequireAuth(prov.Tokens))
	protected.Get("/me", handler.Me)
	protected.Post("/logout-all", handler.LogoutAll)
}
