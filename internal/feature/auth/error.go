package auth

import "golang-rest-api/internal/repository"

var (
	ErrInvalidCredentials = repository.ErrInvalidCredentials
	ErrInvalidToken       = repository.ErrInvalidToken
	ErrTokenExpired       = repository.ErrTokenExpired
	ErrTokenRevoked       = repository.ErrTokenRevoked
)
