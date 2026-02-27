package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {

	if err := createUserIndexes(ctx, db); err != nil {
		return err
	}

	log.Println("All database indexes created successfully")
	return nil
}

func createUserIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
		{
			Keys:    bson.M{"role": 1},
			Options: options.Index().SetName("role_index"),
		},
		{
			Keys:    bson.M{"status": 1},
			Options: options.Index().SetName("status_index"),
		},
		{
			Keys:    bson.D{{Key: "is_deleted", Value: 1}, {Key: "created_at", Value: -1}},
			Options: options.Index().SetName("is_deleted_created_at_index"),
		},
		{
			Keys:    bson.M{"created_at": -1},
			Options: options.Index().SetName("created_at_desc_index"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("Error creating user indexes: %v", err)
		return err
	}

	log.Println("User indexes created successfully")
	return nil
}
