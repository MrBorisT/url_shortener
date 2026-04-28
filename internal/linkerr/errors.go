package linkerr

import "errors"

var (
	ErrLinkNotFound   = errors.New("link not found")
	ErrShortCodeTaken = errors.New("short code already taken")
	ErrLinkDisabled   = errors.New("link is disabled")
)
