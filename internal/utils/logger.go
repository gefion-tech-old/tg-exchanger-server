package utils

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type Logger struct {
	store db.LoggerRepository
}

type LoggerI interface{}

func InitLogger(s db.LoggerRepository) LoggerI {
	return Logger{
		store: s,
	}
}

func (u *Logger) NewRecord(r *models.LogRecord) {
	u.store.Create(r)
}
