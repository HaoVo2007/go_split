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
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequestSimple(c, "invalid data")
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

func (h *ExpenseHandler) GetExpenseById(c *gin.Context) {
	expenseID := c.Param("id")
	if expenseID == "" {
		response.BadRequestSimple(c, "expense ID is required")
		return
	}

	expense, err := h.expenseUseCase.GetExpenseById(c.Request.Context(), expenseID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "expense fetched successfully", expense)
}

func (h *ExpenseHandler) UpdateExpenseById(c *gin.Context) {
	expenseID := c.Param("id")
	if expenseID == "" {
		response.BadRequestSimple(c, "expense ID is required")
		return
	}

	var req expense.UpdateExpenseRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequestSimple(c, "invalid data")
		return
	}

	customValidator := validator.NewCustomValidator()
	if validationErrors := customValidator.ValidateAndTranslate(&req); validationErrors != nil {
		response.BadRequest(c, "invalid data", validationErrors)
		return
	}

	err := h.expenseUseCase.UpdateExpenseById(c.Request.Context(), expenseID, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "expense updated successfully", nil)
}

func (h *ExpenseHandler) DeleteExpenseById(c *gin.Context) {
	expenseID := c.Param("id")
	if expenseID == "" {
		response.BadRequestSimple(c, "expense ID is required")
		return
	}

	err := h.expenseUseCase.DeleteExpenseById(c.Request.Context(), expenseID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "expense deleted successfully", nil)
}

func (h *ExpenseHandler) GetSettlementByExpenseID(c *gin.Context) {
	expenseID := c.Param("id")
	if expenseID == "" {
		response.BadRequestSimple(c, "expense ID is required")
		return
	}

	settlement, err := h.expenseUseCase.GetSettlementByExpenseID(c.Request.Context(), expenseID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "settlement fetched successfully", settlement)
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

func (h *ExpenseHandler) GetBalanceByGroupID(c *gin.Context) {
	groupID := c.Param("group_id")
	if groupID == "" {
		response.BadRequestSimple(c, "group ID is required")
		return
	}

	balance, err := h.expenseUseCase.GetBalanceByGroupID(c.Request.Context(), groupID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "balance fetched successfully", balance)
}