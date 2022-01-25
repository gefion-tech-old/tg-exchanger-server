package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type AdminPanelRepository struct {
	store                        *sql.DB
	botMessagesRepository        *BotMessagesRepository
	notificationRepository       *NotificationRepository
	exchangerRepository          *ExchangerRepository
	userBillsRepository          *UserBillsRepository
	logsRepository               *LoggerRepository
	merchantAutopayoutRepository *MerchantAutopayoutRepository
	exchangeRequestRepository    *ExchangeRequestRepository
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

func (r *AdminPanelRepository) MerchantAutopayout() db.MerchantAutopayoutRepository {
	if r.merchantAutopayoutRepository != nil {
		return r.merchantAutopayoutRepository
	}

	r.merchantAutopayoutRepository = &MerchantAutopayoutRepository{
		store: r.store,
	}

	return r.merchantAutopayoutRepository
}

func (r *AdminPanelRepository) ExchangeRequest() db.ExchangeRequestRepository {
	if r.exchangeRequestRepository != nil {
		return r.exchangeRequestRepository
	}

	r.exchangeRequestRepository = &ExchangeRequestRepository{
		store: r.store,
	}

	return r.exchangeRequestRepository
}
