package hub

import (
	"context"
	"go-split/internal/domain/repository"
	"go-split/internal/interface/websocket/event"
	"log"
)

type Hub struct {
	Groups            map[string]map[*Client]bool
	MessageBroadcasts chan *event.MessageEvent
	SeenBroadcasts    chan *event.SeenEvent
	UnreadBroadcasts  chan *event.UnreadEvent
	Register          chan *Client
	Unregister        chan *Client
	MessageRepository repository.MessageRepository
}

func NewHub(messageRepository repository.MessageRepository) *Hub {
	return &Hub{
		Groups:            make(map[string]map[*Client]bool),
		MessageBroadcasts: make(chan *event.MessageEvent),
		SeenBroadcasts:    make(chan *event.SeenEvent),
		UnreadBroadcasts:  make(chan *event.UnreadEvent),
		Register:          make(chan *Client),
		Unregister:        make(chan *Client),
		MessageRepository: messageRepository,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			for groupID := range client.GroupIds {
				if _, ok := h.Groups[groupID]; !ok {
					h.Groups[groupID] = make(map[*Client]bool)
				}
				h.Groups[groupID][client] = true
				h.broadcastPresence(groupID, "join")
			}
			go client.sendInitialUnreadCounts()

		case client := <-h.Unregister:
			for groupID := range client.GroupIds {
				if _, ok := h.Groups[groupID]; ok {
					delete(h.Groups[groupID], client)
					h.broadcastPresence(groupID, "leave")
				}
				if len(h.Groups[groupID]) == 0 {
					delete(h.Groups, groupID)
				}
			}
			close(client.Send)

		case msg := <-h.MessageBroadcasts:
			h.broadcastMessage(msg)

		case seen := <-h.SeenBroadcasts:
			h.broadcastSeen(seen)

		case unread := <-h.UnreadBroadcasts:
			h.broadcastUnread(unread)
		}
	}
}

func (h *Hub) broadcastMessage(msg *event.MessageEvent) {
	clients, ok := h.Groups[msg.GroupID]
	if !ok {
		return
	}
	for client := range clients {
		select {
		case client.Send <- msg.ToJSON():
		default:
			h.removeClient(clients, client, msg.GroupID)
		}
	}
}

func (h *Hub) broadcastSeen(seenEvent *event.SeenEvent) {
	clients, ok := h.Groups[seenEvent.GroupID]
	if !ok {
		return
	}
	for client := range clients {
		if client.User.ID.Hex() == seenEvent.UserID {
			continue
		}
		select {
		case client.Send <- seenEvent.ToJSON():
		default:
			h.removeClient(clients, client, seenEvent.GroupID)
		}
	}
}

func (h *Hub) broadcastUnread(unreadEvent *event.UnreadEvent) {
	clients, ok := h.Groups[unreadEvent.GroupID]
	if !ok {
		return
	}
	for client := range clients {
		if client.User.ID.Hex() == unreadEvent.SenderID {
			continue
		}

		go func(c *Client) {
			counts, err := c.MessageRepository.GetUnreadCount(context.Background(), unreadEvent.GroupID, c.User.ID.Hex())
			if err != nil {
				log.Println("get unread count error:", err)
				return
			}

			payload := &event.UnreadUpdateEvent{
				TypeMessage: "unread_count",
				GroupID:     unreadEvent.GroupID,
				Count:       counts,
			}

			select {
			case c.Send <- payload.ToJSON():
			default:
				h.removeClient(clients, c, unreadEvent.GroupID)
			}
		}(client)
	}
}

func (h *Hub) broadcastPresence(groupID string, presenceType string) {
	clients, ok := h.Groups[groupID]
	if !ok {
		return
	}

	users := h.uniqueUsers(clients)

	presenceEvent := &event.PresenceEvent{
		TypeMessage: "presence",
		GroupID:     groupID,
		Count:       len(clients),
		Type:        presenceType,
		Users:       users,
	}

	for client := range clients {
		select {
		case client.Send <- presenceEvent.ToJSON():
		default:
			h.removeClient(clients, client, groupID)
		}
	}
}

func (h *Hub) uniqueUsers(clients map[*Client]bool) []*event.UserSnapshot {
	seen := make(map[string]bool)
	users := []*event.UserSnapshot{}

	for client := range clients {
		id := client.User.ID.Hex()
		if seen[id] {
			continue
		}
		seen[id] = true

		u := client.User
		var name string
		if u.Profile != nil && u.Profile.Name != nil {
			name = *u.Profile.Name
		} else {
			name = u.Email
		}
		var avatar string
		if u.Profile != nil && u.Profile.Image != nil {
			avatar = *u.Profile.Image
		}

		users = append(users, &event.UserSnapshot{
			UserID: id,
			Name:   name,
			Avatar: avatar,
		})
	}

	return users
}

func (h *Hub) removeClient(clients map[*Client]bool, client *Client, groupID string) {
	delete(clients, client)
	if len(clients) == 0 {
		delete(h.Groups, groupID)
	}
}
