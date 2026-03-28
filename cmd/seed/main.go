package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nebnhoj/strand/configs"
	"github.com/nebnhoj/strand/modules/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load()

	db := configs.DB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	seedUsers(ctx, db)
	seedTodos(ctx, db)

	log.Println("Seeding complete.")
}

func seedUsers(ctx context.Context, db *mongo.Client) {
	col := configs.GetCollection(db, "users")

	col.DeleteMany(ctx, bson.M{})

	records := []users.User{
		{
			Id:        uuid.NewString(),
			FirstName: "Admin",
			LastName:  "User",
			Title:     "Mr",
			Email:     "admin@strand.dev",
			Password:  hash("password"),
			Roles:     []string{"ADMIN"},
			Address: users.Address{
				Street:   "1 Main St",
				City:     "New York",
				Province: "NY",
				Country:  "US",
			},
		},
		{
			Id:        uuid.NewString(),
			FirstName: "Jane",
			LastName:  "Doe",
			Title:     "Ms",
			Email:     "jane@strand.dev",
			Password:  hash("password"),
			Roles:     []string{"USER"},
			Address: users.Address{
				Street:   "42 Park Ave",
				City:     "Los Angeles",
				Province: "CA",
				Country:  "US",
			},
		},
		{
			Id:        uuid.NewString(),
			FirstName: "John",
			LastName:  "Smith",
			Title:     "Dr",
			Email:     "john@strand.dev",
			Password:  hash("password"),
			Roles:     []string{"USER"},
			Address: users.Address{
				Street:   "7 Ocean Blvd",
				City:     "Miami",
				Province: "FL",
				Country:  "US",
			},
		},
	}

	docs := make([]interface{}, len(records))
	for i, u := range records {
		u.Email = strings.ToLower(u.Email)
		docs[i] = u
	}

	_, err := col.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("seed users: %v", err)
	}
	log.Printf("Seeded %d users", len(docs))
}

type Todo struct {
	Id      string `bson:"_id" json:"id,omitempty"`
	Name    string `json:"name"`
	Details string `json:"details"`
}

func seedTodos(ctx context.Context, db *mongo.Client) {
	col := configs.GetCollection(db, "todos")

	col.DeleteMany(ctx, bson.M{})

	records := []Todo{
		{Id: uuid.NewString(), Name: "Buy groceries", Details: "Milk, eggs, bread, and coffee"},
		{Id: uuid.NewString(), Name: "Read a book", Details: "Finish reading The Pragmatic Programmer"},
		{Id: uuid.NewString(), Name: "Exercise", Details: "30-minute run in the morning"},
		{Id: uuid.NewString(), Name: "Write tests", Details: "Add integration tests for the auth module"},
		{Id: uuid.NewString(), Name: "Deploy app", Details: "Push latest changes to production via Docker"},
	}

	docs := make([]interface{}, len(records))
	for i, t := range records {
		docs[i] = t
	}

	_, err := col.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("seed todos: %v", err)
	}
	log.Printf("Seeded %d todos", len(docs))
}

func hash(password string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	return string(b)
}
