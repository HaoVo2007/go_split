package event

import (
	"encoding/json"
	"time"
)

type MessageEvent struct {
	TypeMessage string    `json:"type_message"`
	GroupID     string    `json:"group_id"`
	Message     string    `json:"message"`
	SenderID    string    `json:"sender_id"`
	SenderName  string    `json:"sender_name"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
}

type PresenceEvent struct {
	TypeMessage string          `json:"type_message"`
	GroupID     string          `json:"group_id"`
	Count       int             `json:"count"`
	Type        string          `json:"type"`
	Users       []*PresenceUser `json:"users"`
}

type PresenceUser struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (m *MessageEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *PresenceEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}
