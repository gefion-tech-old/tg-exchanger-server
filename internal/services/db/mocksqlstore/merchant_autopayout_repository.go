package mocksqlstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type MerchantAutopayoutRepository struct {
	ma map[int]*models.MerchantAutopayout
}

func (r *MerchantAutopayoutRepository) Create(m *models.MerchantAutopayout) error {
	return nil
}

func (r *MerchantAutopayoutRepository) Update(m *models.MerchantAutopayout) error {
	return nil
}

func (r *MerchantAutopayoutRepository) Get(m *models.MerchantAutopayout) error {
	return nil
}

func (r *MerchantAutopayoutRepository) Delete(m *models.MerchantAutopayout) error {
	return nil
}

func (r *MerchantAutopayoutRepository) Count(querys interface{}) (int, error) {
	return len(r.ma), nil
}

func (r *MerchantAutopayoutRepository) Selection(querys interface{}) ([]*models.MerchantAutopayout, error) {
	return nil, nil
}
