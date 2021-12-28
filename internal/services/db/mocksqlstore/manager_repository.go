package mocksqlstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type ManagerRepository struct {
	botMessagesRepository  *BotMessagesRepository
	notificationRepository *NotificationRepository
	exchangerRepository    *ExchangerRepository
	userBillsRepository    *UserBillsRepository
}

func (r *ManagerRepository) Bills() db.UserBillsRepository {
	if r.userBillsRepository != nil {
		return r.userBillsRepository
	}

	r.userBillsRepository = &UserBillsRepository{
		bills: make(map[int]*models.Bill),
	}

	return r.userBillsRepository
}

func (r *ManagerRepository) Exchanger() db.ExchangerRepository {
	if r.exchangerRepository != nil {
		return r.exchangerRepository
	}

	r.exchangerRepository = &ExchangerRepository{
		exchangers: make(map[int]*models.Exchanger),
	}

	return r.exchangerRepository
}

func (r *ManagerRepository) Notification() db.NotificationRepository {
	if r.notificationRepository != nil {
		return r.notificationRepository
	}

	r.notificationRepository = &NotificationRepository{
		notification: make(map[int]*models.Notification),
	}

	return r.notificationRepository
}

func (r *ManagerRepository) BotMessages() db.BotMessagesRepository {
	if r.botMessagesRepository != nil {
		return r.botMessagesRepository
	}

	r.botMessagesRepository = &BotMessagesRepository{
		messages: make(map[int]*models.BotMessage),
	}

	return r.botMessagesRepository
}
