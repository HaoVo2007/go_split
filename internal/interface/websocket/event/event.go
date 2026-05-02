package event

import (
	"encoding/json"
	"time"
)

// ==================== CLIENT → SERVER EVENTS ====================
type IncomingMessageEvent struct {
	TypeMessage string `json:"type_message"`
	GroupID     string `json:"group_id"`
	Message     string `json:"message"`
}
type SeenEventSendToServer struct {
	TypeMessage string   `json:"type_message"`
	GroupID     string   `json:"group_id"`
	Messages    []string `json:"messages"`
}

// ==================== SERVER → CLIENT EVENTS ====================
type OutGoingMessageEvent struct {
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
type SeenEventSendToClient struct {
	TypeMessage string    `json:"type_message"` 
	GroupID     string    `json:"group_id"`
	UserID      string    `json:"user_id"`
	User        *UserInfo `json:"user"`
	Messages    []string  `json:"messages"`
	SeenAt      time.Time `json:"seen_at"`
}

type UserInfo struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ==================== TOJSON METHODS ====================

func (m *IncomingMessageEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *OutGoingMessageEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *PresenceEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *SeenEventSendToClient) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *SeenEventSendToServer) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}
