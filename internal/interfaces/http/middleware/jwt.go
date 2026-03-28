package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	authDomain "github.com/nebnhoj/strand/internal/domain/auth"
	"github.com/nebnhoj/strand/pkg/apperrors"
	"github.com/nebnhoj/strand/pkg/response"
)

type JWTMiddleware struct {
	tokens authDomain.TokenService
}

func NewJWTMiddleware(tokens authDomain.TokenService) *JWTMiddleware {
	return &JWTMiddleware{tokens: tokens}
}

func (m *JWTMiddleware) Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Error(c, http.StatusBadRequest, apperrors.JWT_ERROR)
		}

		claims, err := m.tokens.Parse(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			if errors.Is(err, authDomain.ErrTokenExpired) {
				return response.Error(c, http.StatusUnauthorized, apperrors.EXPIRED_TOKEN)
			}
			return response.Error(c, http.StatusUnauthorized, apperrors.JWT_ERROR)
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

// RequirePermission returns a middleware that allows the request only if the
// JWT contains the specified permission. Must be used after Protected().
func (m *JWTMiddleware) RequirePermission(perm string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("claims").(authDomain.Claims)
		if !ok {
			return response.Error(c, http.StatusUnauthorized, apperrors.JWT_ERROR)
		}
		for _, p := range claims.Permissions {
			if p == perm {
				return c.Next()
			}
		}
		return response.Error(c, http.StatusForbidden, apperrors.FORBIDDEN)
	}
}
