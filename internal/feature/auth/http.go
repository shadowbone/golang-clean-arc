package auth

import (
	"golang-rest-api/internal/provider"
	"golang-rest-api/internal/response"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	usecase *AuthUseCase
}

func NewAuthHandler(uc *AuthUseCase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tokens, err := h.usecase.Register(c.Context(), req)
	if err != nil {
		return err
	}

	return response.Created(c, tokens)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tokens, err := h.usecase.Login(c.Context(), req)
	if err != nil {
		return err
	}

	return response.OK(c, tokens)
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tokens, err := h.usecase.Refresh(c.Context(), req)
	if err != nil {
		return err
	}

	return response.OK(c, tokens)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.Logout(c.Context(), req); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) LogoutAll(c fiber.Ctx) error {
	userID, err := provider.UserIDFromContext(c.Context())
	if err != nil {
		return err
	}

	if err := h.usecase.LogoutAll(c.Context(), userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) Me(c fiber.Ctx) error {
	userID, err := provider.UserIDFromContext(c.Context())
	if err != nil {
		return err
	}

	return response.OK(c, fiber.Map{"user_id": userID})
}
