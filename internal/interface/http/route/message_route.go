package route

import (
	"go-split/internal/interface/http/handler"
	"go-split/internal/interface/middleware"

	"github.com/gin-gonic/gin"
)

func SetupMessageRoutes(
	router *gin.Engine,
	messageHandler *handler.MessageHandler,
) {
	messageRouter := router.Group("api/v1/messages")
	messageRouter.Use(middleware.AuthMiddleware())
	{
		messageRouter.GET("/:group_id", messageHandler.GetMessagesByGroupID)
	}
}
