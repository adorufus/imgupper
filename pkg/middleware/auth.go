package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/adorufus/imgupper/pkg/httputil"
	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret         string
	ExpirationTime time.Duration
}

type UserClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type contextKey string

func (c contextKey) String() string {
	return "middleware context key " + string(c)
}

const UserKey contextKey = "user"

func JWTAuth(config JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				httputil.ErrorResponse(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httputil.ErrorResponse(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims := &UserClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Unexpected signing method")
				}

				return []byte(config.Secret), nil
			})

			if err != nil {
				httputil.ErrorResponse(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				httputil.ErrorResponse(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add claims to request context
			ctx := context.WithValue(r.Context(), UserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*UserClaims, error) {
	userClaims, ok := ctx.Value(UserKey).(*UserClaims)
	if !ok {
		return nil, errors.New("User not found in context")
	}

	return userClaims, nil
}

func GenerateToken(userID int64, email string, config JWTConfig) (string, error) {
	expirationTime := time.Now().Add(config.ExpirationTime)
	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
