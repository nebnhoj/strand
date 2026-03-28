package mongodb

import (
	"context"
	"time"

	todoDomain "github.com/nebnhoj/strand/internal/domain/todo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTodo struct {
	ID      string `bson:"_id"`
	Name    string `bson:"name"`
	Details string `bson:"details"`
}

type todoRepository struct {
	col *mongo.Collection
}

func NewTodoRepository(client *mongo.Client, dbName string) todoDomain.Repository {
	return &todoRepository{col: client.Database(dbName).Collection("todos")}
}

func (r *todoRepository) Create(ctx context.Context, t todoDomain.Todo) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := r.col.InsertOne(ctx, mongoTodo{ID: t.ID, Name: t.Name, Details: t.Details})
	return t.ID, err
}

func (r *todoRepository) FindAll(ctx context.Context, q string, page, limit int) ([]todoDomain.Todo, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	filter := bson.M{"$or": []bson.M{
		{"name": primitive.Regex{Pattern: q, Options: "i"}},
		{"details": primitive.Regex{Pattern: q, Options: "i"}},
	}}

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit))
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []todoDomain.Todo
	for cursor.Next(ctx) {
		var doc mongoTodo
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		results = append(results, todoDomain.Todo{ID: doc.ID, Name: doc.Name, Details: doc.Details})
	}
	return results, total, nil
}
