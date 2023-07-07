package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thimc/go-backend/store"
	"github.com/thimc/go-backend/types"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store store.UserStorer
}

func NewUserHandler(store store.UserStorer) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

//	@Summary		Login.
//	@Description	generates a JWT token.
//	@Tags			auth
//	@Accept			json
//	@Param			params	body	types.UserParams	true	"User credentials."
//	@Produce		json
//	@Success		200	{object}	types.LoginResponse
//	@Router			/api/login [post]
func (h *UserHandler) HandleLogin(ctx *fiber.Ctx) error {
	var params types.UserParams
	if err := ctx.BodyParser(&params); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	user, err := h.store.GetUserByEmail(ctx.Context(), params.Email)
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password))
	if err != nil {
		return types.NewApiResponse(false, "access denied", http.StatusUnauthorized)
	}

	claims, token, err := createJWT(user)
	if err != nil {
		return types.NewApiResponse(false, "internal server error", http.StatusInternalServerError)
	}

	return ctx.JSON(types.LoginResponse{
		ID:        user.ID,
		Email:     user.Email,
		Token:     token,
		ExpiresAt: claims["expires"].(int64),
	})
}
