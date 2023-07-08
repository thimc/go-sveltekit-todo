package types

import (
	"testing"
)

func TestNewUserValidation(t *testing.T) {
	params := UserParams{
		Email:    "test@test.org",
		Password: "test-password",
	}
	user, err := NewUser(params.Email, params.Password)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if err := ValidPassword(user.EncryptedPassword, params.Password); err != nil {
		t.Fatalf("expected password validation %v, got %v", nil, err)
	}
}

