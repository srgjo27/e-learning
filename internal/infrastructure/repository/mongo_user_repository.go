package repository

import (
	"context"
	"errors"

	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(c *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{
		collection: c,
	}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return entity.ErrEmailExists
	}

	return err
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	filter := bson.M{"email": email}
	var user entity.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, entity.ErrUserNotFound
	}

	return &user, err
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, entity.ErrUserNotFound
	}
	
	filter := bson.M{"_id": oid}
	var user entity.User
	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, entity.ErrUserNotFound
	}

	return &user, err
}

func (r *MongoUserRepository) UpdatePassword(ctx context.Context, userID string, newHashed string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return entity.ErrUserNotFound
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"password": newHashed}}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

func (r *MongoUserRepository) UpdateEmail(ctx context.Context, userID string, newEmail string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return entity.ErrUserNotFound
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"email": newEmail}}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}