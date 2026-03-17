package database

import (
	"context"
	"fmt"
	"go-split/pkg/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConnection(config config.MongoDBConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uri string
	if config.User != "" && config.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.User, config.Password, config.Host, config.Port)
	} else if config.URL != "" {
		uri = config.URL
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(config.DBName)

	if err := CreateIndexes(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to create database indexes: %w", err)
	}

	return db, nil
}
