package mocksqlstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type AdminPanelRepository struct {
	botMessagesRepository        *BotMessagesRepository
	notificationRepository       *NotificationRepository
	exchangerRepository          *ExchangerRepository
	userBillsRepository          *UserBillsRepository
	logsRepository               *LoggerRepository
	merchantAutopayoutRepository *MerchantAutopayoutRepository
	exchangeRequestRepository    *ExchangeRequestRepository
}

func (r *AdminPanelRepository) Logs() db.LoggerRepository {
	if r.logsRepository != nil {
		return r.logsRepository
	}

	r.logsRepository = &LoggerRepository{
		logs: make(map[int]*models.LogRecord),
	}

	return r.logsRepository
}

func (r *AdminPanelRepository) Bills() db.UserBillsRepository {
	if r.userBillsRepository != nil {
		return r.userBillsRepository
	}

	r.userBillsRepository = &UserBillsRepository{
		bills: make(map[int]*models.Bill),
	}

	return r.userBillsRepository
}

func (r *AdminPanelRepository) Exchanger() db.ExchangerRepository {
	if r.exchangerRepository != nil {
		return r.exchangerRepository
	}

	r.exchangerRepository = &ExchangerRepository{
		exchangers: make(map[int]*models.Exchanger),
	}

	return r.exchangerRepository
}

func (r *AdminPanelRepository) Notification() db.NotificationRepository {
	if r.notificationRepository != nil {
		return r.notificationRepository
	}

	r.notificationRepository = &NotificationRepository{
		notification: make(map[int]*models.Notification),
	}

	return r.notificationRepository
}

func (r *AdminPanelRepository) BotMessages() db.BotMessagesRepository {
	if r.botMessagesRepository != nil {
		return r.botMessagesRepository
	}

	r.botMessagesRepository = &BotMessagesRepository{
		messages: make(map[int]*models.BotMessage),
	}

	return r.botMessagesRepository
}

func (r *AdminPanelRepository) MerchantAutopayout() db.MerchantAutopayoutRepository {
	if r.merchantAutopayoutRepository != nil {
		return r.merchantAutopayoutRepository
	}

	r.merchantAutopayoutRepository = &MerchantAutopayoutRepository{
		ma: make(map[int]*models.MerchantAutopayout),
	}

	return r.merchantAutopayoutRepository
}

func (r *AdminPanelRepository) ExchangeRequest() db.ExchangeRequestRepository {
	if r.exchangeRequestRepository != nil {
		return r.exchangeRequestRepository
	}

	r.exchangeRequestRepository = &ExchangeRequestRepository{
		er: make(map[int]*models.ExchangeRequest),
	}

	return r.exchangeRequestRepository
}
