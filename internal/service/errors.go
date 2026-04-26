package service

import "errors"

var (
	//links
	ErrCouldNotGenerateShortCode = errors.New("could not generate unique short code")
	ErrEmptyLinkID               = errors.New("link id is required")
	ErrLinkNotFound              = errors.New("link not found")
	ErrLinkDisabled              = errors.New("link is disabled")

	//user
	ErrEmptyUserEmail     = errors.New("user email cannot be empty")
	ErrIncorrectUserEmail = errors.New("user email must be a valid email address")
	ErrEmptyUserPassword  = errors.New("user password cannot be empty")
	ErrShortUserPassword  = errors.New("user password must be at least 6 characters long")
	ErrLongUserPassword   = errors.New("user password cannot be longer than 72 characters")

	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
