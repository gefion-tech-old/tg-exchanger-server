package mocksqlstore

import (
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type UserRepository struct {
	users map[int64]*models.User
}

func (r *UserRepository) Create(req *models.UserRequest) (*models.User, error) {
	u := &models.User{
		ChatID:    req.ChatID,
		Username:  req.Username,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	r.users[u.ChatID] = u
	return u, nil
}
