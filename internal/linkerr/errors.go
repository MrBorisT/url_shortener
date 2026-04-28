package linkerr

import "errors"

var (
	ErrLinkNotFound              = errors.New("link not found")
	ErrShortCodeTaken            = errors.New("short code already taken")
	ErrLinkDisabled              = errors.New("link is disabled")
	ErrCouldNotGenerateShortCode = errors.New("could not generate unique short code")
	ErrEmptyLinkID               = errors.New("link id is required")

	//validation
	ErrURLRequired      = errors.New("url is required")
	ErrURLInvalid       = errors.New("url is invalid")
	ErrURLInvalidScheme = errors.New("url scheme must be http or https")
)
