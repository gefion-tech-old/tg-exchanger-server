package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
)

type BotMessagesRepository struct {
	messages map[int]*models.BotMessage
}

func (r *BotMessagesRepository) Create(m *models.BotMessage) (*models.BotMessage, error) {
	m.ID = len(r.messages) + 1
	m.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	m.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.messages[m.ID] = m
	return r.messages[m.ID], nil
}

func (r *BotMessagesRepository) Get(m *models.BotMessage) (*models.BotMessage, error) {
	for _, msg := range r.messages {
		if msg.Connector == m.Connector {
			return r.messages[msg.ID], nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *BotMessagesRepository) Selection(page, limit int) ([]*models.BotMessage, error) {
	arr := []*models.BotMessage{}

	for i, v := range r.messages {
		if i > tools.OffsetThreshold(page, limit) && i <= tools.OffsetThreshold(page, limit)+limit {
			arr = append(arr, v)
		}
		i++
	}

	return arr, nil
}

func (r *BotMessagesRepository) Update(m *models.BotMessage) (*models.BotMessage, error) {
	for _, msg := range r.messages {
		if msg.ID == m.ID {
			r.messages[msg.ID].MessageText = m.MessageText
			r.messages[msg.ID].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
			return r.messages[msg.ID], nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *BotMessagesRepository) Delete(m *models.BotMessage) (*models.BotMessage, error) {
	for _, msg := range r.messages {
		if msg.ID == m.ID {
			defer delete(r.messages, msg.ID)
			return r.messages[m.ID], nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *BotMessagesRepository) Count() (int, error) {
	return len(r.messages), nil
}
