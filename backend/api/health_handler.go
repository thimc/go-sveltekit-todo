package api

import "github.com/gofiber/fiber/v2"

//	@Summary		Show the status of server.
//	@Description	get the status of server.
//	@Tags			health
//	@Accept			*/*
//	@Produce		plain
//	@Success		200	"OK"
//	@Router			/api/health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
