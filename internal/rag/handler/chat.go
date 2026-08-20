package handler

import (
	"go-split/internal/rag/service"
	"go-split/pkg/libs/response"

	"github.com/gin-gonic/gin"
)

type ChatAIHandler struct {
	ChatService service.ChatService
}

func NewChatAIHandler(chatService service.ChatService) *ChatAIHandler {
	return &ChatAIHandler{
		ChatService: chatService,
	}
}

func (h *ChatAIHandler) ChatAIMessage(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestSimple(c, "invalid data")
		return
	}

	data, err := h.ChatService.ChatAIMessage(c.Request.Context(), req.Message)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "chat with AI successfully", data)
}
