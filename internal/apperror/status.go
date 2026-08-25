package apperror

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"golang-rest-api/internal/repository"
)

func StatusFor(err error) int {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fiber.StatusNotFound
	case errors.Is(err, repository.ErrDuplicateKey):
		return fiber.StatusConflict
	case errors.Is(err, repository.ErrValidation):
		return fiber.StatusUnprocessableEntity
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		return fe.Code
	}

	return fiber.StatusInternalServerError
}
