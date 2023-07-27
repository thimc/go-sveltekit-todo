package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thimc/go-svelte-todo/backend/store"
	"github.com/thimc/go-svelte-todo/backend/types"
	"github.com/thimc/go-svelte-todo/backend/utils"
)

type TodoHandler struct {
	store store.TodoStorer
}

func NewTodoHandler(databaseStore store.TodoStorer) *TodoHandler {
	return &TodoHandler{
		store: databaseStore,
	}
}


// @Summary		Get all todos.
// @Description	fetch every todo available.
// @Tags		todos
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Accept		*/*
// @Produce		json
// @Success		200	{object}	types.TodoGetAllResponse
// @Failure		400	{object}	types.APIError
// @Failure		500	{object}	types.APIError
// @Router		/api/v1/todos [get]
// @Security	ApiKeyAuth
func (h *TodoHandler) HandleGetTodos(w http.ResponseWriter, r *http.Request) *types.APIError {
	todos, err := h.store.GetTodos(r.Context())
	if err != nil {
		return types.NewAPIError(false, err, http.StatusInternalServerError)
	}

	return utils.ResponseWriteJSON(w, types.NewTodoGetAllResponse(todos))
}

// @Summary		Get a todo by the ID.
// @Description	fetch one todo.
// @Tags		todos
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Produce		json
// @Success		200	{object}	types.Todo
// @Failure		400	{object}	types.APIError
// @Failure		404	{object}	types.APIError
// @Router		/api/v1/todos/{id} [get]
// @Security	ApiKeyAuth
func (h *TodoHandler) HandleGetTodoByID(w http.ResponseWriter, r *http.Request) *types.APIError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}
	todo, err := h.store.GetTodoByID(r.Context(), int64(id))
	if err != nil {
		return types.NewAPIError(false, err, http.StatusNotFound)
	}

	return utils.ResponseWriteJSON(w, todo)
}

// @Summary		Create a todo.
// @Description	create a single todo.
// @Tags		todos
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		params	body	types.InsertTodoParams	true	"Todo metadata"
// @Produce		json
// @Success		200	{object}	types.Todo
// @Failure		400	{object}	types.APIError
// @Router		/api/v1/todos [post]
// @Security	ApiKeyAuth
func (h *TodoHandler) HandleInsertTodo(w http.ResponseWriter, r *http.Request) *types.APIError {
	var params types.InsertTodoParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}
	todo := types.NewTodoFromParams(params)
	if err := todo.Validate(); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}
	insertedTodo, err := h.store.InsertTodo(r.Context(), todo)
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	return utils.ResponseWriteJSON(w, insertedTodo)
}

// @Summary		Replace a todo.
// @Description	replaces a todo.
// @Tags		todos
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Param		todo	body	types.UpdateTodoParams	true	"New todo data"
// @Produce		json
// @Success		200	{object}	types.APIError
// @Failure		400	{object}	types.APIError
// @Security	ApiKeyAuth
// @Router		/api/v1/todos/{id} [put]
func (h *TodoHandler) HandlePutTodo(w http.ResponseWriter, r *http.Request) *types.APIError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}
	var params types.UpdateTodoParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	if err := h.store.UpdateTodoByID(r.Context(), params, int64(id)); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	return utils.ResponseWriteJSON(w, types.NewAPIError(true, fmt.Errorf("todo ID: %d", id), http.StatusOK))
}

// @Summary		Delete a todo.
// @Description	deletes a todo by the ID.
// @Tags		todos
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Produce		json
// @Success		200	{object}	types.APIError
// @Failure		400	{object}	types.APIError
// @Failure		404	{object}	types.APIError
// @Router		/api/v1/todos/{id} [delete]
// @Security	ApiKeyAuth
func (h *TodoHandler) HandleDeleteTodoByID(w http.ResponseWriter, r *http.Request) *types.APIError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	if err := h.store.DeleteTodoByID(r.Context(), int64(id)); err != nil {
		return types.NewAPIError(false, err, http.StatusNotFound)
	}

	return utils.ResponseWriteJSON(w, types.NewAPIError(true, fmt.Errorf("todo ID: %d", id), http.StatusOK))
}

// @Summary		Patch a todo.
// @Description	mutates a todos properties.
// @Tags		todos
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Param		todo	body	types.UpdateTodoParams	false	"New todo data"
// @Produce		json
// @Success		200	{object}	types.APIError
// @Failure		400	{object}	types.APIError
// @Failure		404	{object}	types.APIError
// Security		ApiKeyAuth
// @Router		/api/v1/todos/{id} [patch]
func (h *TodoHandler) HandlePatchTodoByID(w http.ResponseWriter, r *http.Request) *types.APIError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	var params types.UpdateTodoParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.NewAPIError(false, err, http.StatusBadRequest)
	}

	if err := h.store.PatchTodoByID(r.Context(), int64(id), params); err != nil {
		return types.NewAPIError(false, err, http.StatusNotFound)
	}

	return utils.ResponseWriteJSON(w, types.NewAPIError(true, fmt.Errorf("todo ID: %d", id), http.StatusOK))
}
