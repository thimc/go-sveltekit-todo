package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/thimc/go-svelte-todo/backend/api"
	"github.com/thimc/go-svelte-todo/backend/api/middleware"
	"github.com/thimc/go-svelte-todo/backend/store"
	"github.com/thimc/go-svelte-todo/backend/utils"

	swagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/thimc/go-svelte-todo/backend/docs"
)

// @title			Backend
// @description		Go backend API using Gorilla Mux and PostgreSQL
// @contact.name	Thim Cederlund
// @license.name	MIT
// @host			localhost:1234
// @BasePath		/
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	listenAddr := os.Getenv("LISTEN_ADDRESS")

	r := mux.NewRouter()
	r.Use(middleware.Logger)

	r.PathPrefix("/swagger/").Handler(swagger.Handler(
		swagger.DeepLinking(true),
		swagger.DocExpansion("none"),
		swagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
	log.Printf("Serving swagger docs on http://localhost%s/swagger/\n", listenAddr)

	// stores
	log.Printf("Connecting to the database..")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("PSQL_HOST"),
		os.Getenv("PSQL_PORT"),
		os.Getenv("PSQL_USERNAME"),
		os.Getenv("PSQL_PASSWORD"),
		os.Getenv("PSQL_DATABASE"),
		os.Getenv("PSQL_SSL"))
	databaseStore, err := store.NewPostgreTodoStore(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer databaseStore.Close()
	userStore, err := store.NewPostgreUserStore(databaseStore)
	if err != nil {
		log.Fatal(err)
	}

	// handlers
	todoHandler := api.NewTodoHandler(databaseStore)
	authHandler := api.NewAuthHandler(userStore)
	userHandler := api.NewUserHandler(userStore)

	// routes
	route := r.PathPrefix("/api").Subrouter()
	v1 := route.PathPrefix("/v1").Subrouter()

	// middleware
	jwt := middleware.NewJWTMiddleware(userStore)
	v1.Use(jwt.Middleware)

	route.HandleFunc("/health", utils.HandleAPIFunc(api.HandleHealthCheck)).Methods(http.MethodGet)
	route.HandleFunc("/register", utils.HandleAPIFunc(authHandler.HandleRegister)).Methods(http.MethodPost)
	route.HandleFunc("/login", utils.HandleAPIFunc(authHandler.HandleLogin)).Methods(http.MethodPost)

	// todo
	v1.HandleFunc("/todos", utils.HandleAPIFunc(todoHandler.HandleGetTodos)).Methods(http.MethodGet)
	v1.HandleFunc("/todos", utils.HandleAPIFunc(todoHandler.HandleInsertTodo)).Methods(http.MethodPost)
	v1.HandleFunc("/todos/{id}", utils.HandleAPIFunc(todoHandler.HandlePutTodo)).Methods(http.MethodPut)
	v1.HandleFunc("/todos/{id}", utils.HandleAPIFunc(todoHandler.HandleGetTodoByID)).Methods(http.MethodGet)
	v1.HandleFunc("/todos/{id}", utils.HandleAPIFunc(todoHandler.HandlePatchTodoByID)).Methods(http.MethodPatch)
	v1.HandleFunc("/todos/{id}", utils.HandleAPIFunc(todoHandler.HandleDeleteTodoByID)).Methods(http.MethodDelete)

	// users
	v1.HandleFunc("/users", utils.HandleAPIFunc(userHandler.HandleGetUsers)).Methods(http.MethodGet)
	v1.HandleFunc("/users/{id}", utils.HandleAPIFunc(userHandler.HandleGetUserByID)).Methods(http.MethodGet)

	v1.HandleFunc("/user/password", utils.HandleAPIFunc(userHandler.HandlePutUserPassword)).Methods(http.MethodPut)

	log.Printf("Serving on %s...", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, r))
}
