package repository

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("data tidak ditemukan")
	ErrDuplicateKey = errors.New("data sudah ada")
	ErrValidation   = errors.New("validasi gagal")

	ErrUnauthorized       = errors.New("tidak terautentikasi")
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrInvalidToken       = errors.New("token tidak valid")
	ErrTokenExpired       = errors.New("token sudah kedaluwarsa")
	ErrTokenRevoked       = errors.New("token sudah dicabut")
	ErrTooManyRequests    = errors.New("terlalu banyak percobaan, coba lagi nanti")
)

func Validationf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrValidation, fmt.Sprintf(format, args...))
}
