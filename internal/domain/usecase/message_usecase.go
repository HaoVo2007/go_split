package usecase

import (
	"context"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"
	messageMapper "go-split/internal/interface/http/dto/mapper/message"
	messageRes "go-split/internal/interface/http/dto/response/message"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageUseCase interface {
	GetMessagesByGroupID(ctx context.Context, groupID string, pageSize int, pageIndex int) (*messageRes.ListMessageResponse, error)
}

type messageUseCase struct {
	messageRepository repository.MessageRepository
	userRepository    repository.UserRepository
}

func NewMessageUseCase(
	messageRepository repository.MessageRepository,
	userRepository repository.UserRepository,
) MessageUseCase {
	return &messageUseCase{
		messageRepository: messageRepository,
		userRepository:    userRepository,
	}
}

func (m *messageUseCase) GetMessagesByGroupID(ctx context.Context, groupID string, pageSize int, pageIndex int) (*messageRes.ListMessageResponse, error) {
	messages, totalItems, err := m.messageRepository.GetMessagesByGroupID(ctx, groupID, pageSize, pageIndex)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return &messageRes.ListMessageResponse{
			Messages:   []*messageRes.MessageResponse{},
			Pagination: &messageRes.PaginationResponse{
				PageSize:   pageSize,
				PageIndex:  pageIndex,
				TotalItems: 0,
				TotalPages: 0,
			},
		}, nil
	}

	var userIDs []primitive.ObjectID
	for _, message := range messages {
		userID, err := primitive.ObjectIDFromHex(message.UserID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	users, err := m.userRepository.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*entity.Users)
	for _, user := range users {
		userMap[user.ID.Hex()] = user
	}

	response := messageMapper.ToMessageResponse(messages, userMap)

	return &messageRes.ListMessageResponse{
		Messages: response,
		Pagination: &messageRes.PaginationResponse{
			PageSize:   pageSize,
			PageIndex:  pageIndex,
			TotalItems: int(totalItems),
			TotalPages: (int(totalItems) + pageSize - 1) / pageSize,
		},
	}, nil
	
}
