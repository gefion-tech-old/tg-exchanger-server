package redisstore

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type RegistrationClient struct {
	client *redis.Client
}

type RegistrationClientI interface {
	/*
		Сохранить код верификации
	*/
	SaveVerificationCode(code int, body []byte) error

	/*
		Получить данные соответствующие переданному коде верификации
	*/
	FetchVerificationCode(code int) (string, error)

	/*
		Закрыть установленное соединение с клиентом
	*/
	Close() error

	/*
	   Очистить все хранилище
	*/
	Clear()
}

func InitRegistrationClient(c *redis.Client) RegistrationClientI {
	return &RegistrationClient{client: c}
}

func (c *RegistrationClient) FetchVerificationCode(code int) (string, error) {
	return c.client.Get(fmt.Sprintf("%d", code)).Result()
}

func (c *RegistrationClient) SaveVerificationCode(code int, body []byte) error {
	return c.client.Set(fmt.Sprintf("%d", code), body, 30*time.Minute).Err()
}

func (c *RegistrationClient) Close() error {
	return c.client.Close()
}

func (c *RegistrationClient) Clear() {
	c.client.FlushAllAsync()
}
