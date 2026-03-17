package handler

import (
	"go-split/internal/domain/usecase"
	"go-split/internal/interface/http/dto/request/user"
	"go-split/pkg/libs/helper"
	"go-split/pkg/libs/response"
	"go-split/pkg/libs/validator"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	err := h.userUseCase.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, "user created successfully", nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	user, err := h.userUseCase.LoginUser(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "login successfully", user)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req user.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	user, err := h.userUseCase.RegisterUser(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, "user registered successfully", user)
}

func (h *UserHandler) Logout(c *gin.Context) {
	err := h.userUseCase.LogoutUser(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "logout successfully", nil)
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, err := helper.GetUserID(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	user, err := h.userUseCase.GetUserById(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "get current user successfully", user)
}

func (h *UserHandler) 	UpdateProfile(c *gin.Context) {
	var req user.UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	err := h.userUseCase.UpdateProfile(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "update profile successfully", nil)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userUseCase.GetUsers(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "get users successfully", users)
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequestSimple(c, "user ID is required")
		return
	}

	user, err := h.userUseCase.GetUserById(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "get user by ID successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequestSimple(c, "user ID is required")
		return
	}

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	err := h.userUseCase.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "update user successfully", nil)
}

func (h *UserHandler) SoftDeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequestSimple(c, "user ID is required")
		return
	}

	err := h.userUseCase.SoftDeleteUser(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "soft delete user successfully", nil)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		response.BadRequestSimple(c, "refresh token is required")
		return
	}

	accessToken, refreshToken, err := h.userUseCase.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "refresh token successfully", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
