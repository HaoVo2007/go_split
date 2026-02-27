package usecase

import (
	"context"
	"errors"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"
	expenseMapper "go-split/internal/interface/http/dto/mapper/expense"
	"go-split/internal/interface/http/dto/request/expense"
	expenseRes "go-split/internal/interface/http/dto/response/expense"
	"go-split/pkg/libs/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseUseCase interface {
	CreateExpense(ctx context.Context, req expense.CreateExpenseRequest) error
	GetExpensesByGroupID(ctx context.Context, groupID string) ([]*expenseRes.ExpenseResponse, error)
}

type expenseUseCase struct {
	expenseRepository repository.ExpenseRepository
	groupRepository   repository.GroupRepository
	userRepository    repository.UserRepository
}

func NewExpenseUseCase(
	expenseRepository repository.ExpenseRepository,
	groupRepository repository.GroupRepository,
	userRepository repository.UserRepository,
) ExpenseUseCase {
	return &expenseUseCase{
		expenseRepository: expenseRepository,
		groupRepository:   groupRepository,
		userRepository:    userRepository,
	}
}

func (e *expenseUseCase) CreateExpense(ctx context.Context, req expense.CreateExpenseRequest) error {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return err
	}

	groupID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return err
	}

	group, err := e.groupRepository.GetGroupById(ctx, groupID)
	if err != nil {
		return err
	}

	if group == nil {
		return errors.New("group not found")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	expense := entity.Expenses{
		ID:        primitive.NewObjectID(),
		GroupID:   req.GroupID,
		Date:      date,
		Name:      req.Name,
		Amount:    req.Amount,
		Category:  req.Category,
		PaidBy:    req.PaidBy,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return e.expenseRepository.CreateExpense(ctx, expense)
}

func (e *expenseUseCase) GetExpensesByGroupID(ctx context.Context, groupID string) ([]*expenseRes.ExpenseResponse, error) {
	expenses, err := e.expenseRepository.GetExpensesByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	paidByIDs := []primitive.ObjectID{}
	for _, expense := range expenses {
		for _, paidByID := range expense.PaidBy {
			paidByIDObject, err := primitive.ObjectIDFromHex(paidByID)
			if err != nil {
				return nil, err
			}
			paidByIDs = append(paidByIDs, paidByIDObject)
		}
	}

	paidByUsers, err := e.userRepository.GetUsersByIDs(ctx, paidByIDs)
	if err != nil {
		return nil, err
	}

	paidByUsersMap := make(map[string]*entity.Users)
	for _, user := range paidByUsers {
		paidByUsersMap[user.ID.Hex()] = user
	}

	responses := make([]*expenseRes.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		paidByUsers := make([]*entity.Users, 0, len(expense.PaidBy))
		for _, paidByID := range expense.PaidBy {
			if user, ok := paidByUsersMap[paidByID]; ok {
				paidByUsers = append(paidByUsers, user)
			}
		}
		responses[i] = expenseMapper.ToExpenseResponse(expense, paidByUsers)
	}

	return responses, nil
}
