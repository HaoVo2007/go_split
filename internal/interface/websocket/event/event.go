package event

import (
	"encoding/json"
	"time"
)

// ==================== CLIENT → SERVER ====================

type MessageRequest struct {
	TypeMessage string `json:"type_message"`
	GroupID     string `json:"group_id"`
	Message     string `json:"message"`
}

type SeenRequest struct {
	TypeMessage string   `json:"type_message"`
	GroupID     string   `json:"group_id"`
	Messages    []string `json:"messages"`
}

// ==================== SERVER → CLIENT ====================

type MessageEvent struct {
	TypeMessage string    `json:"type_message"`
	GroupID     string    `json:"group_id"`
	Message     string    `json:"message"`
	SenderID    string    `json:"sender_id"`
	SenderName  string    `json:"sender_name"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
}

type SeenEvent struct {
	TypeMessage string        `json:"type_message"`
	GroupID     string        `json:"group_id"`
	UserID      string        `json:"user_id"`
	User        *UserSnapshot `json:"user"`
	Messages    []string      `json:"messages"`
	SeenAt      time.Time     `json:"seen_at"`
}

type UnreadEvent struct {
	TypeMessage string         `json:"type_message"`
	GroupID     string         `json:"group_id"`
	SenderID    string         `json:"-"`
	Counts      map[string]int `json:"counts"`
}

type UnreadUpdateEvent struct {
	TypeMessage string `json:"type_message"`
	GroupID     string `json:"group_id"`
	Count       int    `json:"count"`
}

type PresenceEvent struct {
	TypeMessage string          `json:"type_message"`
	GroupID     string          `json:"group_id"`
	Count       int             `json:"count"`
	Type        string          `json:"type"`
	Users       []*UserSnapshot `json:"users"`
}

// ==================== SHARED TYPES ====================

type UserSnapshot struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ==================== TOJSON METHODS ====================

func (m *MessageRequest) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *MessageEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *SeenEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *SeenRequest) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *UnreadEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *PresenceEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *UnreadUpdateEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}
