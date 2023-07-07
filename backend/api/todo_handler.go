package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thimc/go-backend/store"
	"github.com/thimc/go-backend/types"
)

type TodoHandler struct {
	store store.DatabaseStorer
}

func NewTodoHandler(databaseStore store.DatabaseStorer) *TodoHandler {
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
// @Success		200	{object}	[]types.Todo
// @Router		/api/v1/todos [get]
// @Security	ApiKeyAuth
func (h *TodoHandler) HandleGetTodos(c *fiber.Ctx) error {
	todos, err := h.store.GetTodos(c.Context())
	if err != nil {
		return c.JSON(types.NewApiResponse(false, err.Error(), http.StatusInternalServerError))
	}

	return c.JSON(map[string]any{"count": len(todos), "result": todos})
}

// @Summary		Get a todo by the ID.
// @Description	fetch one todo.
// @Tags		todos
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Produce		json
// @Success		200	{object}	types.Todo
// @Router		/api/v1/todos/{id} [get]
// @Security	ApiKeyAuth
func (h *TodoHandler) HandleGetTodoByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}
	todo, err := h.store.GetTodoByID(c.Context(), int64(id))
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusNotFound)
	}
	if todo == nil {
		return types.NewApiResponse(false, "invalid todo id", http.StatusNotFound)
	}

	return c.JSON(todo)
}

// @Summary		Create a todo.
// @Description	create a single todo.
// @Tags		todos
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		params	body	types.InsertTodoParams	true	"Todo metadata"
// @Produce		json
// @Success		200	{object}	types.Todo
// @Router		/api/v1/todos [post]
// @Security	ApiKeyAuth
func (h *TodoHandler) HandleInsertTodo(c *fiber.Ctx) error {
	var params types.InsertTodoParams
	if err := c.BodyParser(&params); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}
	todo := types.NewTodoFromParams(params)
	if err := todo.Validate(); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}
	insertedTodo, err := h.store.InsertTodo(c.Context(), todo)
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	return c.JSON(insertedTodo)
}

// @Summary		Replace a todo.
// @Description	replaces a todo.
// @Tags		todos
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Param		todo	body	types.UpdateTodoParams	true	"New todo data"
// @Produce		json
// @Success		200	{object}	nil
// @Security	ApiKeyAuth
// @Router		/api/v1/todos/{id} [put]
func (h *TodoHandler) HandlePutTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}
	var params types.UpdateTodoParams
	if err := c.BodyParser(&params); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	if err := h.store.UpdateTodoByID(c.Context(), params, int64(id)); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// @Summary		Delete a todo.
// @Description	deletes a todo by the ID.
// @Tags		todos
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Produce		json
// @Success		200	{object}	nil
// @Router		/api/v1/todos/{id} [delete]
// @Security	ApiKeyAuth
func (h *TodoHandler) HandleDeleteTodoByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	if err := h.store.DeleteTodoByID(c.Context(), int64(id)); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusNotFound)
	}

	return nil
}

// @Summary		Patch a todo.
// @Description	updates only the specified todo fields.
// @Tags		todos
// @Accept		json
// @Param		Authorization	header	string	true	"JWT Token, needs to start with Bearer"
// @Param		id	path	int	true	"Todo ID"
// @Param		todo	body	types.UpdateTodoParams	false	"New todo data"
// @Produce		json
// @Success		200	{object}	nil
// Security		ApiKeyAuth
// @Router		/api/v1/todos/{id} [patch]
func (h *TodoHandler) HandlePatchTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}
	_ = id

	var params types.UpdateTodoParams
	if err := c.BodyParser(&params); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusBadRequest)
	}

	if err := h.store.PatchTodoByID(c.Context(), int64(id), params); err != nil {
		return types.NewApiResponse(false, err.Error(), http.StatusNotFound)
	}

	return nil
}

