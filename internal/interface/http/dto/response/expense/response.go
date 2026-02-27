package expense

import "go-split/internal/interface/http/dto/response/user"

type ExpenseResponse struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Amount    float64              `json:"amount"`
	Category  string               `json:"category"`
	PaidBy    []*user.UserResponse `json:"paid_by"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}
