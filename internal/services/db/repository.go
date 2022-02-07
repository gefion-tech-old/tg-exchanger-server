package db

import (
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type AdminPanelRepository interface {
	Logs() LoggerRepository
	Bills() UserBillsRepository
	BotMessages() BotMessagesRepository
	Notification() NotificationRepository
	Exchanger() ExchangerRepository
	MerchantAutopayout() MerchantAutopayoutRepository
	ExchangeRequest() ExchangeRequestRepository
}

type UserRepository interface {
	Create(u *models.User) error
	RegisterInAdminPanel(u *models.User) error
	FindByUsername(username string) (*models.User, error)
	GetAllManagers() ([]*models.User, error)
}

type UserBillsRepository interface {
	Create(b *models.Bill) error
	Delete(b *models.Bill) error
	All(chatID int64) ([]*models.Bill, error)
	FindById(b *models.Bill) error
}

type LoggerRepository interface {
	Create(lr *models.LogRecord) error
	Delete(lr *models.LogRecord) error
	Count(querys interface{}) (int, error)
	Selection(querys interface{}) ([]*models.LogRecord, error)
	DeleteSelection(querys interface{}) ([]*models.LogRecord, error)
}

type ExchangerRepository interface {
	Create(e *models.Exchanger) error
	Delete(e *models.Exchanger) error
	Update(e *models.Exchanger) error
	GetByName(e *models.Exchanger) error
	Count(querys interface{}) (int, error)
	Selection(querys interface{}) ([]*models.Exchanger, error)
}

type NotificationRepository interface {
	Create(n *models.Notification) error
	Delete(n *models.Notification) error
	Update(n *models.Notification) error
	Get(n *models.Notification) error
	Count(querys interface{}) (int, error)
	Selection(querys interface{}) ([]*models.Notification, error)
	CheckNew() (int, error)
}

type BotMessagesRepository interface {
	Create(m *models.BotMessage) error
	Update(m *models.BotMessage) error
	Delete(m *models.BotMessage) error
	Get(m *models.BotMessage) error
	Count(querys interface{}) (int, error)
	Selection(querys interface{}) ([]*models.BotMessage, error)
}

type MerchantAutopayoutRepository interface {
	Create(m *models.MerchantAutopayout) error
	Update(m *models.MerchantAutopayout) error
	Delete(m *models.MerchantAutopayout) error
	Get(m *models.MerchantAutopayout) error
	Count(querys interface{}) (int, error)
	Selection(querys interface{}) ([]*models.MerchantAutopayout, error)
	GetFistIfActive(service string) (*models.MerchantAutopayout, error)
	GetAllMerchants() ([]*models.MerchantAutopayout, error)
}

type ExchangeRequestRepository interface {
	Create(r *models.ExchangeRequest) error
	Update(r *models.ExchangeRequest) error
	Get(r *models.ExchangeRequest) error
	Delete(r *models.ExchangeRequest) error
	Count(querys interface{}) (int, error)
	Selection(querys interface{}) ([]*models.ExchangeRequest, error)
	GetAllByStatus(s ...AppType.ExchangeRequestStatus) ([]*models.ExchangeRequest, error)
}
