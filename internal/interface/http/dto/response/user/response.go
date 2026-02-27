package user

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Name      *string `json:"name"`
	Image     *string `json:"image"`
	Address   *string `json:"address"`
	Phone     *string `json:"phone"`
}
