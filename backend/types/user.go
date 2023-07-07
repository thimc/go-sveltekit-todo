package types

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
}

// UserParams is used when logging in and when we're creating a new user
type UserParams struct {
	// The users email address
	Email string `json:"email" validate:"required"`
	// The users password in plain text
	Password string `json:"password" validate:"required"`
} // @name UserParams

type LoginResponse struct {
	// The users ID
	ID int `json:"id"`
	// The users email address
	Email string `json:"email"`
	// The token
	Token string `json:"token"`
	// Unix timestamp for when the token expires
	ExpiresAt int64 `json:"expires"`
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
