package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type UserRepository struct {
	users map[int64]*models.User

	userBillsRepository *UserBillsRepository
}

/*
	==========================================================================================
	КОНСТРУКТОРЫ ВЛОЖЕННЫХ СТРУКТУР
	==========================================================================================
*/

func (r *UserRepository) Bills() db.UserBillsRepository {
	if r.userBillsRepository != nil {
		return r.userBillsRepository
	}

	r.userBillsRepository = &UserBillsRepository{
		bills: make(map[uint]*models.Bill),
	}

	return r.userBillsRepository
}

/*
	==========================================================================================
	КОНЕЧНЫЕ МЕТОДЫ ТЕКУЩЕЙ СТРУКТУРЫ
	==========================================================================================
*/

func (r *UserRepository) Create(u *models.User) (*models.User, error) {
	u.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	u.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.users[u.ChatID] = u
	return r.users[u.ChatID], nil
}

func (r *UserRepository) RegisterAsManager(u *models.User) (*models.User, error) {
	for _, user := range r.users {
		if u.Username == user.Username {
			r.users[user.ChatID].Hash = u.Hash
			r.users[user.ChatID].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
			return r.users[user.ChatID], nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, sql.ErrNoRows
}
