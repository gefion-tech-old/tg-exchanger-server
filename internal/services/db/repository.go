package db

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type AdminPanelRepository interface {
	Logs() LoggerRepository
	Bills() UserBillsRepository
	BotMessages() BotMessagesRepository
	Notification() NotificationRepository
	Exchanger() ExchangerRepository
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
	DeleteSelection(date_from, date_to string) ([]*models.LogRecord, error)
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
	CheckNew() (int, error)
	Count(querys interface{}) (int, error)
	Selection(querys interface{}) ([]*models.Notification, error)
}

type BotMessagesRepository interface {
	Create(m *models.BotMessage) error
	Update(m *models.BotMessage) error
	Delete(m *models.BotMessage) error
	Get(m *models.BotMessage) error
	Count(querys interface{}) (int, error)
	Selection(querys interface{}) ([]*models.BotMessage, error)
}
