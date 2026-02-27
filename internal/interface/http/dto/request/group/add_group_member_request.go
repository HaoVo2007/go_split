package group

type AddGroupMemberRequest struct {
	Email string `form:"email" validate:"required,email" msg:""`
}
