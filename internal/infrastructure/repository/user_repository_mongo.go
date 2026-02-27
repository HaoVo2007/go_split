package repository

import (
	"context"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepositoryMongo(collection *mongo.Collection) repository.UserRepository {
	return &userRepositoryMongo{
		collection: collection,
	}
}

func (r *userRepositoryMongo) CreateUser(ctx context.Context, user entity.Users) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryMongo) FindUserByEmail(ctx context.Context, email string) (*entity.Users, error) {
	filter := bson.M{
		"email":      email,
		"is_deleted": false,
	}
	user := &entity.Users{}
	err := r.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryMongo) FindUserByID(ctx context.Context, userID primitive.ObjectID) (*entity.Users, error) {
	filter := bson.M{
		"_id": userID,
	}
	user := &entity.Users{}
	err := r.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryMongo) UpdateUser(ctx context.Context, userID primitive.ObjectID, updateFields bson.M) error {
	filter := bson.M{
		"_id": userID,
	}

	updateDoc := bson.M{
		"$set": updateFields,
	}

	_, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryMongo) GetUsers(ctx context.Context) ([]*entity.Users, error) {
	filter := bson.M{
		"is_deleted": false,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := []*entity.Users{}
	for cursor.Next(ctx) {
		user := &entity.Users{}
		err := cursor.Decode(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepositoryMongo) GetUsersByGroupIDs(ctx context.Context, userIDs []primitive.ObjectID) ([]*entity.Users, error) {

	filter := bson.M{
		"_id": bson.M{
			"$in": userIDs,
		},
		"is_deleted": false,
	}

	users := []*entity.Users{}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		user := &entity.Users{}
		err := cursor.Decode(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepositoryMongo) GetUsersByIDs(ctx context.Context, userIDs []primitive.ObjectID) ([]*entity.Users, error) {
	filter := bson.M{
		"_id": bson.M{
			"$in": userIDs,
		},
		"is_deleted": false,
	}

	users := []*entity.Users{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		user := &entity.Users{}
		err := cursor.Decode(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
