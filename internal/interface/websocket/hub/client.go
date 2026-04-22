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
	UserID            string
	MessageRepository repository.MessageRepository
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

		var msg event.MessageEvent
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("error: %v", err)
			break
		}

		if !c.GroupIds[msg.GroupID] {
			log.Println("User send message to group not in list of groups")
			continue
		}

		message := entity.Messages{
			ID:        primitive.NewObjectID(),
			GroupID:   msg.GroupID,
			Message:   msg.Message,
			UserID:    c.UserID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = c.MessageRepository.CreateMessage(context.Background(), message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		c.Hub.Broadcasts <- &msg
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
