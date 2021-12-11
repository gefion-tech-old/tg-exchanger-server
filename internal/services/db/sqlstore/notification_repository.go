package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type NotificationRepository struct {
	store *sql.DB
}

func (r *NotificationRepository) Create(n *models.Notification) (*models.Notification, error) {
	if err := r.store.QueryRow(
		`
		INSERT INTO notifications(type, chat_id, username, code, user_card)
		SELECT $1, $2, $3, $4, $5
		RETURNING id, type, status, chat_id, username, code, user_card, created_at, updated_at
		`,
		n.Type,
		n.MetaData.ChatID,
		n.MetaData.Username,
		n.MetaData.Code,
		n.MetaData.UserCard,
	).Scan(
		&n.ID,
		&n.Type,
		&n.Status,
		&n.MetaData.ChatID,
		&n.MetaData.Username,
		&n.MetaData.Code,
		&n.MetaData.UserCard,
		&n.CreatedAt,
		&n.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return n, nil
}

func (r *NotificationRepository) Delete(n *models.Notification) (*models.Notification, error) {
	if err := r.store.QueryRow(
		`
		DELETE FROM notifications
		WHERE id=$1
		RETURNING id, type, status, chat_id, username, code, user_card, created_at, updated_at
		`,
		n.ID,
	).Scan(
		&n.ID,
		&n.Type,
		&n.Status,
		&n.MetaData.ChatID,
		&n.MetaData.Username,
		&n.MetaData.Code,
		&n.MetaData.UserCard,
		&n.CreatedAt,
		&n.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return n, nil
}
