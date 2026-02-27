package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidationError chứa thông tin validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse chứa thông tin error response
type ErrorResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Errors  []ValidationError  `json:"errors,omitempty"`
	Data    interface{}        `json:"data,omitempty"`
}

// SuccessResponse chứa thông tin success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// BadRequest returns 400 Bad Request with validation errors
func BadRequest(c *gin.Context, message string, validationErrors map[string]string) {
	errors := make([]ValidationError, 0, len(validationErrors))
	for field, msg := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   field,
			Message: msg,
		})
	}

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// BadRequestSimple returns 400 Bad Request with simple message
func BadRequestSimple(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// InternalServerError returns 500 Internal Server Error
func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

// NotFound returns 404 Not Found
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// Unauthorized returns 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// Forbidden returns 403 Forbidden
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// Conflict returns 409 Conflict
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// Success returns 200 OK with data
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created returns 201 Created with data
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent trả về 204 No Content
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
