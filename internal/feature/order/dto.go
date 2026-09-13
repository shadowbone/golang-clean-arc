package order

import "github.com/google/uuid"

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,min=1,max=50,dive"`
}

type OrderItemRequest struct {
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	Qty       int       `json:"qty" validate:"required,gt=0,lte=1000"`
}
