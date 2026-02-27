package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expenses struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	GroupID   string             `bson:"group_id" json:"group_id"`
	Date      time.Time          `bson:"date" json:"date"`
	Name      string             `bson:"name" json:"name"`
	Amount    float64            `bson:"amount" json:"amount"`
	Category  string             `bson:"category" json:"category"`
	PaidBy    []string           `bson:"paid_by" json:"paid_by"`
	CreatedBy string             `bson:"created_by" json:"created_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type ExpenseSplits struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	ExpensesID    string             `bson:"expenses_id" json:"expenses_id"`
	UserId        string             `bson:"user_id" json:"user_id"`
	Amount        float64            `bson:"amount" json:"amount"`
	Description   string             `bson:"description" json:"description"`
	Image         string             `bson:"image" json:"image"`
	ImagePublicID string             `bson:"image_public_id" json:"image_public_id"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
