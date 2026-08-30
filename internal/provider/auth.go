package provider

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"golang-rest-api/internal/repository"
)

type TokenVerifier interface {
	ParseAccess(token string) (uuid.UUID, error)
}

type userCtxKey struct{}

func ContextWithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userCtxKey{}, id)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(userCtxKey{}).(uuid.UUID)
	if !ok {
		return uuid.Nil, repository.ErrUnauthorized
	}
	return id, nil
}

func RequireAuth(tv TokenVerifier) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		if header == "" {
			return repository.ErrUnauthorized
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return repository.ErrInvalidToken
		}

		userID, err := tv.ParseAccess(parts[1])
		if err != nil {
			return err
		}

		c.SetContext(ContextWithUserID(c.Context(), userID))
		return c.Next()
	}
}
