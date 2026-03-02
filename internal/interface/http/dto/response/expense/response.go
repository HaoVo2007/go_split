package expense

import "go-split/internal/interface/http/dto/response/user"

type ExpenseResponse struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Amount        float64              `json:"amount"`
	Category      string               `json:"category"`
	Image         string               `json:"image"`
	ImagePublicID string               `json:"image_public_id"`
	PaidBy        []*user.UserResponse `json:"paid_by"`
	CreatedAt     string               `json:"created_at"`
	UpdatedAt     string               `json:"updated_at"`
}

type BalanceResponse struct {
	GroupID       string           `json:"group_id"`
	TotalExpenses float64          `json:"total_expenses"`
	TotalOwed     float64          `json:"total_owed"`
	TotalPaid     float64          `json:"total_paid"`
	Member        []*MemberBalance `json:"members"`
	Settlement    []*Settlement    `json:"settlements"`
}

type MemberBalance struct {
	UserID    string  `json:"user_id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	Email     string  `json:"email"`
	TotalPaid float64 `json:"total_paid"`
	TotalOwed float64 `json:"total_owed"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
}

type Settlement struct {
	FromUser user.UserResponse `json:"from_user"`
	ToUser   user.UserResponse `json:"to_user"`
	Amount   float64           `json:"amount"`
}

type SettlementResponse struct {
	ExpenseID   string           `json:"expense_id"`
	Amount      float64          `json:"amount"`
	Members     []*MemberBalance `json:"members"`
	Settlements []*Settlement    `json:"settlements"`
}
