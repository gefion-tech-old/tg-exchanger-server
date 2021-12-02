package server

import (
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/mocksqlstore"
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
	rr, err := db.InitRedis(&config.Redis, 1)
	assert.NoError(t, err)

	AppRedis := &redisstore.AppRedisDictionaries{
		Registration: rr,
	}

	return root(mocksqlstore.Init(), AppRedis, config), AppRedis, func(clients ...*redis.Client) {
		for _, client := range clients {
			client.FlushAllAsync()
			client.Close()
		}
	}
}
