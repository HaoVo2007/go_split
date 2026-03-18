package user

import "mime/multipart"

type UpdateProfileRequest struct {
	Name      string                `form:"name" validate:"omitempty,min=3,max=100" msg:""`
	Image     *multipart.FileHeader `form:"image" validate:"" msg:""`
	Address   string                `form:"address" validate:"omitempty,min=3,max=100" msg:""`
	Phone     string                `form:"phone" validate:"omitempty,min=3,max=100" msg:""`
	ImageType string                `form:"image_type" validate:"omitempty,oneof=upload preset" msg:""`
	ImageUrl  string                `form:"image_url" validate:"omitempty,url" msg:""`
}
