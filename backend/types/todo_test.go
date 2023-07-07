package types

import (
	"testing"
)

func TestTodoValidate(t *testing.T) {
	tests := []struct {
		name     string
		fields   Todo
		expected bool
	}{
		{name: "", fields: Todo{
			Title:   "",
			Content: "",
		}, expected: true},

		{name: "test", fields: Todo{
			Title:   "test",
			Content: "test",
		}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Todo{
				ID:        tt.fields.ID,
				Title:     tt.fields.Title,
				Content:   tt.fields.Content,
				Created:   tt.fields.Created,
				Updated:   tt.fields.Updated,
				CreatedBy: tt.fields.CreatedBy,
				UpdatedBy: tt.fields.UpdatedBy,
				Done:      tt.fields.Done,
			}
			if err := tr.Validate(); (err != nil) != tt.expected {
				t.Errorf("expected validation %v, got %v", tt.expected, err)
			}
		})
	}
}
