package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/thimc/go-backend/api"
	_ "github.com/thimc/go-backend/docs"
	"github.com/thimc/go-backend/store"
)

// @title			Backend
// @description	Go backend API using Fiber and PostgreSQL
// @contact.name	Thim Cederlund
// @license.name	MIT
// @host			localhost:1111
// @BasePath		/
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// error handler
	config := fiber.Config{
		AppName:      "Backend",
		ErrorHandler: api.HandleError,
	}

	app := fiber.New(config)

	// middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// database connection
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
	app.Get("/swagger/*", swagger.HandlerDefault)

	// routes
	app.Get("/api/health", api.HealthCheck)
	app.Post("/api/register", authHandler.HandleRegister)
	app.Post("/api/login", authHandler.HandleLogin)

	route := app.Group("/api/v1", api.JWT(userStore))
	route.Post("/todos", todoHandler.HandleInsertTodo)
	route.Get("/todos", todoHandler.HandleGetTodos)
	route.Put("/todos/:id", todoHandler.HandlePutTodo)
	route.Get("/todos/:id", todoHandler.HandleGetTodoByID)
	route.Patch("/todos/:id", todoHandler.HandlePatchTodo)
	route.Delete("/todos/:id", todoHandler.HandleDeleteTodoByID)

	route.Get("/users", userHandler.HandleGetUsers)
	route.Get("/users/:id", userHandler.HandleGetUserByID)

	// serve
	if err := app.Listen(os.Getenv("LISTEN_ADDRESS")); err != nil {
		log.Fatal(err)
	}
}
