package group

import "mime/multipart"

type UpdateGroupRequest struct {
	Name        string                `form:"name" validate:"min=3,max=100" msg:""`
	Description string                `form:"description" validate:"min=3,max=100" msg:""`
	Image       *multipart.FileHeader `form:"image" validate:"" msg:""`
}
