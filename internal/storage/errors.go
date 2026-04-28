package storage

import "errors"

var (
	//user
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

const (
	PGCodeUniqueViolation = "23505"
)
