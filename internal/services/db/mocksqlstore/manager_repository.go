package mocksqlstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type ManagerRepository struct {
	botMessagesRepository  *BotMessagesRepository
	notificationRepository *NotificationRepository
}

func (r *ManagerRepository) Notification() db.NotificationRepository {
	if r.notificationRepository != nil {
		return r.notificationRepository
	}

	r.notificationRepository = &NotificationRepository{
		notification: make(map[uint]*models.Notification),
	}

	return r.notificationRepository
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
