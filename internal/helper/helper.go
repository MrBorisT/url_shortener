package helper

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MrBorisT/url_shortener/internal/linkerr"
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
	case errors.Is(err, linkerr.ErrURLRequired):
		_ = WriteJSONError(w, http.StatusBadRequest, "original_url is required")
	case errors.Is(err, linkerr.ErrURLInvalid):
		_ = WriteJSONError(w, http.StatusBadRequest, "original_url must be a valid URL")
	case errors.Is(err, linkerr.ErrURLInvalidScheme):
		_ = WriteJSONError(w, http.StatusBadRequest, "original_url scheme must be http or https")
	default:
		return false
	}

	return true
}
