package auth

import (
	"errors"

	message "github.com/nebnhoj/strand/helpers/errors"
	appJWT "github.com/nebnhoj/strand/middlewares/jwt"
	"github.com/nebnhoj/strand/modules/users"
	"golang.org/x/crypto/bcrypt"
)

func GetJWTToken(auth Auth) (map[string]interface{}, error) {
	user, err := users.GetUserByEmail(auth.Email)
	if err != nil {
		return nil, errors.New(message.UNAUTHORIZED)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
	if err != nil {

		return nil, errors.New(message.UNAUTHORIZED)
	}

	t, err := appJWT.CreateJWTClaim(user)
	if err != nil {
		return nil, err
	}
	token := map[string]interface{}{
		"token": t,
	}
	return token, nil
}
