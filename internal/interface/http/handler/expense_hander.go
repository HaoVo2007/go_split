package handler

import (
	"go-split/internal/domain/usecase"
	"go-split/internal/interface/http/dto/request/expense"
	"go-split/pkg/libs/response"
	"go-split/pkg/libs/validator"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	expenseUseCase usecase.ExpenseUseCase
}

func NewExpenseHandler(expenseUseCase usecase.ExpenseUseCase) *ExpenseHandler {
	return &ExpenseHandler{
		expenseUseCase: expenseUseCase,
	}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req expense.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestSimple(c, "invalid JSON")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	err := h.expenseUseCase.CreateExpense(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, "expense created successfully", nil)
}

func (h *ExpenseHandler) GetExpensesByGroupID(c *gin.Context) {
	groupID := c.Param("group_id")
	if groupID == "" {
		response.BadRequestSimple(c, "group ID is required")
		return
	}

	expenses, err := h.expenseUseCase.GetExpensesByGroupID(c.Request.Context(), groupID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "expenses fetched successfully", expenses)
}

