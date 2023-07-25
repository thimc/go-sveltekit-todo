package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thimc/go-svelte-todo/backend/store"
	"github.com/thimc/go-svelte-todo/backend/types"
	"github.com/thimc/go-svelte-todo/backend/utils"
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
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Produce		json
// @Success		200	{object}	[]types.User
// @Router		/api/v1/users [get]
// @Security	ApiKeyAuth
func (h *UserHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) *types.APIError {
	users, err := h.store.GetUsers(r.Context())
	if err != nil {
		return types.NewAPIError(false, err, http.StatusNotFound)
	}

	return utils.ResponseWriteJSON(w, users)
}

// @Summary		Get user by ID.
// @Description	fetch one user.
// @Tags		users
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Produce		json
// @Success		200	{object}	types.User
// @Router		/api/v1/users/{id} [get]
// @Security	ApiKeyAuth
func (h *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) *types.APIError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}
	user, err := h.store.GetUserByID(r.Context(), int64(id))
	if err != nil {
		return types.NewAPIError(false, err, http.StatusNotFound)
	}

	return utils.ResponseWriteJSON(w, user)
}

// @Summary		Update the password.
// @Description	updates the users password.
// @Tags		users
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		params	body	types.UserPutPasswordParams	true	"New user credentials"
// @Produce		json
// @Success		200	{object}	nil
// @Router		/api/v1/user/password [put]
// @Security	ApiKeyAuth
func (h *UserHandler) HandlePutUserPassword(w http.ResponseWriter, r *http.Request) *types.APIError {
	if user, ok := r.Context().Value("user").(*types.User); ok {
		log.Printf("Hello %s\n", user.Email)

		var params *types.UserPutPasswordParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			return types.NewAPIError(false, err, http.StatusBadRequest)
		}

		if err := params.Validate(); err != nil {
			return types.NewAPIError(false, err, http.StatusBadRequest)
		}

		newUser, err := types.NewUser(user.Email, params.Password)
		if err != nil {
			return types.NewAPIError(false, err, http.StatusBadRequest)
		}

		err = h.store.UpdateUserPasswordByID(r.Context(), newUser.EncryptedPassword, int64(user.ID))
		if err != nil {
			return types.NewAPIError(false, err, http.StatusBadRequest)
		}

		return nil
	}

	return utils.ResponseWriteJSON(w, types.NewAPIError(false, fmt.Errorf("unknown user"), http.StatusBadRequest))
}
