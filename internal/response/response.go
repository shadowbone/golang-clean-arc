package response

import (
	"golang-rest-api/internal/apperror"
	"golang-rest-api/internal/logger"
	"log/slog"

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
	code := apperror.StatusFor(err)

	msg := err.Error()
	if code == fiber.StatusInternalServerError {
		logger.FromContext(c.Context()).Error("unhandled error",
			slog.String("error", err.Error()),
		)
		msg = "terjadi kesalahan internal"
	}

	return Error(c, code, msg)
}
