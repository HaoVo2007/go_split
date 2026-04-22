package hub

import (
	"go-split/internal/interface/websocket/event"
)

type Hub struct {
	Groups     map[string]map[*Client]bool
	Broadcasts chan *event.MessageEvent
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Groups:     make(map[string]map[*Client]bool),
		Broadcasts: make(chan *event.MessageEvent),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
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
		}	
	}
}

func (h *Hub) broadcastPresence(groupID string, typeStr string) {
	clients, ok := h.Groups[groupID]
	if !ok {
		return
	}

	count := len(clients)

	presenceEvent := &event.PresenceEvent{
		GroupID: groupID,
		Count:   count,
		Type:    typeStr,
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
