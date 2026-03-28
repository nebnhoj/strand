package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v3"
	"github.com/nebnhoj/strand/helpers"
	error_message "github.com/nebnhoj/strand/helpers/errors"
	"github.com/nebnhoj/strand/modules/users"
)

func JWTSecretKey() string {
	return os.Getenv("JWT_SECRET")
}

func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return jwtError(c, errors.New("missing or malformed JWT"))
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(JWTSecretKey()), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return jwtError(c, errors.New(error_message.EXPIRED_TOKEN))
			}
			return jwtError(c, errors.New(error_message.JWT_ERROR))
		}

		c.Locals("user", token)
		return c.Next()
	}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	ID    string   `json:"id,omitempty"`
	Email string   `json:"email,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

func HasAdminRole(c fiber.Ctx) error {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return helpers.ResponseError(c, http.StatusUnauthorized, error_message.JWT_ERROR)
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return helpers.ResponseError(c, http.StatusUnauthorized, error_message.JWT_ERROR)
	}
	roles, ok := claims["roles"].([]interface{})
	if !ok {
		return helpers.ResponseError(c, http.StatusForbidden, error_message.FORBIDDEN)
	}
	if !valueExists(roles, "ADMIN") {
		return helpers.ResponseError(c, http.StatusForbidden, error_message.FORBIDDEN)
	}
	return c.Next()
}

func valueExists(arr []interface{}, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func jwtError(c fiber.Ctx, err error) error {
	if err.Error() == "missing or malformed JWT" {
		return helpers.ResponseError(c, http.StatusBadRequest, error_message.JWT_ERROR)
	}
	return helpers.ResponseError(c, http.StatusUnauthorized, err.Error())
}

func CreateJWTClaim(user users.User) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
		ID:    user.Id,
		Email: user.Email,
		Roles: user.Roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(JWTSecretKey()))
}
