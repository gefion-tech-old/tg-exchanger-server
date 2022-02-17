package mocksqlstore

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type DirectionsRepository struct {
	directions map[int]*models.BotMessage
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
