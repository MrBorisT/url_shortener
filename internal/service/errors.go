package service

import "errors"

var (
	ErrInvalidOriginalURL        = errors.New("invalid original url")
	ErrCouldNotGenerateShortCode = errors.New("could not generate unique short code")
	ErrEmptyLinkID               = errors.New("link id is required")
	ErrEmptyOriginalURL          = errors.New("original URL is required")
	ErrLinkNotFound              = errors.New("link not found")
	ErrLinkDisabled              = errors.New("link is disabled")
)
