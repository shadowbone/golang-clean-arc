package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"        db:"id"`
	Nama      string    `json:"nama"      db:"name"`
	Email     string    `json:"email"     db:"email"`
	Role      string    `json:"role"      db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
