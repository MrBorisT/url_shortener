package storage

import "errors"

var (
	//user
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	//links
	ErrLinkNotFound = errors.New("link not found")
	ErrEmptyOriginalURL = errors.New("empty original URL")
)

const (
	PGCodeUniqueViolation = "23505"
)
