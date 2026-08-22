package address

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Label     string    `json:"label" db:"label"`
	Street    string    `json:"street" db:"street"`
	City      string    `json:"city" db:"city"`
	IsDefault bool      `json:"is_default" db:"is_default"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
