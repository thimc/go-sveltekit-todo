package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thimc/go-backend/store"
	"github.com/thimc/go-backend/types"
)

type UserHandler struct {
	store store.UserStorer
}

func NewUserHandler(store store.UserStorer) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// @Summary		Get all users.
// @Description	gets all users.
// @Tags		users
// @Accept		json
// @Param			Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Produce		json
// @Success		200	{object}	[]types.User
// @Router		/api/v1/users [get]
// @Security	ApiKeyAuth
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.GetUsers(c.Context())
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusNotFound)
	}

	return c.JSON(users)
}

// @Summary		Get user by ID.
// @Description	fetch one user.
// @Tags		users
// @Accept		json
// @Param			Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param			id	path	int	true	"Todo ID"
// @Produce		json
// @Success		200	{object}	types.User
// @Router		/api/v1/users/{id} [get]
// @Security	ApiKeyAuth
func (h *UserHandler) HandleGetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}
	user, err := h.store.GetUserByID(c.Context(), int64(id))
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusNotFound)
	}
	if user == nil {
		return types.NewApiResponse(false, "invalid user id", http.StatusNotFound)
	}

	return c.JSON(user)
}
