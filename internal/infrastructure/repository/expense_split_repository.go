package repository

import (
	"context"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type expenseSplitRepository struct {
	collection *mongo.Collection
}

func NewExpenseSplitRepository(collection *mongo.Collection) repository.ExpenseSplitRepository {
	return &expenseSplitRepository{
		collection: collection,
	}
}

func (r *expenseSplitRepository) CreateExpenseSplits(ctx context.Context, expenseSplits []entity.ExpenseSplits) error {
	documents := make([]interface{}, len(expenseSplits))
	for i, split := range expenseSplits {
		documents[i] = split
	}

	_, err := r.collection.InsertMany(ctx, documents)
	if err != nil {
		return err
	}
	return nil
}

func (r *expenseSplitRepository) GetExpenseSplitsByExpenseID(ctx context.Context, expenseID string) ([]*entity.ExpenseSplits, error) {
	filter := bson.M{
		"expenses_id": expenseID,
	}

	splits := []*entity.ExpenseSplits{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		split := &entity.ExpenseSplits{}
		err := cursor.Decode(split)
		if err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return splits, nil
}

func (r *expenseSplitRepository) GetExpenseSplitsByExpenseIDs(ctx context.Context, expenseIDs []string) ([]*entity.ExpenseSplits, error) {
	filter := bson.M{
		"expenses_id": bson.M{
			"$in": expenseIDs,
		},
	}

	splits := []*entity.ExpenseSplits{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		split := &entity.ExpenseSplits{}
		err := cursor.Decode(split)
		if err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return splits, nil

}
