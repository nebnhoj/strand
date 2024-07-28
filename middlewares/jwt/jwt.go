package jwt

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/nebnhoj/strand/helpers"
	error_message "github.com/nebnhoj/strand/helpers/errors"
	"github.com/nebnhoj/strand/modules/users"

	"github.com/gofiber/fiber/v2"
)

func JWTSecretKey() string {
	return os.Getenv("JWT_SECRET")
}
func Protected() func(*fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	jwtwareConfig := jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(JWTSecretKey())},
		ContextKey:     "user", // used in private route
		ErrorHandler:   jwtError,
		SuccessHandler: verifyTokenExpiration,
	}

	return jwtware.New(jwtwareConfig)
}

type CustomClaims struct {
	jwt.Claims `json:"-"`
	ID         string   `json:"id,omitempty"`
	Email      string   `json:"email,omitempty"`
	Roles      []string `json:"roles,omitempty"`
	Expiration int64    `json:"exp"`
}

func HasAdminRole(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	roles := claims["roles"].([]interface{})

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

func verifyTokenExpiration(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	expires := int64(claims["exp"].(float64))
	if time.Now().Unix() > expires {
		return jwtError(c, errors.New(error_message.EXPIRED_TOKEN))
	}
	return c.Next()
}
func jwtError(c *fiber.Ctx, err error) error {
	// Return Bad Request if Invalid or missing JWT.
	if err.Error() == "missing or malformed JWT" {
		return helpers.ResponseError(c, http.StatusBadRequest, error_message.JWT_ERROR)
	}

	// Return status 401 and failed authentication error.
	return helpers.ResponseError(c, http.StatusUnauthorized, err.Error())
}

func CreateJWTClaim(user users.User) (string, error) {
	claims := CustomClaims{
		ID:         user.Id,
		Email:      user.Email,
		Roles:      user.Roles,
		Expiration: time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(JWTSecretKey()))
}
