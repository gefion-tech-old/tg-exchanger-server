package redisstore

import (
	"github.com/go-redis/redis/v7"
)

type AppRedisDictionaries struct {
	Registration *redis.Client
	Auth         *redis.Client
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
