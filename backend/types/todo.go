package types

import (
	"fmt"
	"time"
)

type Todo struct {
	// ID
	ID int64 `json:"id,omitempty"`
	// The title of the Todo
	Title string `json:"title"`
	// The content of the Todo
	Content string `json:"content"`
	// PostgreSQL uses a RFC3339-format
	Created time.Time `json:"created" example:"2023-07-06T15:58:47.22Z"`
	// PostgreSQL uses a RFC3339-format
	Updated time.Time `json:"updated" example:"2023-07-06T15:58:47.22Z"`
	// User ID
	CreatedBy int `json:"createdBy"`
	// User ID
	UpdatedBy int `json:"updatedBy"`
	// This boolean determines if the todo has been completed
	Done bool `json:"done"`
} // @name Todo


type InsertTodoParams struct {
	// The title of the Todo
	Title string `json:"title" validate:"required"`
	// The content of the Todo
	Content string `json:"content" validate:"required"`
	// PostgreSQL uses a RFC3339-format
	Created time.Time `json:"-" example:"2023-07-06T15:58:47.22Z" validate:"required"`
	// User ID
	CreatedBy int `json:"createdBy" validate:"required"`
	// This boolean determines if the todo has been completed
	Done bool `json:"done" validate:"required"`
} // @name InsertTodoParams


type UpdateTodoParams struct {
	// The title of the Todo
	Title string `json:"title,omitempty" validate:"required"`
	// The content of the Todo
	Content string `json:"content,omitempty" validate:"required"`
	// PostgreSQL uses a RFC3339-format
	Created time.Time `json:"created,omitempty" example:"2023-07-06T15:58:47.22Z" validate:"required"`
	// PostgreSQL uses a RFC3339-format
	Updated time.Time `json:"updated,omitempty" example:"2023-07-06T15:58:47.22Z" validate:"required"`
	// User ID
	CreatedBy int `json:"createdBy,omitempty" validate:"required"`
	// User ID
	UpdatedBy int `json:"updatedBy,omitempty" validate:"required"`
	// This boolean determines if the todo has been completed
	Done bool `json:"done,omitempty" validate:"required"`
} // @name UpdateTodoParams


func NewTodoFromParams(params InsertTodoParams) *Todo {
	return &Todo{
		Title:     params.Title,
		Content:   params.Content,
		Created:   time.Now().UTC(),
		CreatedBy: params.CreatedBy,
		Done:      params.Done,
	}
}

func (t *Todo) Validate() error {
	if len(t.Title) < 3 {
		return fmt.Errorf("title needs to be at least 3 characters")
	}
	if len(t.Content) < 3 {
		return fmt.Errorf("content needs to be at least 3 characters")
	}

	return nil
}
