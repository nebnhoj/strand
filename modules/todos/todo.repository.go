package todos

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nebnhoj/strand/configs"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todoCollection *mongo.Collection = configs.GetCollection(configs.DB, "todos")
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379", // Redis server address
	Password: "",               // No password
	DB:       0,                // Default database
})

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

func getAllTodos(q string, page int, limit int) (todos []Todo, err error) {
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
			return nil, err

		}
		results = append(results, elem)
	}
	if len(results) < 1 {
		return []Todo{}, err
	}

	if err != nil {
		return nil,err
	}

	jsonData, err := json.Marshal(todos)
	if err != nil {
		return nil,err

	}
	err = rdb.Set(context.Background(), "todos", jsonData, 60*time.Hour).Err()
	if err != nil {
		return nil,err

	}
	return results, err
}
