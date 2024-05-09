package users

type User struct {
	Id       string `bson:"_id" json:"id,omitempty"`
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"omitempty"`
	Title    string `json:"title" validate:"omitempty,required"`
	Email    string `json:"email" validate:"omitempty,email"`
}
