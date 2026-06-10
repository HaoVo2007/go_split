package expense

import (
	"go-split/internal/domain/entity"
	expenseRes "go-split/internal/interface/http/dto/response/expense"
	userRes "go-split/internal/interface/http/dto/response/user"
)

func ToExpenseResponse(expense *entity.Expenses, allUsersMap map[string]*entity.Users) *expenseRes.ExpenseResponse {
	// Convert PaidBy IDs to UserResponse list
	paidByResponses := make([]*userRes.UserResponse, 0, len(expense.PaidBy))
	for _, paidByID := range expense.PaidBy {
		if user, ok := allUsersMap[paidByID]; ok {
			paidByResponses = append(paidByResponses, ToUserResponse(user))
		}
	}

	// Convert Participants IDs to UserResponse list
	participantsResponses := make([]*userRes.UserResponse, 0, len(expense.Participants))
	for _, participantID := range expense.Participants {
		if user, ok := allUsersMap[participantID]; ok {
			participantsResponses = append(participantsResponses, ToUserResponse(user))
		}
	}

	participantSplitsResponses := make([]*expenseRes.ParticipantSplit, 0, len(expense.ParticipantSplits))
	for _, split := range expense.ParticipantSplits {
		if user, ok := allUsersMap[split.UserID]; ok {
			participantSplitsResponses = append(participantSplitsResponses, &expenseRes.ParticipantSplit{
				User: ToUserResponse(user),
				Amount: split.Amount,
			})
		}
	}

	return &expenseRes.ExpenseResponse{
		ID:                expense.ID.Hex(),
		Image:             expense.Image,
		ImagePublicID:     expense.ImagePublicID,
		Name:              expense.Name,
		Amount:            expense.Amount,
		Category:          expense.Category,
		PaidBy:            paidByResponses,
		Participants:      participantsResponses,
		ParticipantSplits: participantSplitsResponses,
		CreatedAt:         expense.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         expense.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToUserResponse(user *entity.Users) *userRes.UserResponse {
	if user.Profile == nil {
		return &userRes.UserResponse{
			ID:      user.ID.Hex(),
			Email:   user.Email,
			Role:    user.Role,
			Name:    nil,
			Image:   nil,
			Address: nil,
			Phone:   nil,
		}
	}

	name := ""
	if user.Profile.Name != nil {
		name = *user.Profile.Name
	}

	image := ""
	if user.Profile.Image != nil {
		image = *user.Profile.Image
	}

	address := ""
	if user.Profile.Address != nil {
		address = *user.Profile.Address
	}

	phone := ""
	if user.Profile.Phone != nil {
		phone = *user.Profile.Phone
	}

	imagePublicID := ""
	if user.Profile.ImagePublicID != nil {
		imagePublicID = *user.Profile.ImagePublicID
	}

	return &userRes.UserResponse{
		ID:            user.ID.Hex(),
		Email:         user.Email,
		Role:          user.Role,
		Name:          &name,
		Image:         &image,
		ImagePublicID: &imagePublicID,
		Address:       &address,
		Phone:         &phone,
	}
}
