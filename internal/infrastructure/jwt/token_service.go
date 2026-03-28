package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	authDomain "github.com/nebnhoj/strand/internal/domain/auth"
)

type tokenService struct {
	secret string
}

func NewTokenService(secret string) authDomain.TokenService {
	return &tokenService{secret: secret}
}

type claims struct {
	jwtlib.RegisteredClaims
	ID          string   `json:"id,omitempty"`
	Email       string   `json:"email,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func (s *tokenService) Generate(userID, email string, roles []string, permissions []string) (string, error) {
	c := claims{
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
		ID:          userID,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS512, c)
	return token.SignedString([]byte(s.secret))
}

func (s *tokenService) Parse(tokenStr string) (authDomain.Claims, error) {
	token, err := jwtlib.Parse(tokenStr, func(t *jwtlib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return authDomain.Claims{}, authDomain.ErrTokenExpired
		}
		return authDomain.Claims{}, authDomain.ErrInvalidToken
	}

	mc, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		return authDomain.Claims{}, authDomain.ErrInvalidToken
	}

	return authDomain.Claims{
		UserID:      fmt.Sprintf("%v", mc["id"]),
		Email:       fmt.Sprintf("%v", mc["email"]),
		Roles:       toStringSlice(mc["roles"]),
		Permissions: toStringSlice(mc["permissions"]),
	}, nil
}

func toStringSlice(v interface{}) []string {
	arr, _ := v.([]interface{})
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
