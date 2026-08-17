package response

import (
	"errors"
	"golang-rest-api/internal/repository"
	"log"

	"github.com/gofiber/fiber/v3"
)

func OK(c fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{"status": "sukses", "data": data})
}

func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "sukses dibuat",
		"data":   data,
	})
}

func Error(c fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{"status": "gagal", "error": msg})
}

func ErrorHandler(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return Error(c, fiber.StatusNotFound, err.Error())

	case errors.Is(err, repository.ErrDuplicateKey):
		return Error(c, fiber.StatusConflict, err.Error())

	case errors.Is(err, repository.ErrValidation):
		return Error(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		return Error(c, fe.Code, fe.Message)
	}

	// Error tak terduga: log detailnya, sembunyikan dari client
	log.Printf("[ERROR] %s %s: %v", c.Method(), c.Path(), err)
	return Error(c, fiber.StatusInternalServerError, "terjadi kesalahan internal")
}
