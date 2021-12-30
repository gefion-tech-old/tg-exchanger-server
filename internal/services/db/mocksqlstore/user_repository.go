package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type UserRepository struct {
	users map[int64]*models.User
}

/*
	==========================================================================================
	КОНЕЧНЫЕ МЕТОДЫ ТЕКУЩЕЙ СТРУКТУРЫ
	==========================================================================================
*/

func (r *UserRepository) Create(u *models.User) error {
	u.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	u.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.users[u.ChatID] = u
	return nil
}

func (r *UserRepository) RegisterInAdminPanel(u *models.User) error {
	for _, user := range r.users {
		if u.Username == user.Username {
			r.users[user.ChatID].Hash = u.Hash
			r.users[user.ChatID].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

			u.CreatedAt = r.users[user.ChatID].CreatedAt
			u.UpdatedAt = r.users[user.ChatID].UpdatedAt
			u.ChatID = r.users[user.ChatID].ChatID
			u.Username = r.users[user.ChatID].Username
			return nil
		}
	}
	return sql.ErrNoRows
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *UserRepository) GetAllManagers() ([]*models.User, error) {
	return nil, nil
}
