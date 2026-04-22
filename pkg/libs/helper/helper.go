package helper

import (
	"context"
	"errors"
	"fmt"
	"go-split/pkg/libs/contants"
	"os"

	"github.com/golang-jwt/jwt/v4"
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

func ValidateToken(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", errors.New("token expired")
			}
		}
		return "", errors.New("invalid token")
	}

	tokenType, ok := token.Claims.(jwt.MapClaims)["type"].(string)
	if !ok || tokenType != "access" {
		return "", errors.New("invalid token type")
	}

	userID, ok := token.Claims.(jwt.MapClaims)["user_id"].(string)
	if !ok {
		return "", errors.New("user ID not found in token")
	}

	return userID, nil
}