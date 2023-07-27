package utils

import (
	"encoding/json"
	"net/http"

	"github.com/thimc/go-svelte-todo/backend/types"
)

// Marshals “data“ into JSON in the response body. Returns the “error“ of “(*json.Encoder).Encode(v any)“.
func WriteJSON(w http.ResponseWriter, data any) error {
	w.Header().Add("Content-Type", "application/json")

	if err, ok := data.(*types.APIError); ok {
		w.WriteHeader(err.StatusCode)
		return json.NewEncoder(w).Encode(*err)
	}

	return json.NewEncoder(w).Encode(&data)
}

// Marshals “data“ into JSON in the response body. Returns an “APIError“ if “(*json.Encoder).Encode(v any)“ fails.
func ResponseWriteJSON(w http.ResponseWriter, data any) *types.APIError {
	if err := WriteJSON(w, data); err != nil {
		return types.NewAPIError(false, err, http.StatusInternalServerError)
	}
	return nil
}
