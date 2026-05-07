package hub

import (
	"context"
	"encoding/json"
	"go-split/internal/domain/entity"
	"go-split/internal/domain/repository"
	"go-split/internal/interface/websocket/event"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

type Client struct {
	Hub               *Hub
	Conn              *websocket.Conn
	Send              chan []byte
	GroupIds          map[string]bool
	MessageRepository repository.MessageRepository
	User              *entity.Users
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var base struct {
			TypeMessage string `json:"type_message"`
		}

		if err := json.Unmarshal(data, &base); err != nil {
			log.Println(err)
			continue
		}

		switch base.TypeMessage {
		case "message":
			var msg event.MessageRequest
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("error: %v", err)
				break
			}
			c.handleMessage(msg)
		case "seen":
			var seen event.SeenRequest
			if err := json.Unmarshal(data, &seen); err != nil {
				log.Printf("error: %v", err)
				break
			}
			c.handleSeen(seen)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(incoming event.MessageRequest) {
	if !c.GroupIds[incoming.GroupID] {
		log.Println("User send message to group not in list of groups")
		return
	}

	var senderName string
	if c.User.Profile != nil && c.User.Profile.Name != nil {
		senderName = *c.User.Profile.Name
	} else {
		senderName = c.User.Email
	}

	var senderAvatar string
	if c.User.Profile != nil && c.User.Profile.Image != nil {
		senderAvatar = *c.User.Profile.Image
	} else {
		senderAvatar = ""
	}

	message := entity.Messages{
		ID:      primitive.NewObjectID(),
		GroupID: incoming.GroupID,
		Message: incoming.Message,
		UserID:  c.User.ID.Hex(),
		SeenBy: []entity.SeenByUser{
			{
				UserID: c.User.ID.Hex(),
				SeenAt: time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.MessageRepository.CreateMessage(context.Background(), message); err != nil {
		log.Printf("error: %v", err)
		return
	}

	c.Hub.MessageBroadcasts <- &event.MessageEvent{
		TypeMessage: "message",
		GroupID:     incoming.GroupID,
		Message:     incoming.Message,
		SenderID:    c.User.ID.Hex(),
		SenderName:  senderName,
		Avatar:      senderAvatar,
		CreatedAt:   time.Now(),
	}

	c.Hub.UnreadBroadcasts <- &event.UnreadEvent{
		TypeMessage: "unread_count",
		GroupID:     incoming.GroupID,
		SenderID:    c.User.ID.Hex(),
	}
}

func (c *Client) handleSeen(seen event.SeenRequest) {
	userID := c.User.ID.Hex()

	if !c.GroupIds[seen.GroupID] {
		log.Println("User send seen to group not in list of groups")
		return
	}

	lastMsgIDs := []primitive.ObjectID{}
	for _, messageID := range seen.Messages {
		lastMsgID, err := primitive.ObjectIDFromHex(messageID)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		lastMsgIDs = append(lastMsgIDs, lastMsgID)
	}

	var name string
	if c.User.Profile != nil && c.User.Profile.Name != nil {
		name = *c.User.Profile.Name
	} else {
		name = c.User.Email
	}

	var avatar string
	if c.User.Profile != nil && c.User.Profile.Image != nil {
		avatar = *c.User.Profile.Image
	} else {
		avatar = ""
	}

	seenEvent := event.SeenEvent{
		TypeMessage: "seen",
		GroupID:     seen.GroupID,
		UserID:      userID,
		User: &event.UserSnapshot{
			UserID: c.User.ID.Hex(),
			Name:   name,
			Avatar: avatar,
		},
		Messages: seen.Messages,
		SeenAt:   time.Now(),
	}

	go func() {
		if err := c.MessageRepository.MarkSeenUpTo(
			context.Background(),
			seen.GroupID,
			userID,
			lastMsgIDs,
		); err != nil {
			log.Println("mark seen error:", err)
		}

		c.Hub.SeenBroadcasts <- &seenEvent
	}()

	resetEvent := &event.UnreadUpdateEvent{
		TypeMessage: "unread_count",
		GroupID:     seen.GroupID,
		Count:       0,
	}

	c.Send <- resetEvent.ToJSON()
}

func (c *Client) sendInitialUnreadCounts() {
	groupIDs := make([]string, 0, len(c.GroupIds))
	for groupID := range c.GroupIds {
		groupIDs = append(groupIDs, groupID)
	}

	counts, err := c.MessageRepository.GetUnreadCounts(context.Background(), groupIDs, c.User.ID.Hex())
	if err != nil {
		log.Println("get unread counts error:", err)
		return
	}

	e := &event.UnreadEvent{
		TypeMessage: "unread_count",
		Counts:      counts,
	}

	c.Send <- e.ToJSON()
}
