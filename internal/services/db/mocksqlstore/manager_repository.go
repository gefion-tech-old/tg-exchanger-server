package mocksqlstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type ManagerRepository struct {
	botMessagesRepository *BotMessagesRepository
}

func (r *ManagerRepository) BotMessages() db.BotMessagesRepository {
	if r.botMessagesRepository != nil {
		return r.botMessagesRepository
	}

	r.botMessagesRepository = &BotMessagesRepository{
		messages: make(map[uint]*models.BotMessage),
	}

	return r.botMessagesRepository
}
