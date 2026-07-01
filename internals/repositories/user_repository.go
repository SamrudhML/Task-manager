package repositories

import (
	"context"
	"errors"
	"time"

	"task-manager/internals/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *UserRepository {
	return &UserRepository{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		user.ID = primitive.NewObjectID().Hex()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
