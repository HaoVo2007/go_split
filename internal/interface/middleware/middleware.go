package middleware

import (
	"go-split/pkg/libs/contants"
	"go-split/pkg/libs/response"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Panic("JWT_SECRET not set")
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
					response.Unauthorized(c, "Token expired")
					c.Abort()
					return
				}
			}
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		tokenType, ok := token.Claims.(jwt.MapClaims)["type"].(string)
		if !ok {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}
		if tokenType != "access" {
			response.Unauthorized(c, "Invalid token type")
			c.Abort()
			return
		}

		userID, ok := token.Claims.(jwt.MapClaims)["user_id"].(string)
		if !ok {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		role, ok := token.Claims.(jwt.MapClaims)["role"].(string)
		if !ok {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		ctx := context.WithValue(
			c.Request.Context(),
			contants.ContextKeyUserID,
			userID,
		)

		ctx = context.WithValue(
			ctx,
			contants.ContextKeyToken,
			tokenString,
		)

		ctx = context.WithValue(
			ctx,
			contants.ContextKeyRole,
			role,
		)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Request.Context().Value(contants.ContextKeyRole).(string)
		if !ok {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		if role != "admin" {
			response.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
