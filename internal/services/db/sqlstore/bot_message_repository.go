package sqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type BotMessagesRepository struct {
	store *sql.DB
}

func (r *BotMessagesRepository) Create(m *models.BotMessage) (*models.BotMessage, error) {
	if err := r.store.QueryRow(
		`
		INSERT INTO bot_messages (connector, message_text, created_by)
		SELECT $1, $2, $3
		WHERE NOT EXISTS (SELECT connector FROM bot_messages WHERE connector=$4)
		RETURNING id, connector, message_text, created_by, created_at, updated_at
		`,
		m.Connector,
		m.MessageText,
		m.CreatedBy,
		m.Connector,
	).Scan(
		&m.ID,
		&m.Connector,
		&m.MessageText,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return m, nil
}

/*
	Получить конкретное сообщение из таблицы `bot_messages`
*/
func (r *BotMessagesRepository) Get(m *models.BotMessage) (*models.BotMessage, error) {
	if err := r.store.QueryRow(
		`
		SELECT id, connector, message_text, created_by, created_at, updated_at
		FROM bot_messages WHERE connector=$1
		`,
		m.Connector,
	).Scan(
		&m.ID,
		&m.Connector,
		&m.MessageText,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return m, nil
}

/*
	Получить все сообщения из таблицы `bot_messages`
*/
func (r *BotMessagesRepository) GetAll() ([]*models.BotMessage, error) {
	bm := []*models.BotMessage{}

	rows, err := r.store.Query(
		`
		SELECT id, connector, message_text, created_by, created_at, updated_at
		FROM bot_messages
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := &models.BotMessage{}
		if err := rows.Scan(
			&m.ID,
			&m.Connector,
			&m.MessageText,
			&m.CreatedBy,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			continue
		}

		bm = append(bm, m)
	}

	return bm, nil
}

/*
	Обновить конкретное сообщение в таблице `bot_messages`
*/
func (r *BotMessagesRepository) Update(m *models.BotMessage) (*models.BotMessage, error) {
	if err := r.store.QueryRow(
		`
		UPDATE bot_messages
		SET message_text=$1, updated_at=$2
		WHERE connector=$3
		RETURNING id, connector, message_text, created_by, created_at, updated_at
		`,
		m.MessageText,
		time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
		m.Connector,
	).Scan(
		&m.ID,
		&m.Connector,
		&m.MessageText,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return m, nil
}

/*
	Удалить конкретное сообщение в таблице `bot_messages`
*/
func (r *BotMessagesRepository) Delete(m *models.BotMessage) (*models.BotMessage, error) {
	if err := r.store.QueryRow(
		`
		DELETE FROM bot_messages
		WHERE connector=$1
		RETURNING id, connector, message_text, created_by, created_at, updated_at
		`,
		m.Connector,
	).Scan(
		&m.ID,
		&m.Connector,
		&m.MessageText,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return m, nil
}
