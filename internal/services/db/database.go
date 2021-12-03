package db

import (
	"database/sql"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/go-redis/redis/v7"
	_ "github.com/lib/pq" // database driver
	"github.com/nsqio/go-nsq"
)

// Функция инициализации Postgres БД
func InitPostgres(config *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		return nil, err
	}

	// Тест соединения
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Функция инициализации Redis БД
func InitRedis(config *config.RedisConfig, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		DB:   db,
	})

	// Тест соединения
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}

// Функция инициализации NSQ очереди
func InitNSQ(config *config.NsqConfig) (*nsq.Producer, error) {
	nsqConf := nsq.NewConfig()
	producer, err := nsq.NewProducer(fmt.Sprintf("%s:%d", config.Host, config.Port), nsqConf)
	if err != nil {
		return nil, err
	}

	// Тест соединения
	if err := producer.Ping(); err != nil {
		return nil, err
	}

	return producer, nil
}
