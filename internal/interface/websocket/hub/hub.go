package hub

import (
	"go-split/internal/domain/entity"
	"go-split/internal/interface/websocket/event"
)

type Hub struct {
	Groups         map[string]map[*Client]bool
	Broadcasts     chan *event.OutGoingMessageEvent
	SeenBroadcasts chan *event.SeenEventSendToClient
	Register       chan *Client
	Unregister     chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Groups:         make(map[string]map[*Client]bool),
		Broadcasts:     make(chan *event.OutGoingMessageEvent),
		SeenBroadcasts: make(chan *event.SeenEventSendToClient),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			for groupId := range client.GroupIds {
				if _, ok := h.Groups[groupId]; !ok {
					h.Groups[groupId] = make(map[*Client]bool)
				}
				h.Groups[groupId][client] = true
				h.broadcastPresence(groupId, "join")
			}
		case client := <-h.Unregister:
			for groupId := range client.GroupIds {
				if _, ok := h.Groups[groupId]; ok {
					delete(h.Groups[groupId], client)
					h.broadcastPresence(groupId, "leave")
				}
				if len(h.Groups[groupId]) == 0 {
					delete(h.Groups, groupId)
				}
			}
			close(client.Send)
		case message := <-h.Broadcasts:
			if clients, ok := h.Groups[message.GroupID]; ok {
				for client := range clients {
					select {
					case client.Send <- message.ToJSON():
					default:
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.Groups, message.GroupID)
						}
					}
				}
			}
		case seen := <-h.SeenBroadcasts:
			h.handleSeen(seen)
		}
	}
}

func (h *Hub) broadcastPresence(groupID string, typeStr string) {
	clients, ok := h.Groups[groupID]
	if !ok {
		return
	}

	count := len(clients)

	userUnique := []*entity.Users{}
	userUniqueMap := make(map[string]bool)
	for client := range clients {
		if !userUniqueMap[client.User.ID.Hex()] {
			userUnique = append(userUnique, client.User)
			userUniqueMap[client.User.ID.Hex()] = true
		}
	}

	users := []*event.PresenceUser{}
	for _, user := range userUnique {
		var name string
		if user.Profile != nil && user.Profile.Name != nil {
			name = *user.Profile.Name
		} else {
			name = user.Email
		}
		var avatar string
		if user.Profile != nil && user.Profile.Image != nil {
			avatar = *user.Profile.Image
		} else {
			avatar = ""
		}
		users = append(users, &event.PresenceUser{
			UserID: user.ID.Hex(),
			Name:   name,
			Avatar: avatar,
		})
	}

	presenceEvent := &event.PresenceEvent{
		TypeMessage: "presence",
		GroupID:     groupID,
		Count:       count,
		Type:        typeStr,
		Users:       users,
	}

	for client := range clients {
		select {
		case client.Send <- presenceEvent.ToJSON():
		default:
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.Groups, groupID)
			}
		}
	}
}

func (h *Hub) handleSeen(seenEvent *event.SeenEventSendToClient) {
	clients, ok := h.Groups[seenEvent.GroupID]
	if !ok {
		return
	}
	for client := range clients {
		if client.User.ID.Hex() != seenEvent.UserID {
			select {
			case client.Send <- seenEvent.ToJSON():
			default:
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.Groups, seenEvent.GroupID)
				}
			}
		}
	}
}
