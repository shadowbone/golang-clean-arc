package user

import (
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"        db:"id"`
	Nama  string    `json:"nama"      db:"name"`
	Email string    `json:"email"     db:"email"`
	Role  string    `json:"role"      db:"role"`
	repository.AuditTrail
}
