package models

type UserMessages struct {
	ID        uint64 `json:"id"`
	ChatID    int64  `json:"chat_id"`
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
