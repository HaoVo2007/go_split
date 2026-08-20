package route

import (
	"go-split/internal/interface/middleware"
	"go-split/internal/rag/handler"

	"github.com/gin-gonic/gin"
)

func SetupChatAIRoutes(router *gin.Engine, chatAIHandler *handler.ChatAIHandler) {
	chatAIRouter := router.Group("api/v1/chat-ai")
	chatAIRouter.Use(middleware.AuthMiddleware())
	{
		chatAIRouter.POST("", chatAIHandler.ChatAIMessage)
	}
}
