package handler

import (
	"go-split/internal/domain/usecase"
	"go-split/pkg/libs/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	MessageUseCase usecase.MessageUseCase
}

func NewMessageHandler(messageUseCase usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{
		MessageUseCase: messageUseCase,
	}
}

func (h *MessageHandler) GetMessagesByGroupID(c *gin.Context) {
	groupID := c.Param("group_id")
	pageSize := c.Query("page_size")
	pageIndex := c.Query("page_index")

	if pageSize == "" {
		pageSize = "10"
	}

	if pageIndex == "" {
		pageIndex = "1"
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		response.BadRequestSimple(c, "invalid page size")
		return
	}

	pageIndexInt, err := strconv.Atoi(pageIndex)
	if err != nil {
		response.BadRequestSimple(c, "invalid page index")
		return
	}

	messages, err := h.MessageUseCase.GetMessagesByGroupID(c.Request.Context(), groupID, pageSizeInt, pageIndexInt)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "messages fetched successfully", messages)
}
