package auth

import "errors"

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Claims holds the decoded values from a JWT token.
type Claims struct {
	UserID      string
	Email       string
	Roles       []string
	Permissions []string
}

// TokenService is the port the application layer uses to create and verify tokens.
// The concrete implementation lives in infrastructure/jwt.
type TokenService interface {
	Generate(userID, email string, roles []string, permissions []string) (string, error)
	Parse(tokenStr string) (Claims, error)
}
