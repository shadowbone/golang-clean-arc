package user

import (
	"golang-rest-api/internal/repository"
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
	p := repository.Pagination{
		Page:  fiber.Query(c, "page", 0),
		Limit: fiber.Query(c, "limit", 0),
	}
	users, err := h.usecase.GetAllUser(c.Context(), p)
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

func (h *UserHandler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	var input User
	if err := c.Bind().Body(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.UpdateUser(c.Context(), id, input)
	if err != nil {
		return err
	}

	return response.OK(c, result)
}

func (h *UserHandler) Destroy(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id harus berupa UUID yang valid")
	}

	if err := h.usecase.DeleteUser(c.Context(), id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
