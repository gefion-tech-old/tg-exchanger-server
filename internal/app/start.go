package app

import (
	"context"
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/go-redis/redis/v7"
)

type App struct {
	db     *sql.DB
	redis  *redis.Client
	config *config.Config
}

type AppI interface {
	Start(ctx context.Context) error
}

func Init(db *sql.DB, redis *redis.Client, config *config.Config) AppI {
	return &App{
		db:     db,
		redis:  redis,
		config: config,
	}
}

func (a *App) Start(ctx context.Context) error {
	sqlStore := sqlstore.Init(a.db)
	redisStore := redisstore.Init(a.redis)
	server := server.Init(sqlStore, redisStore, a.config)
	return server.Run()
}
