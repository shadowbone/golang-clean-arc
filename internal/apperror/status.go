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
	case errors.Is(err, repository.ErrUnauthorized),
		errors.Is(err, repository.ErrInvalidCredentials),
		errors.Is(err, repository.ErrInvalidToken),
		errors.Is(err, repository.ErrTokenExpired),
		errors.Is(err, repository.ErrTokenRevoked):
		return fiber.StatusUnauthorized
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		return fe.Code
	}

	return fiber.StatusInternalServerError
}
