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
		expenseRouter.GET("/:group_id", expenseHandler.GetExpensesByGroupID)
	}
}
