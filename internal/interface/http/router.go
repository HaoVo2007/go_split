package http

import (
	"go-split/internal/interface/http/handler"
	"go-split/internal/interface/http/route"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	groupHandler *handler.GroupHandler,
	expenseHandler *handler.ExpenseHandler,
) {
	route.SetupUserRoutes(router, userHandler)
	route.SetupGroupRoutes(router, groupHandler)
	route.SetupExpenseRoutes(router, expenseHandler)
}
