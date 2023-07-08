package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/thimc/go-backend/types"
)

func TestRegisterFail(t *testing.T) {
	suite := newTestSuite(t)
	defer teardown(suite, t)

	authHandler := NewAuthHandler(suite.userStore)

	app := fiber.New(fiber.Config{ErrorHandler: HandleError})
	app.Post("/", authHandler.HandleLogin)

	params := types.UserParams{
		Email:    "test@user.com",
		Password: "",
	}

	b, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user params: %s", err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected http code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestRegisterSuccess(t *testing.T) {
	suite := newTestSuite(t)
	defer teardown(suite, t)

	authHandler := NewAuthHandler(suite.userStore)

	app := fiber.New(fiber.Config{ErrorHandler: HandleError})
	app.Post("/", authHandler.HandleRegister)

	params := types.UserParams{
		Email:    "test@user.com",
		Password: "Password123456",
	}

	b, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user params: %s", err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected http code %d, got %d (%s)", http.StatusOK, resp.StatusCode, resp.Status)
	}

	err = suite.userStore.DeleteUserByID(context.TODO(), int64(user.ID))
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoginFail(t *testing.T) {
	suite := newTestSuite(t)
	defer teardown(suite, t)

	authHandler := NewAuthHandler(suite.userStore)

	app := fiber.New(fiber.Config{ErrorHandler: HandleError})
	app.Post("/", authHandler.HandleLogin)

	params := types.UserParams{
		Email:    "",
		Password: "",
	}

	b, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user params: %s", err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected http code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestLoginSuccess(t *testing.T) {
	suite := newTestSuite(t)
	defer teardown(suite, t)

	authHandler := NewAuthHandler(suite.userStore)

	app := fiber.New(fiber.Config{ErrorHandler: HandleError})
	app.Post("/register", authHandler.HandleRegister)
	app.Post("/login", authHandler.HandleLogin)

	params := types.UserParams{
		Email:    "test@user.com",
		Password: "Password123456",
	}

	registerJson, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user registration params: %s", err)
	}

	reqRegister := httptest.NewRequest("POST", "/register", bytes.NewReader(registerJson))
	reqRegister.Header.Add("Content-Type", "application/json")
	respRegister, err := app.Test(reqRegister)
	if err != nil {
		t.Fatal(err)
	}

	if respRegister.StatusCode != http.StatusOK {
		t.Errorf("expected http code %d, got %d", http.StatusOK, respRegister.StatusCode)
	}

	b, err := json.Marshal(&params)
	if err != nil {
		t.Fatalf("error when marshaling user params: %s", err)
	}

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var authResponse types.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	}

	if authResponse.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}
	if authResponse.ExpiresAt == 0 {
		t.Fatalf("expected a expiration UNIX timestamp")
	}
	if authResponse.Email != params.Email {
		t.Fatalf("expected user email %s, got %s", params.Email, authResponse.Email)
	}
	if authResponse.ID < 1 {
		t.Fatalf("expected a user id, got %d", authResponse.ID)
	}

	if err := suite.userStore.DeleteUserByID(context.TODO(), int64(authResponse.ID)); err != nil {
		t.Fatal(err)
	}
}
