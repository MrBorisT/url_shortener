package validation

import (
	"net/url"
	"strings"

	"github.com/MrBorisT/url_shortener/internal/linkerr"
)

func NormalizeURL(raw string) (string, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return "", linkerr.ErrURLRequired
	}

	u, err := url.ParseRequestURI(normalized)
	if err != nil {
		return "", linkerr.ErrURLInvalid
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", linkerr.ErrURLInvalidScheme
	}

	if u.Host == "" {
		return "", linkerr.ErrURLMissingHost
	}

	return normalized, nil
}
