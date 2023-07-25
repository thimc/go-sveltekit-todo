package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thimc/go-svelte-todo/backend/types"
	"github.com/thimc/go-svelte-todo/backend/utils"
)

func TestGetUsersSuccess(t *testing.T) {
	testSuite := newTestSuite(t)
	defer testSuite.Teardown(t)
	loginHandler := http.HandlerFunc(utils.HandleAPIFunc(testSuite.authHandler.HandleLogin))
	userHandler := http.HandlerFunc(utils.HandleAPIFunc(testSuite.userHandler.HandleGetUsers))

	params := types.UserParams{
		Email:    fmt.Sprintf("test%d@golangtest.com", rand.Intn(10000)),
		Password: fmt.Sprintf("secret-%d-password", rand.Intn(10000)),
	}

	user, err := types.NewUser(params.Email, params.Password)
	if err != nil {
		t.Fatalf("user validation failed: %v", err)
	}

	insertedUser, err := testSuite.userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("error when creating a mock user: %v", err)
	}

	body, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user params: %s", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	loginHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected http status code %v got %v", http.StatusOK, rr.Code)
	}

	var resp types.LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("error when marshaling login response: %s", err)
	}

	if resp.ExpiresAt == 0 {
		t.Fatalf("expected an expiration date for the jwt token, got %v", resp.ExpiresAt)
	}
	if resp.ID == 0 && resp.ID != insertedUser.ID {
		t.Fatalf("expected user id %v, got %v", insertedUser.ID, resp.ID)
	}
	if resp.Email == "" && resp.Email != insertedUser.Email {
		t.Fatalf("expected user email %v, got %v", insertedUser.Email, resp.Email)
	}
	if resp.Token == "" {
		t.Fatalf("expected a jwt token, got %v", resp.Token)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", resp.Token))
	rr = httptest.NewRecorder()

	userHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected http status code %v got %v", http.StatusOK, rr.Code)
	}

	var userResp []types.User
	if err := json.NewDecoder(rr.Body).Decode(&userResp); err != nil {
		t.Fatalf("error when marshaling login response: %s", err)
	}

	if len(userResp) < 1 {
		t.Fatalf("expected a non-zero user list, got %+vv", userResp)
	}

	err = testSuite.userStore.DeleteUserByID(context.TODO(), int64(user.ID))
	if err != nil {
		t.Fatal(err)
	}

	removedUser, err := testSuite.userStore.GetUserByID(context.TODO(), int64(insertedUser.ID))
	if err == nil {
		t.Fatalf("expected the mock user to be removed, got %+v", removedUser)
	}
}
