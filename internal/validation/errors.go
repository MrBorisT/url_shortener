package validation

import "errors"

var (
	//validation
	ErrURLRequired      = errors.New("url is required")
	ErrURLInvalid       = errors.New("url is invalid")
	ErrURLInvalidScheme = errors.New("url scheme must be http or https")
	ErrURLMissingHost   = errors.New("URL must have non-empty host")
)
