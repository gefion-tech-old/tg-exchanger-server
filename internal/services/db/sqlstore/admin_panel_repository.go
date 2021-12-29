package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type AdminPanelRepository struct {
	store                  *sql.DB
	botMessagesRepository  *BotMessagesRepository
	notificationRepository *NotificationRepository
	exchangerRepository    *ExchangerRepository
	userBillsRepository    *UserBillsRepository
	logsRepository         *LoggerRepository
}

/*
	==========================================================================================
	КОНСТРУКТОРЫ ВЛОЖЕННЫХ СТРУКТУР
	==========================================================================================
*/

func (r *AdminPanelRepository) Logs() db.LoggerRepository {
	if r.logsRepository != nil {
		return r.logsRepository
	}

	r.logsRepository = &LoggerRepository{
		store: r.store,
	}

	return r.logsRepository
}

func (r *AdminPanelRepository) Bills() db.UserBillsRepository {
	if r.userBillsRepository != nil {
		return r.userBillsRepository
	}

	r.userBillsRepository = &UserBillsRepository{
		store: r.store,
	}

	return r.userBillsRepository
}

func (r *AdminPanelRepository) Exchanger() db.ExchangerRepository {
	if r.exchangerRepository != nil {
		return r.exchangerRepository
	}

	r.exchangerRepository = &ExchangerRepository{
		store: r.store,
	}

	return r.exchangerRepository
}

func (r *AdminPanelRepository) Notification() db.NotificationRepository {
	if r.notificationRepository != nil {
		return r.notificationRepository
	}

	r.notificationRepository = &NotificationRepository{
		store: r.store,
	}

	return r.notificationRepository
}

func (r *AdminPanelRepository) BotMessages() db.BotMessagesRepository {
	if r.botMessagesRepository != nil {
		return r.botMessagesRepository
	}

	r.botMessagesRepository = &BotMessagesRepository{
		store: r.store,
	}

	return r.botMessagesRepository
}
