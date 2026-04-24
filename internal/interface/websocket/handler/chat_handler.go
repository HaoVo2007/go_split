package handler

import (
	"go-split/internal/domain/repository"
	"go-split/internal/interface/websocket/hub"
	"go-split/pkg/libs/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatHandler struct {
	hub     *hub.Hub
	group   repository.GroupRepository
	message repository.MessageRepository
	user    repository.UserRepository
}

func NewChatHandler(
	hub *hub.Hub,
	group repository.GroupRepository,
	message repository.MessageRepository,
	user repository.UserRepository,
) *ChatHandler {
	return &ChatHandler{
		hub:     hub,
		group:   group,
		message: message,
		user:    user,
	}
}

func (h *ChatHandler) HandleConnection(c *gin.Context) {
	token := c.Query("token")

	validatedUserID, err := helper.ValidateToken(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": err.Error()})
		return
	}

	groupUser, err := h.group.GetGroups(c.Request.Context(), validatedUserID)
	if err != nil {
		log.Printf("Failed to get group user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group user", "message": err.Error()})
		return
	}

	if len(groupUser) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found", "message": "group not found"})
		return
	}

	groupIds := make(map[string]bool)
	for _, group := range groupUser {
		groupIds[group.ID.Hex()] = true
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(validatedUserID)
	if err != nil {
		log.Printf("Failed to convert user ID to object ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert user ID to object ID", "message": err.Error()})
		return
	}

	user, err := h.user.FindUserByID(c.Request.Context(), objectID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user", "message": err.Error()})
		return
	}

	client := &hub.Client{
		Hub:               h.hub,
		Conn:              conn,
		Send:              make(chan []byte, 256),
		MessageRepository: h.message,
		GroupIds:          groupIds,
		User:              user,
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
