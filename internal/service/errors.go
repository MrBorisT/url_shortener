package service

import "errors"

var (
	ErrInvalidOriginalURL        = errors.New("invalid original url")
	ErrCouldNotGenerateShortCode = errors.New("could not generate unique short code")
)
