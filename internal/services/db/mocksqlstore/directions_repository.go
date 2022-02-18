package mocksqlstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type DirectionsRepository struct {
	directions map[int]*models.BotMessage

	directionsMaRepository *DirectionsMaRepository
}

func (r *DirectionsRepository) Ma() db.DirectionsMaRepository {
	if r.directionsMaRepository != nil {
		return r.directionsMaRepository
	}

	r.directionsMaRepository = &DirectionsMaRepository{
		dma: make(map[int]*models.DirectionMA),
	}

	return r.directionsMaRepository
}

func (r *DirectionsRepository) Create(er *models.Direction) error {
	return nil
}

func (r *DirectionsRepository) Update(er *models.Direction) error {
	return nil
}

func (r *DirectionsRepository) Delete(er *models.Direction) error {
	return nil
}

func (r *DirectionsRepository) Get(er *models.Direction) error {
	return nil
}

func (r *DirectionsRepository) Count(querys interface{}) (int, error) {
	return len(r.directions), nil
}

func (r *DirectionsRepository) Selection(querys interface{}) ([]*models.Direction, error) {
	return nil, nil
}
