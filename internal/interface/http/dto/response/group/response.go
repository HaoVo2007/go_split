package group

import "go-split/internal/interface/http/dto/response/user"

type GroupResponse struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Image         string               `json:"image"`
	ImagePublicID string               `json:"image_public_id"`
	Members       []*user.UserResponse `json:"members"`
	TotalMembers  int                  `json:"total_members"`
	Description   string               `json:"description"`
	CreatedBy     string               `json:"created_by"`
	CreatedAt     string               `json:"created_at"`
	UpdatedAt     string               `json:"updated_at"`
}
