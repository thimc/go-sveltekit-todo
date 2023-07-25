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

func TestRegisterFail(t *testing.T) {
	testSuite := newTestSuite(t)
	defer testSuite.Teardown(t)
	handler := http.HandlerFunc(utils.HandleAPIFunc(testSuite.authHandler.HandleRegister))

	tests := []struct {
		name                 string
		params               types.UserParams
		httpStatusCode       int
		expectValidationFail bool
	}{
		{
			name: "Empty Email",
			params: types.UserParams{
				Email:    "",
				Password: "test-password",
			},
			httpStatusCode:       http.StatusBadRequest,
			expectValidationFail: true,
		},
		{
			name: "Empty Password",
			params: types.UserParams{
				Email: "user@domain.com",
			},
			httpStatusCode:       http.StatusBadRequest,
			expectValidationFail: true,
		},
		{
			name: "Invalid Email",
			params: types.UserParams{
				Email: "user@domaincom",
			},
			httpStatusCode:       http.StatusBadRequest,
			expectValidationFail: true,
		},
		{
			name: "Invalid Password",
			params: types.UserParams{
				Email:    "user@domain.com",
				Password: "pass",
			},
			httpStatusCode:       http.StatusBadRequest,
			expectValidationFail: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(&tt.params)
			if err != nil {
				t.Fatalf("error when marshaling user params: %s", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.httpStatusCode {
				t.Errorf("expected http status code %v got %v", tt.httpStatusCode, status)
			}

			if err := tt.params.Validate(); (err == nil) == tt.expectValidationFail {
				t.Error(err)
			}
		})
	}
}

func TestRegisterSuccess(t *testing.T) {
	testSuite := newTestSuite(t)
	defer testSuite.Teardown(t)
	handler := http.HandlerFunc(utils.HandleAPIFunc(testSuite.authHandler.HandleRegister))

	params := types.UserParams{
		Email:    fmt.Sprintf("test%d@golangtest.com", rand.Intn(10000)),
		Password: fmt.Sprintf("secret-%d-password", rand.Intn(10000)),
	}

	if err := params.Validate(); err != nil {
		t.Fatal(err)
	}

	body, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user params: %s", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected http status code %v got %v (resp: %s)", http.StatusOK, status, rr.Body.String())
	}

	var user types.User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Fatalf("error when marshaling user: %s", err)
	}

	if user.ID == 0 {
		t.Errorf("expected a non-zero user ID, got %d", user.ID)
	}
	if user.Email == "" {
		t.Errorf("expected a non-empty user email, got %s", user.Email)
	}
	if user.EncryptedPassword != "" {
		t.Errorf("expected an empty encrypted password field, got %s", user.EncryptedPassword)
	}

	err = testSuite.userStore.DeleteUserByID(context.TODO(), int64(user.ID))
	if err != nil {
		t.Fatalf("error when removing mock user: %s", err)
	}

	removedUser, err := testSuite.userStore.GetUserByID(context.TODO(), int64(user.ID))
	if err == nil {
		t.Errorf("expected the mock user to be removed, got %+v", removedUser)
	}
}

func TestLoginFailUnknownUser(t *testing.T) {
	testSuite := newTestSuite(t)
	defer testSuite.Teardown(t)

	handler := http.HandlerFunc(utils.HandleAPIFunc(testSuite.authHandler.HandleLogin))

	params := types.UserParams{
		Email:    "unknown@unknown.com",
		Password: "blank",
	}

	if err := params.Validate(); err != nil {
		t.Fatal(err)
	}

	body, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user params: %s", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected http status code %v got %v", http.StatusUnauthorized, rr.Code)
	}
}

func TestLoginSuccess(t *testing.T) {
	testSuite := newTestSuite(t)
	defer testSuite.Teardown(t)
	handler := http.HandlerFunc(utils.HandleAPIFunc(testSuite.authHandler.HandleLogin))

	params := types.UserParams{
		Email:    fmt.Sprintf("test%d@golangtest.com", rand.Intn(10000)),
		Password: fmt.Sprintf("secret-%d-password", rand.Intn(10000)),
	}
	if err := params.Validate(); err != nil {
		t.Fatal(err)
	}

	user, err := types.NewUser(params.Email, params.Password)
	if err != nil {
		t.Fatalf("user validation failed: %v", err)
	}

	insertedUser, err := testSuite.userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("error when creating a mock user: %v", err)
	}
	defer func() {
		err = testSuite.userStore.DeleteUserByID(context.TODO(), int64(insertedUser.ID))
		if err != nil {
			t.Fatalf("expected the mock user to be removed, got %+v", insertedUser)
		}
	}()

	body, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user params: %s", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected http status code %v, got %v", http.StatusOK, rr.Code)
	}

	var resp types.LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("error when marshaling login response: %s", err)
	}

	if resp.ExpiresAt == 0 {
		t.Errorf("expected an expiration date for the jwt token, got %v", resp.ExpiresAt)
	}
	if resp.ID == 0 && resp.ID != insertedUser.ID {
		t.Fatalf("expected user id %v, got %v", insertedUser.ID, resp.ID)
	}
	if resp.Email == "" && resp.Email != insertedUser.Email {
		t.Errorf("expected user email %v, got %v", insertedUser.Email, resp.Email)
	}
	if resp.Token == "" {
		t.Errorf("expected a jwt token, got %v", resp.Token)
	}

	err = testSuite.userStore.DeleteUserByID(context.TODO(), int64(insertedUser.ID))
	if err != nil {
		t.Fatalf("error when removing mock user: %s", err)
	}
}
