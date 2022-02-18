package mocksqlstore

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type DirectionsMaRepository struct {
	dma map[int]*models.DirectionMA
}

func (r *DirectionsMaRepository) Create(er *models.DirectionMA) error {
	return nil
}

func (r *DirectionsMaRepository) Update(er *models.DirectionMA) error {
	return nil
}

func (r *DirectionsMaRepository) Delete(er *models.DirectionMA) error {
	return nil
}

func (r *DirectionsMaRepository) Get(er *models.DirectionMA) error {
	return nil
}

func (r *DirectionsMaRepository) Count(querys interface{}) (int, error) {
	return len(r.dma), nil
}

func (r *DirectionsMaRepository) Selection(querys interface{}) ([]*models.DirectionMA, error) {
	return nil, nil
}
