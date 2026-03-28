package handlers

// Swagger request/response types. Used only in annotation comments so the
// generated model names are clean instead of full Go import paths.

type AuthResponse struct {
	Token       string   `json:"token"`
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
} // @name AuthResponse

type AuthRequest struct {
	Email    string `json:"email"    example:"admin@strand.dev"`
	Password string `json:"password" example:"password"`
} // @name AuthRequest

type CreateUserRequest struct {
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Title           string         `json:"title,omitempty"`
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	ConfirmPassword string         `json:"confirm_password"`
	Roles           []string       `json:"roles"`
	Address         AddressRequest `json:"address"`
} // @name CreateUserRequest

type AddressRequest struct {
	Street   string `json:"street,omitempty"`
	City     string `json:"city"`
	Province string `json:"province,omitempty"`
	Country  string `json:"country,omitempty"`
} // @name AddressRequest

type CreateTodoRequest struct {
	Name    string `json:"name"`
	Details string `json:"details,omitempty"`
} // @name CreateTodoRequest

type UserResponse struct {
	ID        string         `json:"id,omitempty"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Title     string         `json:"title,omitempty"`
	Email     string         `json:"email,omitempty"`
	Roles     []string       `json:"roles,omitempty"`
	Address   AddressRequest `json:"address"`
} // @name UserResponse

type TodoResponse struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Details string `json:"details"`
} // @name TodoResponse

type ErrorResponse struct {
	Status  int `json:"status"`
	Message any `json:"message"`
} // @name ErrorResponse
