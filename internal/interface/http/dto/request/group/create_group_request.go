package group

import "mime/multipart"

type CreateGroupRequest struct {
	Name        string                `form:"name" validate:"required,min=3,max=100" msg:""`
	Description string                `form:"description" msg:""`
	Image       *multipart.FileHeader `form:"image" msg:""`
}
