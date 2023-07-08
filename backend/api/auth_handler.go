package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thimc/go-backend/store"
	"github.com/thimc/go-backend/types"
)

type AuthHandler struct {
	store store.UserStorer
}

func NewAuthHandler(store store.UserStorer) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

// @Summary		Register a user.
// @Description	register a regular user.
// @Tags		auth
// @Accept		json
// @Param		params	body	types.UserParams	true	"User credentials"
// @Produce		json
// @Success		200	{object}	types.User
// @Router		/api/register [post]
// @Security	ApiKeyAuth
func (h *AuthHandler) HandleRegister(ctx *fiber.Ctx) error {
	var params types.UserParams
	if err := ctx.BodyParser(&params); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	if err := params.Validate(); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	user, err := types.NewUser(params.Email, params.Password)
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	insertedUser, err := h.store.CreateUser(ctx.Context(), user)
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	return ctx.JSON(insertedUser)
}

// @Summary		Login a user.
// @Description	generates a JWT token that is valid for 6 hours.
// @Tags		auth
// @Accept		json
// @Param		params	body	types.UserParams	true	"User credentials."
// @Produce		json
// @Success		200	{object}	types.LoginResponse
// @Router		/api/login [post]
func (h *AuthHandler) HandleLogin(ctx *fiber.Ctx) error {
	var params types.UserParams
	if err := ctx.BodyParser(&params); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}
	if err := params.Validate(); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusUnauthorized)
	}

	user, err := h.store.GetUserByEmail(ctx.Context(), params.Email)
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusUnauthorized)
	}
	if err := types.ValidPassword(user.EncryptedPassword, params.Password); err != nil {
		return types.NewApiResponse(false, "access denied", http.StatusUnauthorized)
	}

	claims, token, err := createJWT(user)
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusUnauthorized)
	}

	return ctx.JSON(types.LoginResponse{
		ID:        user.ID,
		Email:     user.Email,
		Token:     token,
		ExpiresAt: claims["expires"].(int64),
	})
}
