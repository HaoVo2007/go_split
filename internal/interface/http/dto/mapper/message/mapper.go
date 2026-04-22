package message

import (
	"go-split/internal/domain/entity"
	userMapper "go-split/internal/interface/http/dto/mapper/user"
	messageRes "go-split/internal/interface/http/dto/response/message"
)

func ToMessageResponse(messages []*entity.Messages, userMap map[string]*entity.Users) []*messageRes.MessageResponse {
	messageResponses := make([]*messageRes.MessageResponse, len(messages))
	for i, message := range messages {
		messageResponses[i] = &messageRes.MessageResponse{
			ID:        message.ID.Hex(),
			GroupID:   message.GroupID,
			Message:   message.Message,
			Sender:    userMapper.ToUserResponse(userMap[message.UserID]),
			CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: message.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return messageResponses
}
