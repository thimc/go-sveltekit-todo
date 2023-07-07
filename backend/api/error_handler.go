package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thimc/go-backend/types"
)

func HandleError(c *fiber.Ctx, err error) error {
	if resp, ok := err.(types.ApiResponse); ok {
		return c.Status(resp.ErrorCode).JSON(resp)
	}
	resp := types.NewApiResponse(false, err.Error(), http.StatusInternalServerError)
	return c.Status(resp.ErrorCode).JSON(resp)
}
