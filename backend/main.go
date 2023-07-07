package main

import (
	"log"
	"net/http"
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
	"github.com/thimc/go-backend/types"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

// @title			Backend
// @description	Go backend API using Fiber and PostgreSQL
// @contact.name	Thim Cederlund
// @license.name	MIT
// @host			localhost:1111
// @BasePath		/
func main() {
	// error handler
	config := fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		if resp, ok := err.(types.ApiResponse); ok {
			return c.Status(resp.ErrorCode).JSON(resp)
		}
		resp := types.NewApiResponse(false, err.Error(), http.StatusInternalServerError)
		return c.Status(resp.ErrorCode).JSON(resp)
	}}

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
	databaseStore, err := store.NewPostgreStore(os.Getenv("PSQL_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer databaseStore.Close()

	userStore, err := store.NewPostgreUserStore(databaseStore)
	if err != nil {
		log.Fatal(err)
	}
	defer userStore.Close()

	// handlers
	todoHandler := api.NewTodoHandler(databaseStore)
	userHandler := api.NewUserHandler(userStore)
	app.Get("/swagger/*", swagger.HandlerDefault)

	// routes
	app.Get("/api/health", api.HealthCheck)
	app.Post("/api/login", userHandler.HandleLogin)

	route := app.Group("/api/v1", api.JWT(userStore))
	route.Get("/todos", todoHandler.HandleGetTodos)
	route.Post("/todos", todoHandler.HandleInsertTodo)
	route.Delete("/todos/:id", todoHandler.HandleDeleteTodoByID)
	route.Put("/todos/:id", todoHandler.HandleUpdateTodo)
	route.Patch("/todos", todoHandler.HandlePatchTodo)
	route.Get("/todos/:id", todoHandler.HandleGetTodoByID)

	// serve
	if err := app.Listen(os.Getenv("LISTEN_ADDRESS")); err != nil {
		log.Fatal(err)
	}
}
