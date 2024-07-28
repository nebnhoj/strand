package todos

import (
	"context"
	"time"

	"github.com/nebnhoj/strand/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todoCollection *mongo.Collection = configs.GetCollection(configs.DB, "todos")

func create(todo Todo) (id any, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := todoCollection.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"id": result.InsertedID,
	}

	return data, nil
}

func getAllTodos(q string, page int, limit int) (todos []Todo, count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results []Todo
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default page size
	}
	skip := (page - 1) * limit
	options := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := todoCollection.Find(ctx,
		bson.M{
			"$or": []bson.M{
				{"name": primitive.Regex{Pattern: q, Options: "i"}},
				{"details": primitive.Regex{Pattern: q, Options: "i"}},
			},
		}, options)

	for cursor.Next(ctx) {
		var elem Todo
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, 0, err

		}
		results = append(results, elem)
	}
	total, err := todoCollection.CountDocuments(ctx,
		bson.M{
			"$or": []bson.M{
				{"name": primitive.Regex{Pattern: q, Options: "i"}},
				{"details": primitive.Regex{Pattern: q, Options: "i"}},
			},
		})
	if err != nil {
		return nil, 0, err
	}

	return results, total, err
}
