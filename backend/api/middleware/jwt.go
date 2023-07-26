package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/thimc/go-svelte-todo/backend/store"
	"github.com/thimc/go-svelte-todo/backend/types"
	"github.com/thimc/go-svelte-todo/backend/utils"
)

type JWTMiddleware struct {
	store store.UserStorer
}
func NewJWTMiddleware(store store.UserStorer) *JWTMiddleware {
	return &JWTMiddleware{
		store: store,
	}
}

func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			utils.WriteJSON(w, types.NewAPIError(false, fmt.Errorf("Missing token"), http.StatusBadRequest))
			return
		}
		tokenArr := strings.Split(tokenHeader, " ")
		if len(tokenArr) < 2 {
			utils.WriteJSON(w, types.NewAPIError(false, fmt.Errorf("Malformed token"), http.StatusBadRequest))
			return
		}

		tok, err := ValidateJWT(tokenArr[1])
		if err != nil || !tok.Valid {
			utils.WriteJSON(w, types.NewAPIError(false, fmt.Errorf("Invalid token"), http.StatusBadRequest))
			return
		}

		claims := tok.Claims.(jwt.MapClaims)
		if time.Now().Unix() > int64(claims["expiresAt"].(float64)) {
			utils.WriteJSON(w, types.NewAPIError(false, fmt.Errorf("Token expired"), http.StatusUnauthorized))
			return
		}

		email := claims["email"].(string)
		user, err := m.store.GetUserByEmail(r.Context(), email)
		if err != nil {
			utils.WriteJSON(w, types.NewAPIError(false, fmt.Errorf("Access denied"), http.StatusUnauthorized))
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateJWT(user *types.User) (jwt.MapClaims, string, error) {
	claims := &jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"expiresAt": time.Now().Add(time.Hour * 6).Unix(),
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte(secret))

	return *claims, tok, err
}

func ValidateJWT(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token")
		}
		return []byte(secret), nil
	})
}
