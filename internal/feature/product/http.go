package product

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"golang-rest-api/internal/response"
)

type ProductHandler struct {
	usecase *ProductUseCase
}

func NewProductHandler(uc *ProductUseCase) *ProductHandler {
	return &ProductHandler{usecase: uc}
}

func (h *ProductHandler) FetchAll(c fiber.Ctx) error {
	items, err := h.usecase.GetAllProduct(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, items)
}

func (h *ProductHandler) FetchById(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	p, err := h.usecase.FindByIdProduct(c.Context(), id)
	if err != nil {
		return err
	}

	return response.OK(c, p)
}

func (h *ProductHandler) Store(c fiber.Ctx) error {
	var input Product

	if err := c.Bind().Body(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.CreateProduct(c.Context(), input)
	if err != nil {
		return err
	}

	return response.Created(c, result)
}
