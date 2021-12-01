package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type UserRepository struct {
	store *sql.DB
}

/*
	==========================================================================================
	КОНЕЧНЫЕ МЕТОДЫ ТЕКУЩЕЙ СТРУКТУРЫ
	==========================================================================================
*/

func (r *UserRepository) Create(req *models.UserRequest) (*models.User, error) {
	// Создаю объект пользователя который будет записан в БД
	u := &models.User{
		ChatID:   req.ChatID,
		Username: req.Username,
	}

	if err := r.store.QueryRow(
		`
		INSERT INTO users (chat_id, username) 
		SELECT $1, $2
		WHERE NOT EXISTS (SELECT chat_id FROM users WHERE chat_id=$3)
		RETURNING chat_id, username, hash, created_at, updated_at
		`,
		u.ChatID,
		u.Username,
		u.ChatID,
	).Scan(
		&u.ChatID,
		&u.Username,
		&u.Hash,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}
