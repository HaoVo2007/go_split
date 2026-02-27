package helper

import (
	"go-split/pkg/libs/contants"
	"context"
	"errors"
)

func GetUserID(ctx context.Context) (string, error) {
	userID, exists := ctx.Value(contants.ContextKeyUserID).(string)
	if !exists {
		return "", errors.New("user ID not found")
	}
	return userID, nil
}

func GetRole(ctx context.Context) (string, error) {
	role, exists := ctx.Value(contants.ContextKeyRole).(string)
	if !exists {
		return "", errors.New("role not found")
	}
	return role, nil
}
