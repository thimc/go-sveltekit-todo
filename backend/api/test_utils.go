package api

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/thimc/go-backend/store"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

// Bootstraps a temporary database connection that is used by the tests
func setup(t *testing.T) store.DatabaseStorer {
	store, err := store.NewPostgreStore(os.Getenv("PSQL_URI"))
	if err != nil {
		t.Fatal(err)
	}

	return store
}
