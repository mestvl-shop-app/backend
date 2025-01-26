package service

import "errors"

var (
	ClientAlreadyExists      = errors.New("client already exists")
	ClientInvalidCredentials = errors.New("client invalid credentials")
)
