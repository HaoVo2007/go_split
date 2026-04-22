package message

import "go-split/internal/interface/http/dto/response/user"

type MessageResponse struct {
	ID        string             `json:"id"`
	GroupID   string             `json:"group_id"`
	Message   string             `json:"message"`
	Sender    *user.UserResponse `json:"sender"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}

type PaginationResponse struct {
	PageSize   int `json:"page_size"`
	PageIndex  int `json:"page_index"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type ListMessageResponse struct {
	Messages   []*MessageResponse  `json:"messages"`
	Pagination *PaginationResponse `json:"pagination"`
}
