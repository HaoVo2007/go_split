package handler

import (
	"go-split/internal/domain/usecase"
	"go-split/internal/interface/http/dto/request/group"
	"go-split/pkg/libs/response"
	"go-split/pkg/libs/validator"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupUseCase usecase.GroupUseCase
}

func NewGroupHandler(groupUseCase usecase.GroupUseCase) *GroupHandler {
	return &GroupHandler{
		groupUseCase: groupUseCase,
	}
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req group.CreateGroupRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	err := h.groupUseCase.CreateGroup(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, "group created successfully", nil)
}

func (h *GroupHandler) GetGroups(c *gin.Context) {
	groups, err := h.groupUseCase.GetGroups(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "groups fetched successfully", groups)
}

func (h *GroupHandler) GetGroupById(c *gin.Context) {
	id := c.Param("id")
	group, err := h.groupUseCase.GetGroupById(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "group fetched successfully", group)
}

func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var req group.UpdateGroupRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	err := h.groupUseCase.UpdateGroup(c.Request.Context(), id, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "group updated successfully", nil)
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	err := h.groupUseCase.DeleteGroup(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "group deleted successfully", nil)
}

func (h *GroupHandler) AddGroupMember(c *gin.Context) {
	id := c.Param("id")
	var req group.AddGroupMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	err := h.groupUseCase.AddGroupMember(c.Request.Context(), id, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "group member added successfully", nil)
}

func (h *GroupHandler) RemoveGroupMember(c *gin.Context) {
	id := c.Param("id")
	memberId := c.Param("memberId")
	err := h.groupUseCase.RemoveGroupMember(c.Request.Context(), id, memberId)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "group member removed successfully", nil)
}
