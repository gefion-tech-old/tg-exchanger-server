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
func TestServer(t *testing.T) (*Server, *redis.Client) {
	t.Helper()

	config := config.InitTestConfig(t)
	assert.NotNil(t, config)

	// Создаю подключение к Redis БД
	redis, err := db.InitRedis(&config.Redis)
	assert.NoError(t, err)
	redisStore := redisstore.Init(redis)

	return root(mocksqlstore.Init(), redisStore, config), redis
}
