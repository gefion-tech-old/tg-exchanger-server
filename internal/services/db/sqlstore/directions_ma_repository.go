package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type DirectionsMaRepository struct {
	store *sql.DB
}

func (r *DirectionsMaRepository) Create(dma *models.DirectionMA) error {
	return nil
}

func (r *DirectionsMaRepository) Update(dma *models.DirectionMA) error {
	return nil
}

func (r *DirectionsMaRepository) Get(dma *models.DirectionMA) error {
	return nil
}

func (r *DirectionsMaRepository) Delete(dma *models.DirectionMA) error {
	return nil
}
