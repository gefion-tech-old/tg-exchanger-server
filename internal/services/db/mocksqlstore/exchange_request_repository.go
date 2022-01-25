package mocksqlstore

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type ExchangeRequestRepository struct {
	er map[int]*models.ExchangeRequest
}

func (r *ExchangeRequestRepository) Create(er *models.ExchangeRequest) error {
	return nil
}

func (r *ExchangeRequestRepository) Update(er *models.ExchangeRequest) error {
	return nil
}

func (r *ExchangeRequestRepository) Delete(er *models.ExchangeRequest) error {
	return nil
}

func (r *ExchangeRequestRepository) Get(er *models.ExchangeRequest) error {
	return nil
}

func (r *ExchangeRequestRepository) Count(querys interface{}) (int, error) {
	return len(r.er), nil
}

func (r *ExchangeRequestRepository) Selection(querys interface{}) ([]*models.ExchangeRequest, error) {
	arr := []*models.ExchangeRequest{}
	for _, er := range r.er {
		arr = append(arr, er)
	}

	return arr, nil
}
