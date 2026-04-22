package websocket

import (
	"go-split/internal/interface/websocket/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, chatHandler *handler.ChatHandler) {
	router.GET("/ws", chatHandler.HandleConnection)
}
