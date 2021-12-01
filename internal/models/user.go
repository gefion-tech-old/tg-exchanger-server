package models

// Структура записи в таблице `users`
type User struct {
	ChatID    int64  `json:"chat_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type UserRequest struct {
	ChatID   int64  `json:"chat_id"`
	Username string `json:"username"`
}
