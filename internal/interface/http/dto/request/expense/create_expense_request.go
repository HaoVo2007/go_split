package expense

import "mime/multipart"

type CreateExpenseRequest struct {
	GroupID      string                `form:"group_id" validate:"required" msg:""`
	Image        *multipart.FileHeader `form:"image" msg:""`
	Date         string                `form:"date" validate:"required" msg:""`
	Name         string                `form:"name" validate:"required" msg:""`
	Amount       float64               `form:"amount" validate:"required" msg:""`
	Category     string                `form:"category" validate:"required" msg:""`
	PaidBy       []string              `form:"paid_by" validate:"required" msg:""`
	Participants []string              `form:"participants" validate:"required" msg:""`
}

type UpdateExpenseRequest struct {
	Image    *multipart.FileHeader `form:"image" msg:""`
	Date     string                `form:"date" validate:"required" msg:""`
	Name     string                `form:"name" validate:"required" msg:""`
	Amount   float64               `form:"amount" validate:"required" msg:""`
	Category string                `form:"category" validate:"required" msg:""`
	PaidBy   []string              `form:"paid_by" validate:"required" msg:""`
}
