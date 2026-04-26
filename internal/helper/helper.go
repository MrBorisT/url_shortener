package helper

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MrBorisT/url_shortener/internal/validation"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}

func WriteValidationError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, validation.ErrURLRequired):
		_ = WriteJSONError(w, http.StatusBadRequest, "original_url is required")
	case errors.Is(err, validation.ErrURLInvalid):
		_ = WriteJSONError(w, http.StatusBadRequest, "original_url must be a valid URL")
	case errors.Is(err, validation.ErrURLInvalidScheme):
		_ = WriteJSONError(w, http.StatusBadRequest, "original_url scheme must be http or https")
	case errors.Is(err, validation.ErrURLMissingHost):
		_ = WriteJSONError(w, http.StatusBadRequest, "original_url host is required")
	default:
		return false
	}

	return true
}
