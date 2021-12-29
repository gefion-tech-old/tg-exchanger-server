package mocksqlstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type LoggerRepository struct {
	logs map[int]*models.LogRecord
}

func (r *LoggerRepository) Create(lr *models.LogRecord) error {
	return nil
}

func (r *LoggerRepository) Count() (int, error) {
	return len(r.logs), nil
}

func (r *LoggerRepository) Delete(lr *models.LogRecord) error {
	return nil
}

func (r *LoggerRepository) Selection(page, limit int) ([]*models.LogRecord, error) {
	return nil, nil
}
