package utils

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type Logger struct {
	store db.LoggerRepository
}

type LoggerI interface {
	NewRecord(r *models.LogRecord) error
}

func InitLogger(s db.LoggerRepository) LoggerI {
	return &Logger{
		store: s,
	}
}

func (u *Logger) NewRecord(r *models.LogRecord) error {
	return u.store.Create(r)
}
