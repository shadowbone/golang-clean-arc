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

func (h *Handler) Update(c fiber.Ctx) error {
	var input AddressRequest
	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User id harus berupa UUID yang valid")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	if err := c.Bind().Body(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.Update(c.Context(), userId, id, input)
	if err != nil {
		return err
	}
	return response.OK(c, result)

}

func (h *Handler) Destroy(c fiber.Ctx) error {
	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User id harus berupa UUID yang valid")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	if err := h.usecase.Delete(c, userId, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) SetDefault(c fiber.Ctx) error {
	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User id harus berupa UUID yang valid")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	if err := h.usecase.SetDefault(c, userId, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
