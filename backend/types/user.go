package types

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	// The users ID
	ID int `json:"id" example:"0"`
	// The users email address
	Email string `json:"email" example:"user@domain.com"`
	// The users password in a encrypted format
	EncryptedPassword string `json:"-"`
} // @name User

// UserParams is used when logging in and when we're creating a new user
type UserParams struct {
	// The users email address
	Email string `json:"email" example:"user@domain.com" validate:"required"`
	// The users password in plain text
	Password string `json:"password" validate:"required"`
} // @name UserParams

func (p *UserParams) Validate() error {
	var b strings.Builder
	if len(p.Email) < 5 {
		b.WriteString(fmt.Sprintf("the email needs to be at least 5 characters\n"))
	}
	if len(p.Password) < 5 {
		b.WriteString(fmt.Sprintf("the password needs to be at least 5 characters\n"))
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(p.Email) {
		b.WriteString(fmt.Sprintf("the email needs to be a valid email address\n"))
	}
	if b.Len() > 1 {
		return fmt.Errorf("%s", b.String())
	}

	return nil
}

type LoginResponse struct {
	// The users ID
	ID int `json:"id" example:"0"`
	// The users email address
	Email string `json:"email" example:"user@domain.com"`
	// The token
	Token string `json:"token"`
	// Unix timestamp for when the token expires
	ExpiresAt int64 `json:"expires" example:"1688751625"`
} // @name LoginResponse

func NewUser(email, password string) (*User, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:             email,
		EncryptedPassword: string(encrypted),
	}, nil
}

func ValidPassword(encrytpedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(encrytpedPassword), []byte(password))
}
