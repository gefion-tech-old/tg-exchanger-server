package redisstore

import "github.com/go-redis/redis/v7"

// Redis хранилище
type RedisStore struct {
	Client *redis.Client
}

func Init(client *redis.Client) *RedisStore {
	return &RedisStore{
		Client: client,
	}
}
