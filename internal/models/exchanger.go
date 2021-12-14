package models

type Exchanger struct {
	ID         uint64 `json:"id"`
	Name       string `json:"string" binding:"required"`
	UrlToParse string `json:"url" binding:"required"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
