package user

// @name User
type User struct {
	ID          string   `json:"id,omitempty"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Title       string   `json:"title,omitempty"`
	Email       string   `json:"email,omitempty"`
	Password    string   `json:"-"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Address     Address  `json:"address"`
}

// @name Address
type Address struct {
	Street   string `json:"street,omitempty"`
	City     string `json:"city"`
	Province string `json:"province,omitempty"`
	Country  string `json:"country,omitempty"`
}
