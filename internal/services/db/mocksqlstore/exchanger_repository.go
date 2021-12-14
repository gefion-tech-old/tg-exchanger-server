package mocksqlstore

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type ExchangerRepository struct {
	exchangers map[uint]*models.Exchanger
}

func (r *ExchangerRepository) Create(e *models.Exchanger) (*models.Exchanger, error) {
	return nil, nil
}

func (r *ExchangerRepository) Update(e *models.Exchanger) (*models.Exchanger, error) {
	return nil, nil
}

func (r *ExchangerRepository) Get(e *models.Exchanger) (*models.Exchanger, error) {
	return nil, nil
}

func (r *ExchangerRepository) Delete(e *models.Exchanger) (*models.Exchanger, error) {
	return nil, nil
}

func (r *ExchangerRepository) Count() (int, error) {
	return len(r.exchangers), nil
}

func (r *ExchangerRepository) GetSlice(limit int) ([]*models.Exchanger, error) {
	return nil, nil
}
