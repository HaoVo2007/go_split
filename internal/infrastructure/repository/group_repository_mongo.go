package repository

import (
	"context"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type groupRepositoryMongo struct {
	collection *mongo.Collection
}

func NewGroupRepositoryMongo(collection *mongo.Collection) repository.GroupRepository {
	return &groupRepositoryMongo{
		collection: collection,
	}
}

func (r *groupRepositoryMongo) CreateGroup(ctx context.Context, group entity.Groups) error {
	_, err := r.collection.InsertOne(ctx, group)
	if err != nil {
		return err
	}
	return nil
}

func (r *groupRepositoryMongo) GetGroups(ctx context.Context, userID string) ([]*entity.Groups, error) {
	filter := bson.M{
		"is_deleted": false,
		"$or": []bson.M{
			{
				"created_by": userID,
			},
			{
				"members": userID,
			},
		},
	}

	groups := []*entity.Groups{}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		group := &entity.Groups{}
		err := cursor.Decode(group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *groupRepositoryMongo) GetGroupById(ctx context.Context, groupID primitive.ObjectID) (*entity.Groups, error) {
	filter := bson.M{
		"_id":        groupID,
		"is_deleted": false,
	}
	group := &entity.Groups{}
	err := r.collection.FindOne(ctx, filter).Decode(group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return group, nil
}

func (r *groupRepositoryMongo) UpdateGroup(ctx context.Context, groupID primitive.ObjectID, group *entity.Groups) error {
	filter := bson.M{
		"_id": groupID,
	}
	update := bson.M{
		"$set": group,
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *groupRepositoryMongo) DeleteGroup(ctx context.Context, groupID primitive.ObjectID) error {
	filter := bson.M{
		"_id": groupID,
	}
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
