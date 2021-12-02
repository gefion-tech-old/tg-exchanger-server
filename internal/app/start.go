package app

import (
	"context"
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
)

type App struct {
	db     *sql.DB
	redis  *redisstore.AppRedisDictionaries
	config *config.Config
}

type AppI interface {
	Start(ctx context.Context) error
}

func Init(db *sql.DB, config *config.Config) AppI {
	return &App{
		db:     db,
		config: config,
	}
}

func (a *App) Start(ctx context.Context) error {
	// Создание redis хранилища для хранения данных о регистрации пользователя
	rRegistration, err := db.InitRedis(&a.config.Redis, 1)
	if err != nil {
		return err
	}
	defer rRegistration.Close()

	// Создание redis хранилища для хранения пользовательских сессий
	rAuth, err := db.InitRedis(&a.config.Redis, 2)
	if err != nil {
		return err
	}
	defer rAuth.Close()

	// Инициализация БД
	sqlStore := sqlstore.Init(a.db)

	a.redis = &redisstore.AppRedisDictionaries{
		Registration: rRegistration,
		Auth:         rAuth,
	}

	// Инициализация сервера
	server := server.Init(sqlStore, a.redis, a.config)
	return server.Run()
}
