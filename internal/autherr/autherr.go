package autherr

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrEmptyUserEmail     = errors.New("user email cannot be empty")
	ErrIncorrectUserEmail = errors.New("user email must be a valid email address")
	ErrEmptyUserPassword  = errors.New("user password cannot be empty")
	ErrShortUserPassword  = errors.New("user password must be at least 6 characters long")
	ErrLongUserPassword   = errors.New("user password cannot be longer than 72 characters")
)
