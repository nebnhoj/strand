package auth

// Fine-grained permissions used across the application.
const (
	PermUsersRead  = "users:read"
	PermUsersWrite = "users:write"
	PermTodosRead  = "todos:read"
	PermTodosWrite = "todos:write"
)

// RolePermissions maps each role name to the set of permissions it grants.
// Add new roles here; no middleware code needs to change.
var RolePermissions = map[string][]string{
	"ADMIN": {PermUsersRead, PermUsersWrite, PermTodosRead, PermTodosWrite},
	"USER":  {PermTodosRead, PermTodosWrite},
}
