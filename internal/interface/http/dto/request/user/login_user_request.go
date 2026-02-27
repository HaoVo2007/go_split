package user

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email" msg:""`
	Password string `json:"password" validate:"required,min=6,max=50" msg:""`
}
