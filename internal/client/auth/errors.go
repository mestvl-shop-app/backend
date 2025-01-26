package auth

import "errors"

var (
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)
