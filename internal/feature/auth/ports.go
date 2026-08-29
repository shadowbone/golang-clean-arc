package auth

import (
	"context"

	"github.com/google/uuid"
)

type UserInfo struct {
	ID    uuid.UUID
	Nama  string
	Email string
}

type UserProvider interface {
	FindByEmail(ctx context.Context, email string) (UserInfo, error)
	Create(ctx context.Context, nama, email string) (UserInfo, error)
}
