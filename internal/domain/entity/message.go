package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Messages struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	GroupID   string             `bson:"group_id" json:"group_id"`
	Message   string             `bson:"message" json:"message"`
	UserID    string             `bson:"user_id" json:"user_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
