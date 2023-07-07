package api

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/thimc/go-backend/store"
	"github.com/thimc/go-backend/types"
)

const TokenHeader = "Authorization"

func JWT(s store.UserStorer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenHdr, ok := c.GetReqHeaders()[TokenHeader]
		if !ok {
			return fmt.Errorf("Missing token")
		}

		tokenParts := strings.Split(tokenHdr, " ")

		if len(tokenParts) < 2 {
			return fmt.Errorf("malformed token")
		}
		tokenStr := tokenParts[1]

		claims, err := validateToken(tokenStr)
		if err != nil {
			return fmt.Errorf("invalid token")
		}

		expires := int64(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired")
		}

		email := claims["email"].(string)
		user, err := s.GetUserByEmail(c.Context(), email)
		if err != nil {
			return err
		}

		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

// Validate by parsing the token with the JWT_SIGNING_KEY
func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_SIGNING_KEY")
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claim")
	}

	return claims, nil
}

func createJWT(user *types.User) (jwt.MapClaims, string, error) {
	claims := &jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": time.Now().Add(time.Hour*8).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SIGNING_KEY")
	tokenStr, err := token.SignedString([]byte(secret))

	return *claims, tokenStr, err
}
