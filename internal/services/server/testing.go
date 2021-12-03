package server

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/mocksqlstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

/*
	Функция возвращает сконфигурированный тестовый сервер
	К тестовому серверу подключается имитация sql хранилища
*/
func TestServer(t *testing.T) (*Server, *redisstore.AppRedisDictionaries, func(...*redis.Client)) {
	t.Helper()

	config := config.InitTestConfig(t)
	assert.NotNil(t, config)

	// Создание redis хранилища для хранения данных о регистрации пользователя
	rRegistration, err := db.InitRedis(&config.Redis, 1)
	assert.NoError(t, err)

	// Создание redis хранилища для хранения пользовательских сессий
	rAuth, err := db.InitRedis(&config.Redis, 2)
	assert.NoError(t, err)

	AppRedis := &redisstore.AppRedisDictionaries{
		Registration: rRegistration,
		Auth:         rAuth,
	}

	// Инициализация соединения с NSQ
	producer, err := db.InitNSQ(&config.NSQ)
	assert.NoError(t, err)

	return root(mocksqlstore.Init(), nsqstore.Init(producer), AppRedis, config), AppRedis, func(clients ...*redis.Client) {
		for _, client := range clients {
			client.FlushAllAsync()
			client.Close()
		}
	}
}

/*
	Метод для быстрой проверки текста ошибки
*/
func TestGetErrorText(t *testing.T, recBody *bytes.Buffer) (string, error) {
	t.Helper()

	var body map[string]interface{}

	if err := json.NewDecoder(recBody).Decode(&body); err != nil {
		return "", err
	}

	return body["error"].(string), nil
}
