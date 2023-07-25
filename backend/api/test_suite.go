package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/thimc/go-svelte-todo/backend/store"
)

type testSuite struct {
	databaseStore store.TodoStorer
	userStore     store.UserStorer

	authHandler *AuthHandler
	userHandler *UserHandler
	todoHandler *TodoHandler
}

func (s *testSuite) Teardown(t *testing.T) error {
	return s.databaseStore.Close()
}

func newTestSuite(t *testing.T) *testSuite {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("PSQL_HOST"),
		os.Getenv("PSQL_PORT"),
		os.Getenv("PSQL_USERNAME"),
		os.Getenv("PSQL_PASSWORD"),
		os.Getenv("PSQL_DATABASE"),
		os.Getenv("PSQL_SSL"))

	databaseStore, err := store.NewPostgreTodoStore(connStr)
	if err != nil {
		t.Fatal(err)
	}

	userStore, err := store.NewPostgreUserStore(databaseStore)
	if err != nil {
		t.Fatal(err)
	}

	authHandler := NewAuthHandler(userStore)
	userHandler := NewUserHandler(userStore)
	todoHandler := NewTodoHandler(databaseStore)

	return &testSuite{
		databaseStore: databaseStore,
		userStore:     userStore,
		authHandler:   authHandler,
		userHandler:   userHandler,
		todoHandler:   todoHandler,
	}
}
