package expense

type CreateExpenseRequest struct {
	GroupID  string   `json:"group_id" validate:"required" msg:""`
	Date     string   `json:"date" validate:"required" msg:""`
	Name     string   `json:"name" validate:"required" msg:""`
	Amount   float64  `json:"amount" validate:"required" msg:""`
	Category string   `json:"category" validate:"required" msg:""`
	PaidBy   []string `json:"paid_by" validate:"required" msg:""`
}
