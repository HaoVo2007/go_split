package event

import "encoding/json"

type MessageEvent struct {
	GroupID string `json:"group_id"`
	Message string `json:"message"`
}

type PresenceEvent struct {
	GroupID string `json:"group_id"`
	Count   int    `json:"count"`
	Type    string `json:"type"`
}

func (m *MessageEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *PresenceEvent) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}
