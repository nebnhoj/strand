package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	authDomain "github.com/nebnhoj/strand/internal/domain/auth"
	"github.com/nebnhoj/strand/internal/domain/todo"
	"github.com/nebnhoj/strand/internal/domain/user"
	"github.com/nebnhoj/strand/internal/infrastructure/mongodb"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load()

	db := mongodb.ConnectDB()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userRepo := mongodb.NewUserRepository(db, "strand")
	todoRepo := mongodb.NewTodoRepository(db, "strand")

	seedUsers(ctx, userRepo)
	seedTodos(ctx, todoRepo)

	log.Println("Seeding complete.")
}

func seedUsers(ctx context.Context, repo user.Repository) {
	records := []user.User{
		{
			ID:        uuid.NewString(),
			FirstName: "Admin",
			LastName:  "User",
			Title:     "Mr",
			Email:     "admin@strand.dev",
			Password:  hash("password"),
			Roles:     []string{"ADMIN"},
			Address:   user.Address{Street: "1 Main St", City: "New York", Province: "NY", Country: "US"},
		},
		{
			ID:        uuid.NewString(),
			FirstName: "Jane",
			LastName:  "Doe",
			Title:     "Ms",
			Email:     "jane@strand.dev",
			Password:  hash("password"),
			Roles:     []string{"USER"},
			Address:   user.Address{Street: "42 Park Ave", City: "Los Angeles", Province: "CA", Country: "US"},
		},
		{
			ID:        uuid.NewString(),
			FirstName: "John",
			LastName:  "Smith",
			Title:     "Dr",
			Email:     "john@strand.dev",
			Password:  hash("password"),
			Roles:     []string{"USER"},
			Address:   user.Address{Street: "7 Ocean Blvd", City: "Miami", Province: "FL", Country: "US"},
		},
	}

	for i := range records {
		records[i].Email = strings.ToLower(records[i].Email)
		records[i].Permissions = permissionsFromRoles(records[i].Roles)
		if _, err := repo.Create(ctx, records[i]); err != nil {
			log.Printf("seed user %s: %v", records[i].Email, err)
		}
	}
	log.Printf("Seeded %d users", len(records))
}

func permissionsFromRoles(roles []string) []string {
	seen := make(map[string]struct{})
	var perms []string
	for _, role := range roles {
		for _, p := range authDomain.RolePermissions[role] {
			if _, exists := seen[p]; !exists {
				seen[p] = struct{}{}
				perms = append(perms, p)
			}
		}
	}
	return perms
}

func seedTodos(ctx context.Context, repo todo.Repository) {
	records := []todo.Todo{
		{ID: uuid.NewString(), Name: "Buy groceries", Details: "Milk, eggs, bread, and coffee"},
		{ID: uuid.NewString(), Name: "Read a book", Details: "Finish reading The Pragmatic Programmer"},
		{ID: uuid.NewString(), Name: "Exercise", Details: "30-minute run in the morning"},
		{ID: uuid.NewString(), Name: "Write tests", Details: "Add integration tests for the auth module"},
		{ID: uuid.NewString(), Name: "Deploy app", Details: "Push latest changes to production via Docker"},
	}

	for _, t := range records {
		if _, err := repo.Create(ctx, t); err != nil {
			log.Printf("seed todo %s: %v", t.Name, err)
		}
	}
	log.Printf("Seeded %d todos", len(records))
}

func hash(password string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	return string(b)
}
