package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thimc/go-svelte-todo/backend/utils"
)

func TestHealthCheck(t *testing.T) {
	handler := http.HandlerFunc(utils.HandleAPIFunc(HandleHealthCheck))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected http status code %v got %v", http.StatusOK, status)
	}

    if rr.Body.String() != "OK" {
		t.Errorf("expected http body %v got %v", "OK", rr.Body.String())
    }
}
