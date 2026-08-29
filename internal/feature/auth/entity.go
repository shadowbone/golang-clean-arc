package auth

import (
	"time"

	"github.com/google/uuid"
)

type UserCredential struct {
	UserID            uuid.UUID `db:"user_id"`
	PasswordHash      string    `db:"password_hash"`
	PasswordChangedAt time.Time `db:"password_changed_at"`
}

type RefreshToken struct {
	ID         uuid.UUID  `db:"id"`
	UserID     uuid.UUID  `db:"user_id"`
	TokenHash  string     `db:"token_hash"`
	ExpiresAt  time.Time  `db:"expires_at"`
	RevokedAt  *time.Time `db:"revoked_at"`
	ReplacedBy *uuid.UUID `db:"replaced_by"`
	CreatedAt  time.Time  `db:"created_at"`
}

func (r RefreshToken) IsUsable() bool {
	return r.RevokedAt == nil && time.Now().Before(r.ExpiresAt)
}
