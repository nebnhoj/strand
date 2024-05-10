package users

type User struct {
	Id        string  `bson:"_id" json:"id,omitempty"`
	FirstName string  `bson:"first_name" json:"first_name" validate:"required"`
	LastName  string  `bson:"last_name" json:"last_name" validate:"required"`
	Title     string  `json:"title" validate:"omitempty,required"`
	Email     string  `json:"email" validate:"omitempty,email"`
	Password  string  `json:"-" validate:"omitempty,required"`
	Address   Address `json:"address" validate:"required"`
}

type UserDTO struct {
	User
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
type Address struct {
	Street   string `json:"street" validate:"omitempty"`
	City     string `json:"city" validate:"required"`
	Province string `json:"province" validate:"omitempty"`
	Country  string `json:"country" validate:"omitempty"`
}
