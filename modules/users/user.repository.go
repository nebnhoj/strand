package users

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schuler.com/be-schuler/configs"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func create(user User) (id any, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"id": result.InsertedID,
	}

	return data, nil
}

func getAllUsers(q string, page int, limit int) (users []User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results []User
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default page size
	}
	skip := (page - 1) * limit
	options := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := userCollection.Find(ctx,
		bson.M{
			"$or": []bson.M{
				{"first_name": primitive.Regex{Pattern: q, Options: "i"}},
				{"last_name": primitive.Regex{Pattern: q, Options: "i"}},
			},
		}, options)

	for cursor.Next(ctx) {
		var elem User
		err := cursor.Decode(&elem)
		if err != nil {
			return users, err

		}
		results = append(results, elem)
	}
	if len(results) < 1 {
		return []User{}, err
	}
	return results, err
}

func getUserByID(id string) (user User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result User
	err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	return result, err
}
