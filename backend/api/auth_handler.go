package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thimc/go-svelte-todo/backend/api/middleware"
	"github.com/thimc/go-svelte-todo/backend/store"
	"github.com/thimc/go-svelte-todo/backend/types"
	"github.com/thimc/go-svelte-todo/backend/utils"
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
// @Failure		400	{object}	types.APIError
// @Router		/api/register [post]
// @Security	ApiKeyAuth
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) *types.APIError {
	var params types.UserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	if err := params.Validate(); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	user, err := types.NewUser(params.Email, params.Password)
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	insertedUser, err := h.store.CreateUser(r.Context(), user)
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	return utils.ResponseWriteJSON(w, insertedUser)
}

// @Summary		Login a user.
// @Description	generates a JWT token that is valid for 6 hours.
// @Tags		auth
// @Accept		json
// @Param		params	body	types.UserParams	true	"User credentials."
// @Produce		json
// @Success		200	{object}	types.LoginResponse
// @Failure		400	{object}	types.APIError
// @Failure		401	{object}	types.APIError
// @Router		/api/login [post]
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) *types.APIError {
	var params types.UserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	if err := params.Validate(); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	user, err := h.store.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		return types.NewAPIError(false, err, http.StatusUnauthorized)
	}
	if err := types.ValidPassword(user.EncryptedPassword, params.Password); err != nil {
		return types.NewAPIError(false, fmt.Errorf("Access denied"), http.StatusUnauthorized)
	}

	claims, token, err := middleware.CreateJWT(user)
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	return utils.ResponseWriteJSON(w, types.LoginResponse{
		ID:        user.ID,
		Email:     user.Email,
		Token:     token,
		ExpiresAt: claims["expiresAt"].(int64),
	})
}

// @Summary		Verify token.
// @Description	verifies if the JWT token is valid and returns the user info.
// @Tags		auth
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Produce		json
// @Success		200	{object}	types.User
// @Failure		400	{object}	types.APIError
// @Router		/api/check [get]
func (h *AuthHandler) HandleVerifyToken(w http.ResponseWriter, r *http.Request) *types.APIError {
	user, ok := r.Context().Value("user").(*types.User)
	if !ok {
		return types.NewAPIError(false, fmt.Errorf("Invalid token"), http.StatusBadRequest)
	}
	return utils.ResponseWriteJSON(w, user)
}
