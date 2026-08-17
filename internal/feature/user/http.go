package user

import (
	"golang-rest-api/internal/response"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	usecase *UserUseCase
}

func NewUserHandler(uc *UserUseCase) *UserHandler {
	return &UserHandler{usecase: uc}

}

func (h *UserHandler) FetchAll(c fiber.Ctx) error {
	users, err := h.usecase.GetAllUser(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, users)
}

func (h *UserHandler) FetchById(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	u, err := h.usecase.FindByIdUser(c.Context(), id)
	if err != nil {
		return err
	}

	return response.OK(c, u)
}

func (h *UserHandler) Store(c fiber.Ctx) error {
	var input User

	if err := c.Bind().Body(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.CreateUser(c.Context(), input)
	if err != nil {
		return err
	}

	return response.Created(c, result)
}
