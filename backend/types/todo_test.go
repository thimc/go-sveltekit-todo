package types

import (
	"testing"
	"time"
)

func TestNewTodoValidate(t *testing.T) {
	params := InsertTodoParams{
		Title:     "My title",
		Content:   "My Content",
		Created:   time.Now(),
		CreatedBy: 1,
		Done:      true,
	}

	todo := NewTodoFromParams(params)
	if err := todo.Validate(); err != nil {
		t.Fatalf("expected the todo to be valid, got %s", err)
	}

	if todo.Title != params.Title {
		t.Errorf("expected title %s, got %s", params.Title, todo.Title)
	}
	if todo.Content != params.Content {
		t.Errorf("expected content %s, got %s", params.Content, todo.Content)
	}
	if todo.CreatedBy != params.CreatedBy {
		t.Errorf("expected created by ID %d, got %d", params.CreatedBy, todo.CreatedBy)
	}
}

//func TestTodoValidate(t *testing.T) {
//	tests := []struct {
//		name     string
//		fields   Todo
//		expected bool
//	}{
//		{name: "", fields: Todo{
//			Title:   "",
//			Content: "",
//		}, expected: true},
//
//		{name: "test", fields: Todo{
//			Title:   "test",
//			Content: "test",
//		}, expected: false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tr := &Todo{
//				ID:        tt.fields.ID,
//				Title:     tt.fields.Title,
//				Content:   tt.fields.Content,
//				Created:   tt.fields.Created,
//				Updated:   tt.fields.Updated,
//				CreatedBy: tt.fields.CreatedBy,
//				UpdatedBy: tt.fields.UpdatedBy,
//				Done:      tt.fields.Done,
//			}
//			if err := tr.Validate(); (err != nil) != tt.expected {
//				t.Errorf("expected validation %v, got %v", tt.expected, err)
//			}
//		})
//	}
//}
