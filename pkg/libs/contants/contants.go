package contants

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyToken  contextKey = "token"
	ContextKeyRole   contextKey = "role"
)

const (
	StatusActive   = "active"
	Suspended      = "suspended"
	StatusInactive = "inactive"
)
