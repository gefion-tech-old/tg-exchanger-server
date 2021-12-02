package sqlstore

import (
	"database/sql"
	"time"

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

func (r *UserRepository) Create(req *models.UserFromBotRequest) (*models.User, error) {
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

func (r *UserRepository) RegisterAsManager(u *models.User) (*models.User, error) {
	if err := r.store.QueryRow(
		`
		UPDATE users
		SET hash=$1, updated_at=$2
		WHERE username=$3
		RETURNING chat_id, username, hash, created_at, updated_at
		`,
		u.Hash,
		time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
		u.Username,
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

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	u := &models.User{}

	if err := r.store.QueryRow(
		`
		SELECT chat_id, username, created_at, updated_at FROM users WHERE username=$1
		`,
		username,
	).Scan(
		&u.ChatID,
		&u.Username,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}
