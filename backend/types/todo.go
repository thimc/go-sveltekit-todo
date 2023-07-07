package types

import (
	"fmt"
	"time"
)

type Todo struct {
	// ID
	ID int64 `json:"id,omitempty" example:"0"`
	// The title of the Todo
	Title string `json:"title" example:"My title"`
	// The content of the Todo
	Content string `json:"content" example:"My content"`
	// PostgreSQL uses a ISO 8601-format
	Created time.Time `json:"created" example:"2006-01-02 15:04:05.000-07"`
	// PostgreSQL uses a ISO 8601-format
	Updated time.Time `json:"updated" example:"2006-01-02 15:04:05.000-07"`
	// User ID
	CreatedBy int `json:"createdBy" example:"0"`
	// User ID
	UpdatedBy int `json:"updatedBy" example:"0"`
	// This boolean determines if the todo has been completed
	Done bool `json:"done" example:"false"`
} // @name Todo

type InsertTodoParams struct {
	// The title of the Todo
	Title string `json:"title" example:"My new title" validate:"required"`
	// The content of the Todo
	Content string `json:"content" example:"My new content" validate:"required"`
	// PostgreSQL uses a ISO 8601-format
	Created time.Time `json:"-" example:"2006-01-02 15:04:05.000-07" validate:"required"`
	// User ID
	CreatedBy int `json:"createdBy" example:"0" validate:"required"`
	// This boolean determines if the todo has been completed
	Done bool `json:"done" example:"false" validate:"required"`
} // @name InsertTodoParams

type UpdateTodoParams struct {
	// The title of the Todo
	Title *string `json:"title,omitempty" sql:"title" example:"My new title" validate:"required"`
	// The content of the Todo
	Content *string `json:"content,omitempty" sql:"content" example:"My new content" validate:"required"`
	// PostgreSQL uses a ISO 8601-format
	Created *time.Time `json:"created,omitempty" sql:"created" example:"2006-01-02 15:04:05.000-07" validate:"required"`
	// PostgreSQL uses a ISO 8601-format
	Updated *time.Time `json:"updated,omitempty" sql:"updated" example:"2006-01-02 15:04:05.000-07" validate:"required"`
	// User ID
	CreatedBy *int `json:"createdBy,omitempty" sql:"created_by" example:"0" validate:"required"`
	// User ID
	UpdatedBy *int `json:"updatedBy,omitempty" sql:"updated_by" example:"0" validate:"required"`
	// This boolean determines if the todo has been completed
	Done *bool `json:"done,omitempty" sql:"done" example:"false" validate:"required"`
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
