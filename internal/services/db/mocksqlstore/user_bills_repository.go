package mocksqlstore

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type UserBillsRepository struct {
	bills map[uint]*models.Bill
}

func (r *UserBillsRepository) Create(b *models.Bill) (*models.Bill, error) {
	return nil, nil
}

func (r *UserBillsRepository) Delete(b *models.Bill) (*models.Bill, error) {
	return nil, nil
}

func (r *UserBillsRepository) All(chatID int64) ([]*models.Bill, error) {
	return nil, nil
}
