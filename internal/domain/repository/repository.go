package repository

import (
	"context"
	"go-split/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.Users) error
	FindUserByEmail(ctx context.Context, email string) (*entity.Users, error)
	FindUserByID(ctx context.Context, userID primitive.ObjectID) (*entity.Users, error)
	UpdateUser(ctx context.Context, userID primitive.ObjectID, updateFields bson.M) error
	GetUsers(ctx context.Context) ([]*entity.Users, error)
	GetUsersByGroupIDs(ctx context.Context, userIDs []primitive.ObjectID) ([]*entity.Users, error)
	GetUsersByIDs(ctx context.Context, userIDs []primitive.ObjectID) ([]*entity.Users, error)
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, group entity.Groups) error
	GetGroups(ctx context.Context, userID string) ([]*entity.Groups, error)
	GetGroupById(ctx context.Context, groupID primitive.ObjectID) (*entity.Groups, error)
	UpdateGroup(ctx context.Context, groupID primitive.ObjectID, group *entity.Groups) error
	DeleteGroup(ctx context.Context, groupID primitive.ObjectID) error
}

type ExpenseRepository interface {
	CreateExpense(ctx context.Context, expense entity.Expenses) error
	GetExpenseById(ctx context.Context, expenseID primitive.ObjectID) (*entity.Expenses, error)
	UpdateExpenseById(ctx context.Context, expenseID primitive.ObjectID, expense *entity.Expenses) error
	GetExpensesByGroupID(ctx context.Context, groupID string) ([]*entity.Expenses, error)
	GetExpensesByGroupIDs(ctx context.Context, groupIDs []string) ([]*entity.Expenses, error)
}

type ExpenseSplitRepository interface {
	CreateExpenseSplits(ctx context.Context, expenseSplits []entity.ExpenseSplits) error
	GetExpenseSplitsByExpenseIDs(ctx context.Context, expenseIDs []string) ([]*entity.ExpenseSplits, error)
	GetExpenseSplitsByExpenseID(ctx context.Context, expenseID string) ([]*entity.ExpenseSplits, error)
	DeleteExpenseSplitsByExpenseID(ctx context.Context, expenseID string) error
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, message entity.Messages) error
	GetMessagesByGroupID(ctx context.Context, groupID string, pageSize int, pageIndex int) ([]*entity.Messages, int64, error)
	MarkSeenUpTo(ctx context.Context, groupID string, userID string, lastMessageIDs []primitive.ObjectID) error
	GetUnreadCounts(ctx context.Context, groupIDs []string, userID string) (map[string]int, error)
	GetUnreadCount(ctx context.Context, groupID string, userID string) (int, error)
}
