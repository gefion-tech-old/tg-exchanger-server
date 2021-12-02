package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type UserRepository struct {
	users map[int64]*models.User
}

func (r *UserRepository) Create(req *models.UserFromBotRequest) (*models.User, error) {
	u := &models.User{
		ChatID:    req.ChatID,
		Username:  req.Username,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	r.users[u.ChatID] = u
	return r.users[u.ChatID], nil
}

func (r *UserRepository) RegisterAsManager(req *models.User) (*models.User, error) {
	u := &models.User{
		ChatID:    1,
		Username:  req.Username,
		Hash:      req.Hash,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	r.users[u.ChatID] = u
	return u, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}

	return nil, sql.ErrNoRows
}
