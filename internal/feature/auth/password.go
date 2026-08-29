package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const dummyHash = "$2a$12$C6UzMDM.H6dfI/f/IKcEeO1yFMPfLRPuBFEXjS9J8xVXWvBZjNQ2."

type Hasher struct {
	cost int
}

func NewHasher(cost int) *Hasher {
	return &Hasher{cost: cost}
}

func (h *Hasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(b), nil
}

func (h *Hasher) Verify(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (h *Hasher) VerifyDummy(password string) {
	_ = bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(password)) == nil
}
