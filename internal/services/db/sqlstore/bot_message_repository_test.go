package sqlstore_test

import (
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func Test_SQL_BotMessagesRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	database, teardown := db.TestDB(t, &config.DB)
	defer teardown("users", "bot_messages")

	// Вызываю создание хранилища
	s := sqlstore.Init(database)

	u, err := db.CreateUser(t, s)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	var data *models.BotMessage
	mapstructure.Decode(mocks.BOT_MESSAGE_REQ, &data)
	data.MessageText = "some text"
	data.CreatedBy = mocks.MANAGER_IN_ADMIN_REQ["username"].(string)

	// Создание сообщения
	m, err := s.Manager().BotMessages().Create(data)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	// Получить из БД
	m2, err := s.Manager().BotMessages().Get(m)
	assert.NoError(t, err)
	assert.NotNil(t, m2)

	// Обновление
	m2.MessageText = "new text"
	m3, err := s.Manager().BotMessages().Update(m2)
	assert.NoError(t, err)
	assert.NotNil(t, m3)

	// Удаление
	m4, err := s.Manager().BotMessages().Delete(m)
	assert.NoError(t, err)
	assert.NotNil(t, m4)
}
