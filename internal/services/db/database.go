package db

import (
	"database/sql"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/go-redis/redis/v7"
	_ "github.com/lib/pq" // database driver
)

// Функция инициализации Postgres БД
func InitPostgres(config *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		return nil, err
	}

	// Тест соединение
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Функция инициализации Redis БД
func InitRedis(config *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		DB:   int(config.DB),
	})

	// Тест соединение
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}
