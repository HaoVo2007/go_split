package route

import (
	"go-split/internal/interface/http/handler"
	"go-split/internal/interface/middleware"

	"github.com/gin-gonic/gin"
)

func SetupExpenseRoutes(
	router *gin.Engine,
	expenseHandler *handler.ExpenseHandler,
) {
	expenseRouter := router.Group("api/v1/expenses")
	expenseRouter.Use(middleware.AuthMiddleware())
	{
		expenseRouter.POST("/", expenseHandler.CreateExpense)
		expenseRouter.GET("/:id", expenseHandler.GetExpenseById)
		expenseRouter.PUT("/:id", expenseHandler.UpdateExpenseById)
		expenseRouter.DELETE("/:id", expenseHandler.DeleteExpenseById)
		expenseRouter.GET("/:id/settlement", expenseHandler.GetSettlementByExpenseID)
		expenseRouter.GET("/group/:group_id", expenseHandler.GetExpensesByGroupID)
		expenseRouter.GET("/group/:group_id/balance", expenseHandler.GetBalanceByGroupID)
	}
}
