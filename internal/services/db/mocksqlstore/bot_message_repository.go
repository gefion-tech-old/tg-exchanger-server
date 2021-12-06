package mocksqlstore

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type BotMessagesRepository struct {
	messages map[uint]*models.BotMessage
}

func (r *BotMessagesRepository) Create(m *models.BotMessage) (*models.BotMessage, error) {
	return nil, nil
}

func (r *BotMessagesRepository) Get(m *models.BotMessage) (*models.BotMessage, error) {
	return nil, nil
}
func (r *BotMessagesRepository) GetAll() ([]*models.BotMessage, error) {
	return nil, nil
}
func (r *BotMessagesRepository) Update(m *models.BotMessage) (*models.BotMessage, error) {
	return nil, nil
}
func (r *BotMessagesRepository) Delete(m *models.BotMessage) (*models.BotMessage, error) {
	return nil, nil
}
