package expense

import (
	"go-split/internal/domain/entity"
	expenseRes "go-split/internal/interface/http/dto/response/expense"
	userRes "go-split/internal/interface/http/dto/response/user"
)

func ToExpenseResponse(expense *entity.Expenses, paidBy []*entity.Users) *expenseRes.ExpenseResponse {
	paidByResponses := make([]*userRes.UserResponse, len(paidBy))
	for i, user := range paidBy {
		if user.Profile == nil {
			paidByResponses[i] = &userRes.UserResponse{
				ID:      user.ID.Hex(),
				Email:   user.Email,
				Role:    user.Role,
				Name:    nil,
				Image:   nil,
				Address: nil,
				Phone:   nil,
			}
			continue
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

		paidByResponses[i] = &userRes.UserResponse{
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
	return &expenseRes.ExpenseResponse{
		ID:            expense.ID.Hex(),
		Image:         expense.Image,
		ImagePublicID: expense.ImagePublicID,
		Name:          expense.Name,
		Amount:        expense.Amount,
		Category:      expense.Category,
		PaidBy:        paidByResponses,
		CreatedAt:     expense.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     expense.UpdatedAt.Format("2006-01-02 15:04:05"),
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
