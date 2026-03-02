package repository

import (
	"context"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type expenseRepositoryMongo struct {
	collection *mongo.Collection
}

func NewExpenseRepositoryMongo(collection *mongo.Collection) repository.ExpenseRepository {
	return &expenseRepositoryMongo{
		collection: collection,
	}
}

func (r *expenseRepositoryMongo) CreateExpense(ctx context.Context, expense entity.Expenses) error {
	_, err := r.collection.InsertOne(ctx, expense)
	if err != nil {
		return err
	}
	return nil
}

func (r *expenseRepositoryMongo) GetExpenseById(ctx context.Context, expenseID primitive.ObjectID) (*entity.Expenses, error) {
	filter := bson.M{
		"_id":        expenseID,
		"is_deleted": false,
	}
	expense := &entity.Expenses{}
	err := r.collection.FindOne(ctx, filter).Decode(expense)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return expense, nil
}

func (r *expenseRepositoryMongo) UpdateExpenseById(ctx context.Context, expenseID primitive.ObjectID, expense *entity.Expenses) error {
	filter := bson.M{
		"_id": expenseID,
	}

	update := bson.M{
		"$set": expense,
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *expenseRepositoryMongo) GetExpensesByGroupID(ctx context.Context, groupID string) ([]*entity.Expenses, error) {
	filter := bson.M{
		"group_id":   groupID,
		"is_deleted": false,
	}

	expenses := []*entity.Expenses{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		expense := &entity.Expenses{}
		err := cursor.Decode(expense)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *expenseRepositoryMongo) GetExpensesByGroupIDs(ctx context.Context, groupIDs []string) ([]*entity.Expenses, error) {
	filter := bson.M{
		"group_id": bson.M{
			"$in": groupIDs,
		},
		"is_deleted": false,
	}

	expenses := []*entity.Expenses{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		expense := &entity.Expenses{}
		err := cursor.Decode(expense)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}
