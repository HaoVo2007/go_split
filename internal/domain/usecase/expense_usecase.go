package usecase

import (
	"context"
	"errors"
	"fmt"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"
	expenseMapper "go-split/internal/interface/http/dto/mapper/expense"
	userMapper "go-split/internal/interface/http/dto/mapper/user"
	"go-split/internal/interface/http/dto/request/expense"
	expenseRes "go-split/internal/interface/http/dto/response/expense"
	"go-split/pkg/libs/helper"
	"math"
	"os"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseUseCase interface {
	CreateExpense(ctx context.Context, req expense.CreateExpenseRequest) error
	GetExpenseById(ctx context.Context, expenseID string) (*expenseRes.ExpenseResponse, error)
	UpdateExpenseById(ctx context.Context, expenseID string, req expense.UpdateExpenseRequest) error
	DeleteExpenseById(ctx context.Context, expenseID string) error
	GetSettlementByExpenseID(ctx context.Context, expenseID string) (*expenseRes.SettlementResponse, error)
	GetExpensesByGroupID(ctx context.Context, groupID string) ([]*expenseRes.ExpenseResponse, error)
	GetBalanceByGroupID(ctx context.Context, groupID string) (*expenseRes.BalanceResponse, error)
}

type expenseUseCase struct {
	expenseRepository      repository.ExpenseRepository
	expenseSplitRepository repository.ExpenseSplitRepository
	groupRepository        repository.GroupRepository
	userRepository         repository.UserRepository
	cloudinaryUploader     *helper.CloudinaryUploader
}

func NewExpenseUseCase(
	expenseRepository repository.ExpenseRepository,
	expenseSplitRepository repository.ExpenseSplitRepository,
	groupRepository repository.GroupRepository,
	userRepository repository.UserRepository,
	cloudinaryUploader *helper.CloudinaryUploader,
) ExpenseUseCase {
	return &expenseUseCase{
		expenseRepository:      expenseRepository,
		expenseSplitRepository: expenseSplitRepository,
		groupRepository:        groupRepository,
		userRepository:         userRepository,
		cloudinaryUploader:     cloudinaryUploader,
	}
}

func (e *expenseUseCase) CreateExpense(ctx context.Context, req expense.CreateExpenseRequest) error {
	var imageURL string
	var imagePublicID string
	if req.Image != nil {
		tempPath := fmt.Sprintf("temp_%s", req.Image.Filename)
		if err := helper.SaveUploadedFile(req.Image, tempPath); err != nil {
			return err
		}
		defer os.Remove(tempPath)

		image, publicID, err := e.cloudinaryUploader.UploadImage(ctx, tempPath, "expenses")
		if err != nil {
			return err
		}
		imageURL = image
		imagePublicID = publicID
	}

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

	expenseID := primitive.NewObjectID()
	expense := entity.Expenses{
		ID:            expenseID,
		GroupID:       req.GroupID,
		Image:         imageURL,
		ImagePublicID: imagePublicID,
		Date:          date,
		Name:          req.Name,
		Amount:        req.Amount,
		Category:      req.Category,
		PaidBy:        req.PaidBy,
		Participants:  req.Participants,
		IsDeleted:     false,
		DeletedAt:     nil,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = e.expenseRepository.CreateExpense(ctx, expense)
	if err != nil {
		return err
	}

	share := req.Amount / float64(len(req.Participants))
	expenseSplits := make([]entity.ExpenseSplits, len(req.Participants))
	for i, participant := range req.Participants {
		expenseSplits[i] = entity.ExpenseSplits{
			ID:         primitive.NewObjectID(),
			ExpensesID: expenseID.Hex(),
			UserId:     participant,
			Amount:     share,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	err = e.expenseSplitRepository.CreateExpenseSplits(ctx, expenseSplits)
	if err != nil {
		return err
	}

	return nil

}

func (e *expenseUseCase) GetExpenseById(ctx context.Context, expenseID string) (*expenseRes.ExpenseResponse, error) {
	expenseIDObject, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return nil, err
	}

	expense, err := e.expenseRepository.GetExpenseById(ctx, expenseIDObject)
	if err != nil {
		return nil, err
	}

	if expense == nil {
		return nil, errors.New("expense not found")
	}

	paidByIDs := []primitive.ObjectID{}
	for _, paidByID := range expense.PaidBy {
		paidByIDObject, err := primitive.ObjectIDFromHex(paidByID)
		if err != nil {
			return nil, err
		}
		paidByIDs = append(paidByIDs, paidByIDObject)
	}

	paidByUsers, err := e.userRepository.GetUsersByIDs(ctx, paidByIDs)
	if err != nil {
		return nil, err
	}

	paidByUsersMap := make(map[string]*entity.Users)
	for _, user := range paidByUsers {
		paidByUsersMap[user.ID.Hex()] = user
	}

	paidByUsersEntities := make([]*entity.Users, 0, len(expense.PaidBy))
	for _, paidByID := range expense.PaidBy {
		if user, ok := paidByUsersMap[paidByID]; ok {
			paidByUsersEntities = append(paidByUsersEntities, user)
		}
	}

	return expenseMapper.ToExpenseResponse(expense, paidByUsersEntities), nil
}

func (e *expenseUseCase) UpdateExpenseById(ctx context.Context, expenseID string, req expense.UpdateExpenseRequest) error {
	expenseIDObject, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return err
	}

	expense, err := e.expenseRepository.GetExpenseById(ctx, expenseIDObject)
	if err != nil {
		return err
	}

	if expense == nil {
		return errors.New("expense not found")
	}

	if req.Image != nil {
		if expense.ImagePublicID != "" {
			e.cloudinaryUploader.DeleteImage(ctx, expense.ImagePublicID)
		}
		tempPath := fmt.Sprintf("temp_%s", req.Image.Filename)
		if err := helper.SaveUploadedFile(req.Image, tempPath); err != nil {
			return err
		}
		defer os.Remove(tempPath)
		image, publicID, err := e.cloudinaryUploader.UploadImage(ctx, tempPath, "expenses")
		if err != nil {
			return err
		}
		expense.Image = image
		expense.ImagePublicID = publicID
	}

	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return err
		}
		expense.Date = date
	}

	if req.Name != "" {
		expense.Name = req.Name
	}

	if req.Amount != 0 {
		expense.Amount = req.Amount
	}

	if req.Category != "" {
		expense.Category = req.Category
	}

	if len(req.PaidBy) > 0 {
		expense.PaidBy = req.PaidBy
	}

	if len(req.Participants) > 0 {
		err = e.expenseSplitRepository.DeleteExpenseSplitsByExpenseID(ctx, expenseID)
		if err != nil {
			return err
		}

		share := req.Amount / float64(len(req.Participants))
		expenseSplits := make([]entity.ExpenseSplits, len(req.Participants))
		for i, participant := range req.Participants {
			expenseSplits[i] = entity.ExpenseSplits{
				ID:         primitive.NewObjectID(),
				ExpensesID: expenseIDObject.Hex(),
				UserId:     participant,
				Amount:     share,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
		}

		err = e.expenseSplitRepository.CreateExpenseSplits(ctx, expenseSplits)
		if err != nil {
			return err
		}

		expense.Participants = req.Participants
	}

	expense.UpdatedAt = time.Now()

	return e.expenseRepository.UpdateExpenseById(ctx, expenseIDObject, expense)
}

func (e *expenseUseCase) DeleteExpenseById(ctx context.Context, expenseID string) error {
	expenseIDObject, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return err
	}

	expense, err := e.expenseRepository.GetExpenseById(ctx, expenseIDObject)
	if err != nil {
		return err
	}

	if expense == nil {
		return errors.New("expense not found")
	}

	if expense.ImagePublicID != "" {
		e.cloudinaryUploader.DeleteImage(ctx, expense.ImagePublicID)
	}

	deletedAt := time.Now()

	expense.IsDeleted = true
	expense.DeletedAt = &deletedAt
	expense.UpdatedAt = time.Now()

	err = e.expenseRepository.UpdateExpenseById(ctx, expenseIDObject, expense)
	if err != nil {
		return err
	}

	err = e.expenseSplitRepository.DeleteExpenseSplitsByExpenseID(ctx, expenseID)
	if err != nil {
		return err
	}

	return nil
}

func (e *expenseUseCase) GetSettlementByExpenseID(ctx context.Context, expenseID string) (*expenseRes.SettlementResponse, error) {
	expenseIDObject, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return nil, err
	}

	expense, err := e.expenseRepository.GetExpenseById(ctx, expenseIDObject)
	if err != nil {
		return nil, err
	}

	if expense == nil {
		return nil, errors.New("expense not found")
	}

	splits, err := e.expenseSplitRepository.GetExpenseSplitsByExpenseID(ctx, expenseID)
	if err != nil {
		return nil, err
	}

	if splits == nil {
		return nil, errors.New("splits not found")
	}

	userIDSet := make(map[string]bool)
	for _, userID := range expense.PaidBy {
		userIDSet[userID] = true
	}
	for _, split := range splits {
		userIDSet[split.UserId] = true
	}

	userIDs := make([]primitive.ObjectID, 0, len(userIDSet))
	for userID := range userIDSet {
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, objID)
	}

	users, err := e.userRepository.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*entity.Users)
	for _, u := range users {
		userMap[u.ID.Hex()] = u
	}

	paidByMap := make(map[string]bool)
	for _, userID := range expense.PaidBy {
		paidByMap[userID] = true
	}

	type tempBalance struct {
		UserID  string
		Balance float64
	}

	balanceMap := make(map[string]float64)

	sharePaid := expense.Amount / float64(len(expense.PaidBy))
	for _, paidByID := range expense.PaidBy {
		balanceMap[paidByID] += sharePaid
	}

	for _, split := range splits {
		balanceMap[split.UserId] -= split.Amount
	}

	var creditors []*tempBalance
	var debtors []*tempBalance

	for userID, balance := range balanceMap {
		if balance > 0 {
			creditors = append(creditors, &tempBalance{UserID: userID, Balance: balance})
		} else {
			debtors = append(debtors, &tempBalance{UserID: userID, Balance: balance})
		}
	}

	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].Balance > creditors[j].Balance
	})
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].Balance < debtors[j].Balance
	})

	var members []*expenseRes.MemberBalance

	for userID, balance := range balanceMap {
		user, ok := userMap[userID]
		if !ok {
			return nil, errors.New("user not found")
		}

		name := user.Email
		image := ""
		if user.Profile != nil {
			if user.Profile.Name != nil {
				name = *user.Profile.Name
			}
			if user.Profile.Image != nil {
				image = *user.Profile.Image
			}
		}

		sharePaid := expense.Amount / float64(len(expense.PaidBy))
		totalPaid := 0.0
		if paidByMap[userID] {
			totalPaid = sharePaid
		}

		totalOwed := 0.0
		for _, split := range splits {
			if split.UserId == userID {
				totalOwed = split.Amount
				break
			}
		}

		member := &expenseRes.MemberBalance{
			UserID:    userID,
			Name:      name,
			Image:     image,
			Email:     user.Email,
			TotalPaid: totalPaid,
			TotalOwed: totalOwed,
			Balance:   balance,
		}

		switch {
		case balance > 0:
			member.Status = "creditor"
		case balance < 0:
			member.Status = "debtor"
		default:
			member.Status = "settled"
		}

		members = append(members, member)
	}

	settlements := make([]*expenseRes.Settlement, 0, len(splits))

	i, j := 0, 0
	for i < len(creditors) && j < len(debtors) {
		creditor := creditors[i]
		debtor := debtors[j]

		amount := math.Min(creditor.Balance, -debtor.Balance)
		fromUser := userMapper.ToUserResponse(userMap[creditor.UserID])
		toUser := userMapper.ToUserResponse(userMap[debtor.UserID])
		settlements = append(settlements, &expenseRes.Settlement{
			FromUser: *toUser,
			ToUser:   *fromUser,
			Amount:   amount,
		})

		creditor.Balance -= amount
		debtor.Balance += amount

		if creditor.Balance == 0 {
			i++
		}
		if debtor.Balance == 0 {
			j++
		}
	}

	return &expenseRes.SettlementResponse{
		ExpenseID:   expenseID,
		Amount:      expense.Amount,
		Members:     members,
		Settlements: settlements,
	}, nil
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

func (e *expenseUseCase) GetBalanceByGroupID(ctx context.Context, groupID string) (*expenseRes.BalanceResponse, error) {
	groupObjID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, err
	}

	group, err := e.groupRepository.GetGroupById(ctx, groupObjID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}

	expenses, err := e.expenseRepository.GetExpensesByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	totalExpenses := 0.0
	expenseIDs := make([]string, 0, len(expenses))

	for _, exp := range expenses {
		totalExpenses += exp.Amount
		expenseIDs = append(expenseIDs, exp.ID.Hex())
	}

	splits, err := e.expenseSplitRepository.GetExpenseSplitsByExpenseIDs(ctx, expenseIDs)
	if err != nil {
		return nil, err
	}

	splitsMap := make(map[string][]*entity.ExpenseSplits)
	for _, s := range splits {
		splitsMap[s.ExpensesID] = append(splitsMap[s.ExpensesID], s)
	}

	type balanceInfo struct {
		Paid float64
		Owed float64
	}

	balanceMap := make(map[string]*balanceInfo)

	for _, memberID := range group.Members {
		balanceMap[memberID] = &balanceInfo{}
	}

	if _, ok := balanceMap[group.CreatedBy]; !ok {
		balanceMap[group.CreatedBy] = &balanceInfo{}
	}

	for _, exp := range expenses {
		share := exp.Amount / float64(len(exp.PaidBy))

		for _, payer := range exp.PaidBy {
			balanceMap[payer].Paid += share
		}

		for _, split := range splitsMap[exp.ID.Hex()] {
			balanceMap[split.UserId].Owed += split.Amount
		}
	}

	userIDs := make([]primitive.ObjectID, 0, len(balanceMap))

	for id := range balanceMap {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, objID)
	}

	users, err := e.userRepository.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*entity.Users)
	for _, u := range users {
		userMap[u.ID.Hex()] = u
	}

	currentUserID, err := helper.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	var (
		members   []*expenseRes.MemberBalance
		creditors []*expenseRes.MemberBalance
		debtors   []*expenseRes.MemberBalance
		totalPaid float64
		totalOwed float64
	)

	for userID, info := range balanceMap {

		user, ok := userMap[userID]
		if !ok {
			return nil, errors.New("user not found")
		}

		name := user.Email
		image := ""
		if user.Profile != nil {
			if user.Profile.Name != nil {
				name = *user.Profile.Name
			}
			if user.Profile.Image != nil {
				image = *user.Profile.Image
			}
		}

		balance := info.Paid - info.Owed

		if userID == currentUserID {
			totalPaid = info.Paid
			totalOwed = info.Owed
		}

		member := &expenseRes.MemberBalance{
			UserID:    userID,
			Name:      name,
			Image:     image,
			Email:     user.Email,
			TotalPaid: info.Paid,
			TotalOwed: info.Owed,
			Balance:   balance,
		}

		switch {
		case balance > 0:
			member.Status = "creditor"
			creditors = append(creditors, member)
		case balance < 0:
			member.Status = "debtor"
			debtors = append(debtors, member)
		default:
			member.Status = "settled"
		}

		members = append(members, member)
	}

	type tempBalance struct {
		UserID  string
		Balance float64
	}

	var creditorTemps []tempBalance
	var debtorTemps []tempBalance

	for _, c := range creditors {
		creditorTemps = append(creditorTemps, tempBalance{c.UserID, c.Balance})
	}
	for _, d := range debtors {
		debtorTemps = append(debtorTemps, tempBalance{d.UserID, d.Balance})
	}

	sort.Slice(creditorTemps, func(i, j int) bool {
		return creditorTemps[i].Balance > creditorTemps[j].Balance
	})
	sort.Slice(debtorTemps, func(i, j int) bool {
		return debtorTemps[i].Balance < debtorTemps[j].Balance
	})

	var settlements []*expenseRes.Settlement

	i, j := 0, 0
	for i < len(creditorTemps) && j < len(debtorTemps) {

		creditor := &creditorTemps[i]
		debtor := &debtorTemps[j]

		amount := math.Min(creditor.Balance, -debtor.Balance)

		fromUser := userMapper.ToUserResponse(userMap[creditor.UserID])
		toUser := userMapper.ToUserResponse(userMap[debtor.UserID])

		settlements = append(settlements, &expenseRes.Settlement{
			FromUser: *toUser,
			ToUser:   *fromUser,
			Amount:   amount,
		})

		creditor.Balance -= amount
		debtor.Balance += amount

		if creditor.Balance == 0 {
			i++
		}
		if debtor.Balance == 0 {
			j++
		}
	}

	return &expenseRes.BalanceResponse{
		GroupID:       groupID,
		TotalExpenses: totalExpenses,
		Member:        members,
		TotalPaid:     totalPaid,
		TotalOwed:     totalOwed,
		Settlement:    settlements,
	}, nil
}
