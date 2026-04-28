package validation

import (
	"net/url"
	"strings"
)

func NormalizeURL(raw string) (string, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return "", ErrURLRequired
	}

	u, err := url.ParseRequestURI(normalized)
	if err != nil {
		return "", ErrURLInvalid
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrURLInvalidScheme
	}

	return normalized, nil
}
