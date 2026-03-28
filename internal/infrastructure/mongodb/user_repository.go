package mongodb

import (
	"context"
	"time"

	userDomain "github.com/nebnhoj/strand/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoUser struct {
	ID          string       `bson:"_id"`
	FirstName   string       `bson:"first_name"`
	LastName    string       `bson:"last_name"`
	Title       string       `bson:"title"`
	Email       string       `bson:"email"`
	Password    string       `bson:"password"`
	Roles       []string     `bson:"roles"`
	Permissions []string     `bson:"permissions"`
	Address     mongoAddress `bson:"address"`
}

type mongoAddress struct {
	Street   string `bson:"street"`
	City     string `bson:"city"`
	Province string `bson:"province"`
	Country  string `bson:"country"`
}

type userRepository struct {
	col *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName string) userDomain.Repository {
	return &userRepository{col: client.Database(dbName).Collection("users")}
}

func (r *userRepository) Create(ctx context.Context, u userDomain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := r.col.InsertOne(ctx, toMongoUser(u))
	return u.ID, err
}

func (r *userRepository) FindAll(ctx context.Context, q string, page, limit int) ([]userDomain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	filter := bson.M{"$or": []bson.M{
		{"first_name": primitive.Regex{Pattern: q, Options: "i"}},
		{"last_name": primitive.Regex{Pattern: q, Options: "i"}},
	}}
	opts := options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit))

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []userDomain.User
	for cursor.Next(ctx) {
		var doc mongoUser
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		results = append(results, toDomainUser(doc))
	}
	return results, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (userDomain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var doc mongoUser
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		return userDomain.User{}, err
	}
	return toDomainUser(doc), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (userDomain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var doc mongoUser
	if err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&doc); err != nil {
		return userDomain.User{}, err
	}
	return toDomainUser(doc), nil
}

func toMongoUser(u userDomain.User) mongoUser {
	return mongoUser{
		ID: u.ID, FirstName: u.FirstName, LastName: u.LastName,
		Title: u.Title, Email: u.Email, Password: u.Password,
		Roles: u.Roles, Permissions: u.Permissions,
		Address: mongoAddress{
			Street: u.Address.Street, City: u.Address.City,
			Province: u.Address.Province, Country: u.Address.Country,
		},
	}
}

func toDomainUser(doc mongoUser) userDomain.User {
	return userDomain.User{
		ID: doc.ID, FirstName: doc.FirstName, LastName: doc.LastName,
		Title: doc.Title, Email: doc.Email, Password: doc.Password,
		Roles: doc.Roles, Permissions: doc.Permissions,
		Address: userDomain.Address{
			Street: doc.Address.Street, City: doc.Address.City,
			Province: doc.Address.Province, Country: doc.Address.Country,
		},
	}
}
