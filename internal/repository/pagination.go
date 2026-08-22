package repository

const (
	defaultLimit = 20
	maxLimit     = 100
)

type Pagination struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit < 1 {
		p.Limit = defaultLimit
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type Page[T any] struct {
	Items []T  `json:"items"`
	Meta  Meta `json:"meta"`
}

func NewPage[T any](items []T, total int64, p Pagination) Page[T] {
	totalPages := 0
	if p.Limit > 0 {
		totalPages = int((total + int64(p.Limit) - 1) / int64(p.Limit))
	}

	if items == nil {
		items = []T{}
	}

	return Page[T]{
		Items: items,
		Meta: Meta{
			Page:       p.Page,
			Limit:      p.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
