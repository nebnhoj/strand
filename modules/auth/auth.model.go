package auth

type Auth struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
