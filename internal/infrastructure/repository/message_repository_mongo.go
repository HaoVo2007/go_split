package repository

import (
	"context"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messageRepositoryMongo struct {
	collection *mongo.Collection
}

func NewMessageRepositoryMongo(collection *mongo.Collection) repository.MessageRepository {
	return &messageRepositoryMongo{
		collection: collection,
	}
}

func (r *messageRepositoryMongo) CreateMessage(ctx context.Context, message entity.Messages) error {
	_, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepositoryMongo) GetMessagesByGroupID(ctx context.Context, groupID string, pageSize int, pageIndex int) ([]*entity.Messages, error) {
	filter := bson.M{
		"group_id": groupID,
	}

	if pageSize < 0 {
		pageSize = 10
	}

	if pageIndex < 0 {
		pageIndex = 1
	}

	skip := (pageIndex - 1) * pageSize

	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(pageSize))
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	messages := []*entity.Messages{}
	for cursor.Next(ctx) {
		message := &entity.Messages{}
		err := cursor.Decode(message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *messageRepositoryMongo) MarkSeenUpTo(ctx context.Context, groupID string, userID string, lastMessageIDs []primitive.ObjectID) error {
	filter := bson.M{
		"group_id": groupID,
		"_id": bson.M{
			"$in": lastMessageIDs,
		},
		"seen_by.user_id": bson.M{
			"$ne": userID,
		},
	}

	update := bson.M{
		"$push": bson.M{
			"seen_by": bson.M{
				"user_id": userID,
				"seen_at": time.Now(),
			},
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
