package mocksqlstore

import (
	"database/sql"
	"time"

	AppMath "github.com/gefion-tech/tg-exchanger-server/internal/core/math"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type BotMessagesRepository struct {
	messages map[int]*models.BotMessage
}

func (r *BotMessagesRepository) Create(m *models.BotMessage) error {
	m.ID = len(r.messages) + 1
	m.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	m.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.messages[m.ID] = m
	return nil
}

func (r *BotMessagesRepository) Get(m *models.BotMessage) error {
	for i, msg := range r.messages {
		if msg.Connector == m.Connector {
			r.rewrite(i, m)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *BotMessagesRepository) Selection(querys interface{}) ([]*models.BotMessage, error) {
	arr := []*models.BotMessage{}
	q := querys.(*models.BotMessageSelection)

	for i, v := range r.messages {
		if i > AppMath.OffsetThreshold(q.Page, q.Limit) && i <= AppMath.OffsetThreshold(q.Page, q.Limit)+q.Limit {
			arr = append(arr, v)
		}
		i++
	}

	return arr, nil
}

func (r *BotMessagesRepository) Update(m *models.BotMessage) error {
	for _, msg := range r.messages {
		if msg.ID == m.ID {
			r.messages[msg.ID].MessageText = m.MessageText
			r.messages[msg.ID].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

			r.rewrite(m.ID, m)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *BotMessagesRepository) Delete(m *models.BotMessage) error {
	for _, msg := range r.messages {
		if msg.ID == m.ID {
			r.rewrite(m.ID, m)
			defer delete(r.messages, msg.ID)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *BotMessagesRepository) Count(querys interface{}) (int, error) {
	return len(r.messages), nil
}

func (r *BotMessagesRepository) rewrite(id int, to *models.BotMessage) {
	to.ID = r.messages[id].ID
	to.Connector = r.messages[id].Connector
	to.MessageText = r.messages[id].MessageText
	to.CreatedBy = r.messages[id].CreatedBy
	to.CreatedAt = r.messages[id].CreatedAt
	to.UpdatedAt = r.messages[id].UpdatedAt
}
