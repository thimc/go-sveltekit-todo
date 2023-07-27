package utils

import (
	"log"
	"net/http"

	"github.com/thimc/go-svelte-todo/backend/types"
)

// Takes a ``types.APIFunc`` that returns an ``types.APIError`` and marshals it if needed.
func HandleAPIFunc(route types.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := route(w, r); err != nil {
			if err := WriteJSON(w, err); err != nil {
				log.Printf("Internal Marshal types.APIError: %+v\n", err)
			}
		}
	}
}
