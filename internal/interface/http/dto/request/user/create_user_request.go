package user

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email" msg:""`
	Password string `json:"password" validate:"required,min=6,max=50" msg:""`
	Role     string `json:"role" validate:"required,oneof=admin staff user" msg:""`
}
