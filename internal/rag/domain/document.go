package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	ID         primitive.ObjectID `bson:"_id"`
	Source     string             `bson:"source"`
	Content    string             `bson:"content"`
	Embeddings []float64          `bson:"embeddings"`
	ChunkIndex int                `bson:"chunk_index"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}
