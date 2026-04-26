package storage

import "errors"

var (
	//user
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	//links
	ErrLinkNotFound     = errors.New("link not found")
	ErrEmptyOriginalURL = errors.New("empty original URL")
	ErrShortCodeTaken   = errors.New("short code already taken")
	ErrLinkDisabled     = errors.New("link is disabled")
)

const (
	PGCodeUniqueViolation = "23505"
)
