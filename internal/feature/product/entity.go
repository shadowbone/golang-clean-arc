package product

import (
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
)

type Product struct {
	ID    uuid.UUID `json:"id"         db:"id"`
	Nama  string    `json:"nama"       db:"name"`
	Harga int64     `json:"harga"      db:"price"`
	repository.AuditTrail
}
