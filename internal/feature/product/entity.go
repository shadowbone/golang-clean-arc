package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `json:"id"         db:"id"`
	Nama      string    `json:"nama"       db:"name"`
	Harga     int64     `json:"harga"      db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
