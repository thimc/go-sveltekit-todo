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

	"github.com/gorilla/mux"
	"github.com/thimc/go-svelte-todo/backend/api/middleware"
	"github.com/thimc/go-svelte-todo/backend/types"
	"github.com/thimc/go-svelte-todo/backend/utils"
)

func TestHandleGetTodosFail(t *testing.T) {
	testSuite := newTestSuite(t)
	defer testSuite.Teardown(t)

	jwt := middleware.NewJWTMiddleware(testSuite.userStore)
	r := mux.NewRouter()
	r.Use(jwt.Middleware)

	r.HandleFunc("/", utils.HandleAPIFunc(testSuite.todoHandler.HandleGetTodos))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected http status code %v, got %v", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleGetTodosSuccess(t *testing.T) {
	testSuite := newTestSuite(t)
	defer testSuite.Teardown(t)

	userParams := types.UserParams{
		Email:    fmt.Sprintf("test%d@golangtest.com", rand.Intn(10000)),
		Password: fmt.Sprintf("secret-%d-password", rand.Intn(10000)),
	}

	if err := userParams.Validate(); err != nil {
		t.Fatal(err)
	}

	user, err := types.NewUser(userParams.Email, userParams.Password)
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

	todo := types.NewTodoFromParams(types.InsertTodoParams{
		Title:     "This is the title",
		Content:   "This is the content",
		CreatedBy: insertedUser.ID,
		Done:      false,
	})
	insertedTodo, err := testSuite.databaseStore.InsertTodo(context.TODO(), todo)
	defer func() {
		err := testSuite.databaseStore.DeleteTodoByID(context.TODO(), int64(insertedTodo.ID))
		if err != nil {
			t.Fatal(err)
		}
	}()

	jwt := middleware.NewJWTMiddleware(testSuite.userStore)
	r := mux.NewRouter()
	r.HandleFunc("/", utils.HandleAPIFunc(testSuite.authHandler.HandleLogin)).Methods(http.MethodPost)

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(jwt.Middleware)
	protected.HandleFunc("/get", utils.HandleAPIFunc(testSuite.todoHandler.HandleGetTodos)).Methods(http.MethodGet)

	loginBody, err := json.Marshal(map[string]any{
		"email":    userParams.Email,
		"password": userParams.Password,
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(loginBody))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected http status code %v, got %v", http.StatusOK, rr.Code)
	}

	var loginResp types.LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&loginResp); err != nil {
		t.Fatalf("expected a login response, got %v", err)
	}

	if loginResp.ID == 0 {
		t.Fatalf("expected a non-zero user ID in the authentication response")
	}
	if loginResp.Email == "" {
		t.Fatalf("expected a non-empty e-mail in the authentication response")
	}
	if loginResp.Token == "" {
		t.Fatalf("expected a non-empty token in the authentication response")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/get", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", loginResp.Token))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var getResp types.TodoGetAllResponse
	if err := json.NewDecoder(rr.Body).Decode(&getResp); err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected http status code %v, got %v", http.StatusOK, rr.Code)
	}

	if getResp.Count < 1 {
		t.Fatalf("expected todo item count above zero, got %v", getResp.Count)
	}

	gotTodo := getResp.Result[getResp.Count - 1]
	if gotTodo.Title != insertedTodo.Title {
		t.Fatalf("expected todo title '%v', got '%v'", insertedTodo.Title, gotTodo.Title)
	}
	if gotTodo.Content != insertedTodo.Content {
		t.Fatalf("expected todo title '%v', got '%v'", insertedTodo.Content, gotTodo.Content)
	}
	if gotTodo.Created.Unix() == 0 {
		t.Fatalf("expected a creation date, got '%v'", gotTodo.Created)
	}
	if gotTodo.Done != insertedTodo.Done {
		t.Fatalf("expected todo done status '%v', got '%v'", insertedTodo.Done, gotTodo.Done)
	}
}
