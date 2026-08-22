package address

import (
	"golang-rest-api/internal/response"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) FindById(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	u, err := h.usecase.FindByIdUser(c.Context(), id)
	if err != nil {
		return err
	}

	return response.OK(c, u)
}

func (h *Handler) Store(c fiber.Ctx) error {
	var input AddressRequest
	id, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}
	if err := c.Bind().Body(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.Create(c.Context(), id, input)
	if err != nil {
		return err
	}

	return response.Created(c, result)
}
