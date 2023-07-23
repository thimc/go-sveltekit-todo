package api

import (
	"net/http"

	"github.com/thimc/go-svelte-todo/backend/types"
)

// @Summary		Show the status of server.
// @Description	get the status of server.
// @Tags		misc
// @Accept		*/*
// @Produce		plain
// @Success		200	"OK"
// @Router		/api/health [get]
func HandleHealthCheck(w http.ResponseWriter, r* http.Request) *types.APIError {
	w.Write([]byte("OK"))
	return nil
}
