package order

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending   = "pending"
	StatusPaid      = "paid"
	StatusCancelled = "cancelled"
)

type Order struct {
	ID        uuid.UUID   `json:"id"`
	UserId    uuid.UUID   `json:"user_id"`
	Total     int64       `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Items     []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Nama      string    `json:"nama"`
	Qty       int       `json:"qty"`
	Price     int64     `json:"price"`
	Subtotal  int64     `json:"subtotal"`
}

func (o *Order) HitungTotal() {
	var total int64
	for i := range o.Items {
		o.Items[i].Subtotal = o.Items[i].Price * int64(o.Items[i].Qty)
		total += o.Items[i].Subtotal
	}

	o.Total = total
}

func (o Order) BolehDibatalkan() bool {
	return o.Status == StatusPending
}
