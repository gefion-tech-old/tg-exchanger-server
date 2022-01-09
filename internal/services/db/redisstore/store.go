package redisstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/go-redis/redis/v7"
)

type AppRedisDictionaries struct {
	Registration RegistrationClientI
	Auth         AuthClientI
}

// Redis хранилище
type RedisStore struct {
	client *redis.Client
}

func Init(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func InitAppRedisDictionaries(cfg *config.RedisConfig) (*AppRedisDictionaries, func(), error) {
	// Инициализация хранилища кодов верификации
	rRegistration, err := db.InitRedis(cfg, 1)
	if err != nil {
		return nil, nil, err
	}

	// Инициализация хранилища пользовательских сессий
	rAuth, err := db.InitRedis(cfg, 2)
	if err != nil {
		return nil, nil, err
	}

	return &AppRedisDictionaries{
			Registration: InitRegistrationClient(rRegistration),
			Auth:         InitAuthClient(rAuth),
		}, func() {
			rRegistration.Close()
			rAuth.Close()
		}, nil
}
