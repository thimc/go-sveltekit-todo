package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/thimc/go-backend/store"
)

type testSuite struct {
	databaseStore *store.PostgreTodoStore
	userStore     *store.PostgreUserStore
}

func teardown(suite *testSuite, t *testing.T) {
	suite.databaseStore.Close()
}

func newTestSuite(t *testing.T) *testSuite {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("PSQL_TEST_HOST"),
		os.Getenv("PSQL_TEST_PORT"),
		os.Getenv("PSQL_TEST_USERNAME"),
		os.Getenv("PSQL_TEST_PASSWORD"),
		os.Getenv("PSQL_TEST_DATABASE"),
		os.Getenv("PSQL_TEST_SSL"))

	databaseStore, err := store.NewPostgreTodoStore(connStr)
	if err != nil {
		t.Fatal(err)
	}

	userStore, err := store.NewPostgreUserStore(databaseStore)
	if err != nil {
		t.Fatal(err)
	}
	return &testSuite{
		databaseStore: databaseStore,
		userStore:     userStore,
	}
}
