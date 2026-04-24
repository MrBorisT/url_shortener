package helper

import (
	"encoding/json"
	"net/http"
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
