package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Groups struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	Image         string             `bson:"image" json:"image"`
	ImagePublicID string             `bson:"image_public_id" json:"image_public_id"`
	Members       []string           `bson:"members" json:"members"`
	IsDeleted     bool               `bson:"is_deleted" json:"is_deleted"`
	CreatedBy     string             `bson:"created_by" json:"created_by"`
	DeletedAt     *time.Time         `bson:"deleted_at" json:"deleted_at"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
