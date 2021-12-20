package db

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

/*
	Интерфейс репозитория для работы с таблицей users
*/
type UserRepository interface {
	Bills() UserBillsRepository
	Create(u *models.User) (*models.User, error)
	RegisterAsManager(u *models.User) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	GetAllManagers() ([]*models.User, error)
}

type UserBillsRepository interface {
	Create(b *models.Bill) (*models.Bill, error)
	Delete(b *models.Bill) (*models.Bill, error)
	All(chatID int64) ([]*models.Bill, error)
	FindById(b *models.Bill) (*models.Bill, error)
}

type ManagerRepository interface {
	BotMessages() BotMessagesRepository
	Notification() NotificationRepository
	Exchanger() ExchangerRepository
}

type ExchangerRepository interface {
	Create(e *models.Exchanger) (*models.Exchanger, error)
	Update(e *models.Exchanger) (*models.Exchanger, error)
	Count() (int, error)
	GetByName(e *models.Exchanger) (*models.Exchanger, error)
	Delete(e *models.Exchanger) (*models.Exchanger, error)
	GetSlice(limit int) ([]*models.Exchanger, error)
}

type NotificationRepository interface {
	Create(n *models.Notification) (*models.Notification, error)
	Delete(n *models.Notification) (*models.Notification, error)
	Get(n *models.Notification) (*models.Notification, error)
	GetSlice(limit int) ([]*models.Notification, error)
	UpdateStatus(n *models.Notification) (*models.Notification, error)
	Count() (int, error)
}

type BotMessagesRepository interface {
	Create(m *models.BotMessage) (*models.BotMessage, error)
	Get(m *models.BotMessage) (*models.BotMessage, error)
	GetSlice(limit int) ([]*models.BotMessage, error)
	Update(m *models.BotMessage) (*models.BotMessage, error)
	Delete(m *models.BotMessage) (*models.BotMessage, error)
	Count() (int, error)
}
