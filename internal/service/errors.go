package service

import "errors"

var (
	ErrCouldNotGenerateShortCode = errors.New("could not generate unique short code")
	ErrEmptyLinkID               = errors.New("link id is required")
	ErrLinkNotFound              = errors.New("link not found")
	ErrLinkDisabled              = errors.New("link is disabled")
)
