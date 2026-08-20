package http

import (
	"go-split/internal/interface/http/handler"
	"go-split/internal/interface/http/route"
	ragHandler "go-split/internal/rag/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	groupHandler *handler.GroupHandler,
	messageHandler *handler.MessageHandler,
	expenseHandler *handler.ExpenseHandler,
	chatAIHandler *ragHandler.ChatAIHandler,
) {
	route.SetupUserRoutes(router, userHandler)
	route.SetupGroupRoutes(router, groupHandler)
	route.SetupExpenseRoutes(router, expenseHandler)
	route.SetupMessageRoutes(router, messageHandler)
	route.SetupChatAIRoutes(router, chatAIHandler)
}
