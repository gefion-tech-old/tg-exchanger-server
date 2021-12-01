package models

// Структура записи в таблице `users`
type User struct {
	ChatID    int64  `json:"chat_id"`
	Username  string `json:"username"`
	Hash      string `json:"hash"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserRequest struct {
	ChatID   int64  `json:"chat_id" binding:"required"`
	Username string `json:"username" binding:"required"`
}
