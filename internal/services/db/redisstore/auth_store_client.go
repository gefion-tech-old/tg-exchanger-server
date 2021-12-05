package redisstore

import (
	"strconv"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/go-redis/redis/v7"
)

type AuthClient struct {
	client *redis.Client
}

type AuthClientI interface {
	/*
		Сохранить токен
	*/
	SaveAuth(uuid string, chatID int64, d time.Duration) error
	/*
		Проверить имеет ли пользователь открытую сессию
	*/
	FetchAuth(accessD *models.AccessDetails) (int64, error)

	/*
		Удалить данные пользовательской сессии
	*/
	DeleteAuth(Uuid string) (int64, error)

	/*
		Закрыть установленное соединение с клиентом
	*/
	Close() error

	/*
	   Очистить все хранилище
	*/
	Clear()
}

func InitAuthClient(c *redis.Client) AuthClientI {
	return &AuthClient{client: c}
}

func (c *AuthClient) SaveAuth(uuid string, chatID int64, d time.Duration) error {
	return c.client.Set(uuid, strconv.Itoa(int(chatID)), d).Err()
}

func (c *AuthClient) FetchAuth(accessD *models.AccessDetails) (int64, error) {
	chatID, err := c.client.Get(accessD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	id, _ := strconv.ParseInt(chatID, 10, 64)
	return id, nil
}

func (c *AuthClient) DeleteAuth(Uuid string) (int64, error) {
	deleted, err := c.client.Del(Uuid).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}

func (c *AuthClient) Close() error {
	return c.client.Close()
}

func (c *AuthClient) Clear() {
	c.client.FlushAllAsync()
}
