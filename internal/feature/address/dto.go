package address

type AddressRequest struct {
	Label  string `json:"label" validate:"required, min=2, max=30"`
	Street string `json:"street" validate:"required, min=5"`
	City   string `json:"city" validate:"required, min=2, max=80"`
}

func (r AddressRequest) ToEntity() Address {
	return Address{Label: r.Label, Street: r.Street, City: r.City}
}
