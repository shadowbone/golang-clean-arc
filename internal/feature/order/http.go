package order

import (
	"golang-rest-api/internal/provider"
	"golang-rest-api/internal/repository"
	"golang-rest-api/internal/response"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type OrderHandler struct {
	usecase *OrderUseCase
}

func NewOrderHandler(uc *OrderUseCase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

func (h *OrderHandler) Store(c fiber.Ctx) error {
	userID, err := provider.UserIDFromContext(c.Context())
	if err != nil {
		return err
	}

	var req CreateOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.Create(c.Context(), userID, req)
	if err != nil {
		return err
	}

	return response.Created(c, result)
}

func (h *OrderHandler) FetchAll(c fiber.Ctx) error {
	userID, err := provider.UserIDFromContext(c.Context())
	if err != nil {
		return err
	}

	var p repository.Pagination
	if err := c.Bind().Query(&p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "parameter pagination tidak valid")
	}

	result, err := h.usecase.List(c.Context(), userID, p)
	if err != nil {
		return err
	}

	return response.OK(c, result)
}

func (h *OrderHandler) FetchById(c fiber.Ctx) error {
	userID, err := provider.UserIDFromContext(c.Context())
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	result, err := h.usecase.GetByID(c.Context(), userID, id)
	if err != nil {
		return nil
	}

	return response.OK(c, result)
}

func (h *OrderHandler) Cancel(c fiber.Ctx) error {
	userID, err := provider.UserIDFromContext(c.Context())
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id uuid harus berupa UUID yang valid")
	}

	if err := h.usecase.Cancel(c.Context(), userID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
