package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"password"`
	Role         string             `bson:"role" json:"role"`
	Token        string             `bson:"token" json:"token"`
	RefreshToken string             `bson:"refresh_token" json:"refresh_token"`
	Status       string             `bson:"status" json:"status"`
	Profile      *Profile           `bson:"profile" json:"profile"`
	IsDeleted    bool               `bson:"is_deleted" json:"is_deleted"`
	DeletedAt    *time.Time         `bson:"deleted_at" json:"deleted_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type Profile struct {
	Name          *string    `bson:"name" json:"name"`
	Image         *string    `bson:"image" json:"image"`
	Address       *string    `bson:"address" json:"address"`
	Phone         *string    `bson:"phone" json:"phone"`
	ImagePublicID *string    `bson:"image_public_id" json:"image_public_id"`
}
