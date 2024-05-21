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

// func IsAdmin(c *fiber.Ctx) error {
// 	user := c.Locals("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	isAdmin := claims["isAdmin"]

// 	if !isAdmin.(bool) {
// 		return helpers.ResponseError(c, http.StatusForbidden, error_message.FORBIDDEN)

// 	}
// 	return c.Next()

// }

type CustomClaims struct {
	jwt.Claims                `json:"-"`
	Email                     string   `json:"email,omitempty"`
	Roles                     []string `json:"roles,omitempty"`
	Expiration                int64    `json:"exp"`
	Sub                       string   `json:"sub"`
	UserLastName              string   `json:"userlastname"`
	UserEstablishmentUUID     string   `json:"userestablishmentuuid"`
	UserFirstName             string   `json:"userfirstname"`
	Iss                       string   `json:"iss"`
	UserLastPasswordResetDate string   `json:"userlastpasswordresetdate"`
	UserIsEnabled             string   `json:"userisenabled"`
	UserUUID                  string   `json:"useruuid"`
	UserID                    string   `json:"userid"`
	Aud                       string   `json:"aud"`
	Nbf                       int64    `json:"nbf"`
	UserRoles                 string   `json:"userroles"`
	UserMiddleName            string   `json:"usermiddlename"`
	Iat                       int64    `json:"iat"`
	Jti                       string   `json:"jti"`
	UserEmail                 string   `json:"useremail"`
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
		Email:                     user.Email,
		Roles:                     user.Roles,
		Expiration:                time.Now().Add(time.Hour * 72).Unix(),
		UserRoles:                 "2,ROLE_STASH_ADMIN",
		Iat:                       1715756829,
		UserEmail:                 user.Email,
		Jti:                       "957c39d7a29f62a454af2f5af7cdfe75",
		UserIsEnabled:             "true",
		Aud:                       "web",
		UserUUID:                  "1c34d0fb-f347-4d24-8a65-e1bc11fb6280",
		UserLastPasswordResetDate: "1714984122000",
		Iss:                       "stash.ph",
		UserLastName:              user.LastName,
		UserFirstName:             user.FirstName,
		Sub:                       user.Email,
		UserID:                    "1",
		Nbf:                       1715756829,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(JWTSecretKey()))
}
