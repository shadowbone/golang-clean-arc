package repository

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("data tidak ditemukan")
	ErrDuplicateKey = errors.New("data sudah ada")
	ErrValidation   = errors.New("validasi gagal")
)

func Validationf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrValidation, fmt.Sprintf(format, args...))
}
